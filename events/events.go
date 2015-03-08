package events

import (
	"github.com/diggs/connectrix/config"
	"github.com/diggs/connectrix/events/event"
	"github.com/diggs/connectrix/parsers"
	"github.com/diggs/connectrix/routes"
	"github.com/diggs/connectrix/templates"
)

func CreateEvent(event *event.Event) (int, error) {
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

func templateAndCreateEvent(eventSource *config.EventSource, eventType *config.EventType, namespace string, object interface{}, data *[]byte) (int, error) {

	content, err := makeTemplatedEventContent(object, eventType, data)
	if err != nil {
		return -1, err
	}

	event := event.Event{
		Namespace:  namespace,
		Source:     eventSource.Name,
		Type:       eventType.Type,
		Content:    content,
		Object:     object,
		ParserName: eventSource.Parser,
	}

	return CreateEvent(&event)
}

func CreateEventFromChannel(pubChannelName string, namespace string, object interface{}, data *[]byte, hints []string) (int, error) {

	eventSource, eventType, err := parsers.IdentifyWithHints(hints)
	if err != nil {
		return -1, err
	}

	return templateAndCreateEvent(eventSource, eventType, namespace, object, data)
}

func ParseAndCreateEventFromChannel(pubChannelName string, namespace string, data *[]byte, hints []string) (int, error) {

	object, eventSource, eventType, err := parsers.ParseWithHints(data, hints)
	if err != nil {
		return -1, err
	}

	return templateAndCreateEvent(eventSource, eventType, namespace, object, data)
}
