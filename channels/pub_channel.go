package channels

// PubChannel can be implemented to allow external systems to publish events
type PubChannel interface {
	// Name returns the name of the channel
	Name() string
	// Description returns a description of the channel
	Description() string
	// Start initializes the channel
	StartPubChannel(map[string]string, []map[string]string) error
	// PubChannelArgs are a a list of names of arguments needed to connect the channel
	PubChannelArgs() []*Arg
	// ValidatePubChannelArgs validates the supplied channel args
	ValidatePubChannelArgs(map[string]string) error
	// PubChannelInfo returns info needed to configure the channel
	PubChannelInfo(map[string]string) []*Info
}
