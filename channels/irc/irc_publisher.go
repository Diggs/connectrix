package irc

import (
	"errors"
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events"
	"github.com/diggs/glog"
	irc "github.com/fluffle/goirc/client"
	"regexp"
	"strings"
)

type ircMessage struct {
	// sender contains the nick of the sender
	Sender string
	// rawMsg contains the full line received from IRC excluding the nick of the bot e.g. "@bot echo foo" - rawMsg="echo foo"
	RawMsg string
	// msg is the raw message without the sender or cmd values
	Msg string
	// cmd is the first word in msg e.g. "@bot echo foo" - cmd=echo
	Cmd string
	// args contains the rawMsg broken in to positional args
	Args map[string]string
}

func (ch *IrcChannel) PubChannelArgs() []*channels.Arg {
	return ch.SubChannelArgs() // Same args needed for pub and sub
}

func (ch *IrcChannel) ValidatePubChannelArgs(args map[string]string) error {
	return ch.ValidateSubChannelArgs(args)
}

func (ch *IrcChannel) PubChannelInfo(args map[string]string) []*channels.Info {
	return ch.SubChannelInfo(args)
}

func (ch *IrcChannel) StartPubChannel(channelArgs map[string]string, pubChannelArgs []map[string]string) error {
	for _, args := range pubChannelArgs {
		go func(args map[string]string) {
			ch.connectAndWatch(args)
		}(args)
	}
	return nil
}

func (ch *IrcChannel) connectAndWatch(args map[string]string) {

	connection, err := ch.findOrCreateConnection(args[IRC_SERVER], args[SERVER_PASSWORD], args[IRC_CHANNEL], args[NICKNAME])
	if err != nil {
		// TODO: Want to add some retrying here...
		glog.Debugf("Unable to establish connection to %s:%s: %v", args[IRC_SERVER], args[IRC_CHANNEL], err)
	}

	connection.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {

		senderRegex := regexp.MustCompile("^(.+)!~")
		senderMatches := senderRegex.FindStringSubmatch(line.Src)
		if len(senderMatches) != 2 {
			ch.handleIrcError(args[IRC_CHANNEL], conn, line, errors.New("Unable to determine message sender."))
			return
		}
		sender := senderMatches[1]
		rawMsg := strings.Replace(line.Args[1], "@"+args[NICKNAME]+" ", "", 1)
		msg := strings.Join(strings.Split(rawMsg, " ")[1:], " ")

		argsMap := make(map[string]string)
		split := strings.Split(rawMsg, " ")
		for i, str := range split {
			argsMap[fmt.Sprintf("%d", i)] = str
		}

		m := &ircMessage{Sender: sender, RawMsg: rawMsg, Msg: msg, Args: argsMap}
		rawBytes := []byte(line.Raw)

		// Use the first 'arg' as the only hint - users can then implement routes based on specific commands
		// TODO: How to support namespaces for multitenancy? Could base it on server/channel/nick tuple
		_, err = events.CreateEventFromChannel(ch.Name(), "0", m, &rawBytes, []string{m.Args["0"]})
		if err != nil {
			ch.handleIrcError(args[IRC_CHANNEL], conn, line, err)
			return
		}
	})
}

func (ch *IrcChannel) handleIrcError(ircChannel string, connection *irc.Conn, line *irc.Line, err error) {
	errText := fmt.Sprintf("Unable to handle line: %v - %v", line, err)
	glog.Warningf(errText)
	connection.Privmsg(ircChannel, errText)
}
