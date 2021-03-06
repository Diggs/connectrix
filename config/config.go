package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type ConnectrixConfig struct {
	DatabaseConnection string `json:"database_connection"`
	LogLevel           string `json:"log_level"`
	Channels           map[string]Channel
	Sources            []*EventSource
	Routes             []*Route
}

type EventSource struct {
	Name           string
	Hint           string
	Parser         string
	Events         []*EventType
	NamedArgs      string            `json:"named_args"`
	PubChannelName string            `json:"pub_channel_name"`
	PubChannelArgs map[string]string `json:"pub_channel_args"`
}

type EventType struct {
	Type     string
	Hint     string
	Fields   []string
	Template string
}

type Route struct {
	NamedArgs      string            `json:"named_args"`
	SubChannelName string            `json:"sub_channel_name"`
	SubChannelArgs map[string]string `json:"sub_channel_args"`
	Namespace      string
	EventSource    string `json:"event_source"`
	EventType      string `json:"event_type"`
	Template       string
	Rule           string
}

type Channel struct {
	Config    map[string]string
	NamedArgs map[string]map[string]string `json:"named_args"`
}

// configPath contains the path to the config file relative to the current process
var configPath string = fmt.Sprintf("%sconfig.json", os.Getenv("CONNECTRIX_CONFIG_FILE"))

// the in-memory config, populated via loadConfig
var config ConnectrixConfig

// used to make sure we only load config once
var once sync.Once

// loadConfig loads config from disk
func loadConfig() {

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to determine absolute config file path:\n %v", err))
	}

	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to load config file:\n %v", err))
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to parse config file:\n %v", err))
	}

	substituteNamedArgs()
}

// substituteNamedArgs sets the PubChannelName/PubChannelArgs properties of EventSources and the
// SubChannelName/SubChannelArgs of Routes to the values specified via the NamedRoutes section of the config.
// This allows arguments to be written once but used in multiple places.
func substituteNamedArgs() {

	type arg struct {
		ChannelName string
		NamedArg    map[string]string
	}

	namedArgs := make(map[string]*arg)
	for channelName, channel := range config.Channels {
		for namedArg, val := range channel.NamedArgs {
			namedArgs[namedArg] = &arg{ChannelName: channelName, NamedArg: val}
		}
	}

	for _, eventSource := range config.Sources {
		if eventSource.NamedArgs != "" {
			if arg, ok := namedArgs[eventSource.NamedArgs]; ok {
				eventSource.PubChannelName = arg.ChannelName
				eventSource.PubChannelArgs = arg.NamedArg
			}
		}
	}

	for _, route := range config.Routes {
		if route.NamedArgs != "" {
			if arg, ok := namedArgs[route.NamedArgs]; ok {
				route.SubChannelName = arg.ChannelName
				route.SubChannelArgs = arg.NamedArg
			}
		}
	}
}

// Get returns the Connectrix config
func Get() *ConnectrixConfig {
	once.Do(loadConfig)
	return &config
}
