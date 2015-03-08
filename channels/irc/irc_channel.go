package irc

import (
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events"
	"github.com/diggs/connectrix/events/event"
	"github.com/diggs/glog"
	irc "github.com/fluffle/goirc/client"
	"sync"
	"time"
	"regexp"
	"strings"
	"encoding/json"
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
	return "irc"
}

func (IrcChannel) Description() string {
	return "The IRC channel allows events to be sent and received in IRC chat rooms."
}

func (ch IrcChannel) PubChannelArgs() []channels.Arg {
	// Same args needed for pub and sub
	return ch.SubChannelArgs()
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

func (ch IrcChannel) ValidatePubChannelArgs(args map[string]string) error {
	return ch.ValidateSubChannelArgs(args)
}

func (IrcChannel) PubChannelInfo(map[string]string) []channels.Info {
	return nil
}

func (IrcChannel) SubChannelInfo(map[string]string) []channels.Info {
	return nil
}

func (c IrcChannel) StartSubChannel(config map[string]string) error {
	return nil
}

type ircMessage struct {
	From string
	Msg string
	Args map[string]string
}

func (ch IrcChannel) StartPubChannel(channelArgs map[string]string, pubChannelArgs []map[string]string) error {

	for _, args := range pubChannelArgs {
		go func(args map[string]string) {
			connection, err := ch.findOrCreateConnection(args[IRC_SERVER], args[SERVER_PASSWORD], args[IRC_CHANNEL], args[NICKNAME])
			if err != nil {
				glog.Debugf("Unable to establish connection to %s:%s: %v", args[IRC_SERVER], args[IRC_CHANNEL], err)
			}
			connection.HandleFunc("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {

				fromRegex := regexp.MustCompile("^(.+)!~")
				fromMatches := fromRegex.FindStringSubmatch(line.Src)
				from := fromMatches[1]
				msg := strings.TrimLeft(line.Args[1], "@" + args[NICKNAME] + " ")

				argsMap := make(map[string]string)
				split := strings.Split(msg, " ")
				for i, str := range split {
					argsMap[fmt.Sprintf("%d", i)] = str 
				}

				m := &ircMessage{
					From:from,
					Msg:msg,
					Args:argsMap,
				}

				bytes, err := json.Marshal(m)
				if err != nil {
					errText := fmt.Sprintf("Unable to parse irc message: %v - %v", m, err)
					glog.Debugf(errText)
					conn.Privmsg(args[IRC_CHANNEL], errText)
				}

				_, err = events.CreateEventFromChannel(ch.Name(), "0", &bytes, []string{m.Args["0"]})
				if err != nil {
					errText := fmt.Sprintf("Unable to parse irc message: %v - %v", m, err)
					glog.Debugf(errText)
					conn.Privmsg(args[IRC_CHANNEL], errText)
				}
			})
		}(args)
	}

	return nil
}

func (IrcChannel) findOrCreateConnection(server string, password string, ircChannel string, nickname string) (*irc.Conn, error) {
	connectionKey := makeConnectionKey(server, ircChannel, nickname)
	if !currentNodeHasConnection(connectionKey) {
		err := establishConnection(server, password, ircChannel, nickname)
		if err != nil {
			return nil, err
		}
	}
	connection := getConnection(connectionKey)
	return connection, nil
}

func (ch IrcChannel) Drain(args map[string]string, event *event.Event, content string) error {
	connection, err := ch.findOrCreateConnection(args[IRC_SERVER], args[SERVER_PASSWORD], args[IRC_CHANNEL], args[NICKNAME])
	if err != nil {
		return  err
	}
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

func establishConnection(ircServer string, serverPassword string, ircChannel string, nickname string) error {

	glog.Debugf("Connecting to %s on %s as %s", ircChannel, ircServer, nickname)

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
	connections.Lock()
	defer connections.Unlock()
	delete(connections.m, connectionKey)
}

func makeConnectionKey(ircServer string, ircChannel string, nickname string) string {
	return fmt.Sprintf("%s:%s:%s", ircServer, ircChannel, nickname)
}
