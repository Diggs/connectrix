package channels

import (
	"github.com/diggs/connectrix/events/event"
)

// Info represents a piece of data needed to be passed to an external system to work with a channel
type Info struct {
	Name        string
	Description string
	Value       string
}

// Arg represents a piece of data needed to configure to a channel
type Arg struct {
	Name        string
	Description string
	Default     string
	Required    bool
}

// SubChannel can be implemented to provide a channel that sends events to an external system
type SubChannel interface {
	// Name returns the name of the channel
	Name() string
	// Description returns a description of the channel
	Description() string
	// StartSubChannel initializes the channel
	StartSubChannel(map[string]string) error
	// SubChannelArgs are a a list of names of arguments needed to connect the channel
	SubChannelArgs() []Arg
	// ValidateSubChannelArgs validates the supplied channel args
	ValidateSubChannelArgs(map[string]string) error
	// SubChannelInfo returns info needed to configure the channel
	SubChannelInfo(map[string]string) []Info
	// Drain pushes an event to an external system
	Drain(map[string]string, *event.Event, string) error
}
