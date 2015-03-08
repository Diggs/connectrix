package parsers

import (
	"errors"
	"fmt"
	"github.com/diggs/connectrix/config"
	"github.com/diggs/connectrix/parsers/json"
	"github.com/diggs/connectrix/parsers/xml"
	"github.com/diggs/connectrix/parsers/yaml"
	"github.com/diggs/glog"
	"regexp"
)

func isPositiveHint(hint string, hints []string) bool {
	for i := range hints {
		match, _ := regexp.Match(fmt.Sprintf(".*%s.*", regexp.QuoteMeta(hint)), []byte(hints[i]))
		if match {
			return true
		}
	}
	return false
}

func makeParser(parserName string) (Parser, error) {
	switch parserName {
	case "json":
		return json.JsonParser{}, nil
	case "xml":
		return xml.XmlParser{}, nil
	case "yaml":
		return yaml.YamlParser{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown parser: '%s'", parserName))
	}
}

func findEventSource(hints []string) (*config.EventSource, error) {

	sources := config.Get().Sources
	for i := range sources {
		if isPositiveHint(sources[i].Hint, hints) {
			return &sources[i], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Unable to identify event source using hints '%v'", hints))
}

func findEventType(eventSource *config.EventSource, hints []string) (*config.EventType, error) {

	eventTypes := eventSource.Events
	for i := range eventTypes {
		if isPositiveHint(eventTypes[i].Hint, hints) {
			return &eventTypes[i], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Unable to identify event type using hints '%v'", hints))
}

func Parse(data *[]byte, parserName string) (interface{}, error) {

	parser, err := makeParser(parserName)
	if err != nil {
		return nil, err
	}

	glog.Debugf("Parsing content using: %s", parserName)
	object, err := parser.ParseContent(data)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func ParseWithHints(data *[]byte, hints []string) (interface{}, *config.EventSource, *config.EventType, error) {

	glog.Debugf("Attempting to parse event using hints: %v", hints)

	eventSource, err := findEventSource(hints)
	if err != nil {
		return nil, nil, nil, err
	}
	glog.Debugf("Identified event source as: %s", eventSource.Name)

	eventType, err := findEventType(eventSource, hints)
	if err != nil {
		return nil, nil, nil, err
	}
	glog.Debugf("Identified event type as: %s", eventType.Type)

	object, err := Parse(data, eventSource.Parser)
	if err != nil {
		return nil, nil, nil, err
	}
	glog.Debugf("Successfully parsed event: %v", object)

	return object, eventSource, eventType, nil
}
