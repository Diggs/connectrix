package main

import (
	"github.com/diggs/glog"
	"github.com/tysoft/connectrix/channels"
	"github.com/tysoft/connectrix/channels/http"
	"github.com/tysoft/connectrix/channels/irc"
	"github.com/tysoft/connectrix/config"
	// "github.com/tysoft/connectrix/database"
	"os"
)

func main() {

	log_level := os.Getenv("LOG_LEVEL")
	if log_level == "" {
		log_level = config.Get().LogLevel
	}
	glog.SetSeverity(log_level)
	defer glog.Flush()

	// glog.Info("Connecting to postgres...")
	// err := database.Connect()
	// if err != nil {
	// 	glog.Fatalf("Unable to connect to database: %v", err)
	// }

	glog.Info("Loading channels...")
	err := channels.LoadChannels(
		map[string]channels.PubChannel{
			"http": &http.HttpChannel{},
		},
		map[string]channels.SubChannel{
			"http": &http.HttpChannel{},
			"irc":  &irc.IrcChannel{},
		})
	if err != nil {
		glog.Fatalf("Unable to load channels: %v", err)
	}

	ch := make(chan int)
	<-ch
}
