package channels

import (
	"errors"
	"fmt"
	"github.com/diggs/connectrix/config"
	"github.com/diggs/glog"
)

var pubChannels map[string]PubChannel
var subChannels map[string]SubChannel

// TODO move to config.go
func getPubChannelArgs(channelName string) []map[string]string {
	var pubChannelArgs []map[string]string
	sources := config.Get().Sources
	for _, source := range sources {
		if source.PubChannelName == channelName {
			pubChannelArgs = append(pubChannelArgs, source.PubChannelArgs)
		}
	}
	return pubChannelArgs
}

func LoadChannels(pub map[string]PubChannel, sub map[string]SubChannel) error {

	pubChannels = pub
	subChannels = sub

	glog.Info("Loading publishers...")
	for key, val := range pubChannels {
		go func(name string, channel PubChannel) {
			glog.Infof("Starting publish channel %s...", name)
			channel_config := config.Get().Channels[name].Config
			channel_args := getPubChannelArgs(name)
			err := channel.StartPubChannel(channel_config, channel_args)
			if err != nil {
				glog.Warningf("%s failed to start publish channel: %s", channel.Name(), err.Error())
			}
		}(key, val)
	}

	glog.Info("Loading subscrbers...")
	for key, val := range subChannels {
		go func(name string, channel SubChannel) {
			glog.Infof("Starting subscription channel %s...", name)
			err := channel.StartSubChannel(config.Get().Channels[name].Config)
			if err != nil {
				glog.Warningf("%s failed to start subscription channel: %s", channel.Name(), err.Error())
			}
		}(key, val)
	}

	return nil
}

func GetSubChannel(channelName string) (SubChannel, error) {
	if channel, exists := subChannels[channelName]; exists {
		return channel, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown channel: %s", channelName))
	}
}
