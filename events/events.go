package events

import (
	// "github.com/diggs/connectrix/database"
	"github.com/diggs/connectrix/config"
	"github.com/diggs/connectrix/events/event"
	"github.com/diggs/connectrix/parsers"
	"github.com/diggs/connectrix/routes"
	"github.com/diggs/connectrix/templates"
)

func CreateEvent(event *event.Event) (int, error) {
	// database.GetDatabase()
	routes.RouteEvent(event)
	return 0, nil
}

func makeTemplatedEventContent(object interface{}, eventType *config.EventType, eventData *[]byte) (string, error) {
	if eventType.Template == "" {
		data_ := *eventData
		return string(data_[:]), nil
	} else {
		return templates.Template(object, eventType.Template)
	}
}

func CreateEventFromChannel(pubChannelName string, namespace string, data *[]byte, hints []string) (int, error) {

	object, eventSource, eventType, err := parsers.ParseWithHints(data, hints)
	if err != nil {
		return -1, err
	}

	content, err := makeTemplatedEventContent(object, eventType, data)
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
