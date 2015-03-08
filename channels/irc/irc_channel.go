package irc

import (
	"fmt"
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
	return "irc"
}

func (IrcChannel) Description() string {
	return "The IRC channel allows events to be sent and received in IRC chat rooms."
}

func (ch IrcChannel) findOrCreateConnection(server string, password string, ircChannel string, nickname string) (*irc.Conn, error) {
	connectionKey := makeConnectionKey(server, ircChannel, nickname)
	if !currentNodeHasConnection(connectionKey) {
		err := ch.establishConnection(server, password, ircChannel, nickname)
		if err != nil {
			return nil, err
		}
	}
	connection := getConnection(connectionKey)
	return connection, nil
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

func (ch IrcChannel) establishConnection(ircServer string, serverPassword string, ircChannel string, nickname string) error {

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
		ch.findOrCreateConnection(ircServer, serverPassword, ircChannel, nickname)
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
