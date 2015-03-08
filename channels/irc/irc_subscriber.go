package irc

import (
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events/event"
)

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

func (ch IrcChannel) Drain(args map[string]string, event *event.Event, content string) error {
	connection, err := ch.findOrCreateConnection(args[IRC_SERVER], args[SERVER_PASSWORD], args[IRC_CHANNEL], args[NICKNAME])
	if err != nil {
		return err
	}
	connection.Privmsg(args[IRC_CHANNEL], content)
	return nil
}
