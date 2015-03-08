package main

import (
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/channels/http"
	"github.com/diggs/connectrix/channels/irc"
	"github.com/diggs/connectrix/config"
	"github.com/diggs/glog"
	"os"
)

func main() {

	log_level := os.Getenv("LOG_LEVEL")
	if log_level == "" {
		log_level = config.Get().LogLevel
	}
	glog.SetSeverity(log_level)
	defer glog.Flush()

	glog.Info("Loading channels...")
	err := channels.LoadChannels(
		map[string]channels.PubChannel{
			"http": &http.HttpChannel{},
			"irc":  &irc.IrcChannel{},
		},
		map[string]channels.SubChannel{
			"http": &http.HttpChannel{},
			"irc":  &irc.IrcChannel{},
		})
	if err != nil {
		glog.Fatalf("Unable to load channels: %v", err)
	}

	// live forever
	ch := make(chan int)
	<-ch
}
