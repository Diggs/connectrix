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
	Channels           map[string]map[string]string
	Sources            []EventSource
	Routes             []Route
}

type EventSource struct {
	Name   string
	Hint   string
	Parser string
	Events []EventType
}

type EventType struct {
	Type     string
	Hint     string
	Fields   []string
	Template string
}

type Route struct {
	SubChannelName string            `json:"sub_channel_name"`
	SubChannelArgs map[string]string `json:"sub_channel_args"`
	Namespace      string
	EventSource    string `json:"event_source"`
	EventType      string `json:"event_type"`
	Template       string
	Rule           string
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
}

// Get returns the Connectrix config
func Get() *ConnectrixConfig {
	once.Do(loadConfig)
	return &config
}
