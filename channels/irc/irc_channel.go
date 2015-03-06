package irc

import (
	"errors"
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events/event"
	"github.com/diggs/glog"
	irc "github.com/fluffle/goirc/client"
	"sync"
	"time"
)

const IRC_SERVER string = "IRC Server"
const IRC_CHANNEL string = "IRC Channel"
const NICKNAME string = "Nickname"
const SERVER_PASSWORD string = "Server Password"

var connections = struct {
	sync.RWMutex
	m map[string]*irc.Conn
}{m: make(map[string]*irc.Conn)}

type IrcChannel struct {
}

func (IrcChannel) Name() string {
	return "IRC"
}

func (IrcChannel) Description() string {
	return "The IRC channel allows events to be sent to IRC chat rooms."
}

func (IrcChannel) SubChannelArgs() []channels.Arg {
	return []channels.Arg{
		channels.Arg{
			Name:        IRC_SERVER,
			Description: "The IRC server to connect to.",
			Required:    true,
		},
		channels.Arg{
			Name:        SERVER_PASSWORD,
			Description: "The password to connect to the IRC server with.",
			Default:     "",
		},
		channels.Arg{
			Name:        IRC_CHANNEL,
			Description: "The IRC channel to connect to.",
			Required:    true,
		},
		channels.Arg{
			Name:        NICKNAME,
			Description: "The nickname to connect to the IRC server with.",
			Required:    true,
			Default:     "Connectrix",
		},
	}
}

func (IrcChannel) ValidateSubChannelArgs(args map[string]string) error {
	return nil
}

func (IrcChannel) SubChannelInfo(map[string]string) []channels.Info {
	return nil
}

func (c IrcChannel) StartSubChannel(config map[string]string) error {
	return nil
}

func (IrcChannel) Drain(args map[string]string, event *event.Event, content string) error {

	connectionKey := makeConnectionKey(args[IRC_SERVER], args[IRC_CHANNEL], args[NICKNAME])

	if !currentNodeHasConnection(connectionKey) {
		if anyNodeHasConnection(connectionKey) {
			return drainToExternalNode(args, event, content)
		} else {
			err := establishConnection(args[IRC_SERVER], args[SERVER_PASSWORD], args[IRC_CHANNEL], args[NICKNAME])
			if err != nil {
				return err
			}
		}
	}

	connection := getConnection(connectionKey)
	connection.Privmsg(args[IRC_CHANNEL], content)

	return nil
}

func currentNodeHasConnection(connectionKey string) bool {
	connections.RLock()
	defer connections.RUnlock()
	_, exists := connections.m[connectionKey]
	return exists
}

func getConnection(connectionKey string) *irc.Conn {
	connections.RLock()
	defer connections.RUnlock()
	return connections.m[connectionKey]
}

func anyNodeHasConnection(connectionKey string) bool {
	// TODO Check centralized hash
	return false
}

func establishConnection(ircServer string, serverPassword string, ircChannel string, nickname string) error {

	glog.Debugf("Connecting to %s on %s as %s", ircChannel, ircServer, nickname)

	// TODO Add to centralized hash
	connectionKey := makeConnectionKey(ircServer, ircChannel, nickname)

	// lock the connections map for writing
	connections.Lock()
	defer connections.Unlock()

	// make sure another routine didn't get in first
	// (don't use currentNodeHasConnection() as we are holding a writer lock)
	if _, exists := connections.m[connectionKey]; exists {
		return nil
	}

	connectedToChannel := make(chan bool, 1)
	config := irc.NewConfig(nickname)
	config.Server = ircServer
	config.Pass = serverPassword
	client := irc.Client(config)

	client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		glog.Debugf("Connected to %s", ircServer)
		conn.Join(ircChannel)
	})

	client.HandleFunc("join", func(conn *irc.Conn, line *irc.Line) {
		if line.Nick == nickname {
			glog.Debugf("Joined %s:%s", ircServer, ircChannel)
			connectedToChannel <- true
		}
	})

	client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		glog.Debugf("Disconnected from %s:%s", ircServer, ircChannel)
		removeConnection(connectionKey)
	})

	if err := client.Connect(); err != nil {
		glog.Debugf("Unable to connect to %s:%s: %s", ircServer, ircChannel, err.Error())
		return err
	}
	connections.m[connectionKey] = client

	select {
	case <-connectedToChannel:
		return nil
	case <-time.After(30 * time.Second):
		connectedToChannel <- false
		return fmt.Errorf("Timed out connecting to %s:%s", ircServer, ircChannel)
	}
}

func removeConnection(connectionKey string) {
	// TODO Remove from centralized hash
	connections.Lock()
	defer connections.Unlock()
	delete(connections.m, connectionKey)
}

func drainToExternalNode(args map[string]string, event *event.Event, content string) error {
	return errors.New("Not Implemented")
}

func makeConnectionKey(ircServer string, ircChannel string, nickname string) string {
	return fmt.Sprintf("%s:%s:%s", ircServer, ircChannel, nickname)
}
