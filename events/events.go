package events

import (
	// "github.com/tysoft/connectrix/database"
	"github.com/tysoft/connectrix/events/event"
	"github.com/tysoft/connectrix/parsers"
	"github.com/tysoft/connectrix/routes"
	"github.com/tysoft/connectrix/templates"
)

func CreateEvent(event *event.Event) (int, error) {
	// database.GetDatabase()
	routes.RouteEvent(event)
	return 0, nil
}

func CreateEventFromChannel(pubChannelName string, namespace string, data *[]byte, hints []string) (int, error) {

	object, eventSource, eventType, err := parsers.ParseWithHints(data, hints)
	if err != nil {
		return -1, err
	}

	content, err := templates.Template(object, eventType.Template)
	if err != nil {
		return -1, err
	}

	event := event.Event{
		Namespace:  namespace,
		Source:     eventSource.Name,
		Type:       eventType.Type,
		Content:    content,
		RawContent: data,
		ParserName: eventSource.Parser,
	}

	return CreateEvent(&event)
}
