package routes

import (
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/config"
	"github.com/diggs/connectrix/events/event"
	"github.com/diggs/connectrix/parsers"
	"github.com/diggs/connectrix/templates"
	"github.com/diggs/glog"
	"github.com/diggs/go-eval"
	"sync"
)

var once sync.Once
var routesByPub map[string][]*config.Route = make(map[string][]*config.Route)

func makeRouteKey(namespace string, eventSource string, eventType string) string {
	return fmt.Sprintf("ns:%s:src:%s:type:%s", namespace, eventSource, eventType)
}

func loadRoutes() {
	routes := config.Get().Routes
	for i := range routes {
		route := &routes[i]
		key := makeRouteKey(route.Namespace, route.EventSource, route.EventType)
		if _, exists := routesByPub[key]; !exists {
			routesByPub[key] = []*config.Route{route}
		} else {
			routesByPub[key] = append(routesByPub[key], route)
		}
	}
}

func processEvent(event *event.Event, route *config.Route, channel channels.SubChannel) error {

	// parse the raw content of the event that we can use for templating
	object, err := parsers.Parse(event.RawContent, event.ParserName)
	if err != nil {
		return err
	}

	// template the event, if a custom routing template is specified
	content := event.Content
	if route.Template != "" {
		content, err = templates.Template(object, route.Template)
		if err != nil {
			return err
		}
	}

	// evaluate the routing ruile if specified
	if route.Rule != "" {
		tmplRule, err := templates.Template(object, route.Rule)
		if err != nil {
			return err
		}
		rulePassed, err := goeval.EvalBool(tmplRule)
		if err != nil {
			return err
		}
		// the rule failed, so we shouldn't send the event
		if !rulePassed {
			return nil
		}
	}

	// template each of the routing args
	templatedSubChannelArgs := make(map[string]string, len(route.SubChannelArgs))
	for key, val := range route.SubChannelArgs {
		tmplArg, err := templates.Template(object, val)
		if err != nil {
			return err
		}
		templatedSubChannelArgs[key] = tmplArg
	}

	// send the event
	err = channel.Drain(templatedSubChannelArgs, event, content)
	if err != nil {
		return err
	}

	return nil
}

func RouteEvent(event_ *event.Event) error {

	once.Do(loadRoutes)
	key := makeRouteKey(event_.Namespace, event_.Source, event_.Type)

	if _, exists := routesByPub[key]; exists {
		routes := routesByPub[key]
		for i := range routes {
			route := routes[i]
			channel, err := channels.GetSubChannel(route.SubChannelName)
			if err != nil {
				glog.Warningf("Unable to route event '%v' to '%s': %s", event_, route.SubChannelName, err.Error())
			} else {
				go func(event *event.Event, route *config.Route, channel channels.SubChannel) {
					err := processEvent(event, route, channel)
					if err != nil {
						glog.Warningf("Unable to deliver event '%v' to '%s': %s", event, route.SubChannelName, err.Error())
					}
				}(event_, route, channel)
			}
		}
	}

	return nil
}
