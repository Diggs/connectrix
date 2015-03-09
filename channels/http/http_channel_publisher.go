package http

import (
	"errors"
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events"
	"github.com/diggs/glog"
	"io/ioutil"
	"net/http"
)

var ignoreHeadersInHints map[string]int

func (*HttpChannel) PubChannelArgs() []*channels.Arg {
	return nil
}

func (*HttpChannel) ValidatePubChannelArgs(args map[string]string) error {
	return nil
}

func (*HttpChannel) PubChannelInfo(args map[string]string) []*channels.Info {
	// TODO we want to return the full url (including secret) that can be pushed to
	return nil
}

func (ch *HttpChannel) StartPubChannel(config map[string]string, pubChannelArgs []map[string]string) error {

	ignoreHeadersInHints = map[string]int{
		"Content-Length":   0,
		"Host":             0,
		"Accept-Encoding":  0,
		"Encoding":         0,
		"Accept-Language":  0,
		"Accept":           0,
		"Connection":       0,
		"Origin":           0,
		"X-Requested-With": 0,
	}

	port := config["port"]
	http.HandleFunc("/events", ch.handleWebRequest)
	glog.Infof("Starting HTTP channel on %s...", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), LogHandler(http.DefaultServeMux))
}

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Debugf("HTTP %s to %s", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (ch *HttpChannel) handleWebRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Only POST is supported.", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace, err := getNamespace(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = events.ParseAndCreateEventFromChannel(ch.Name(), namespace, &body, getHints(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
}

func getNamespace(r *http.Request) (string, error) {

	// was the namespace in the url?
	query := r.URL.Query()
	if namespace, exists := query["namespace"]; exists {
		return namespace[0], nil
	}

	// or the headers?
	if namespace, exists := r.Header[NAMESPACE_HEADER]; exists {
		return namespace[0], nil
	}

	return "", errors.New(fmt.Sprintf("Unable to determine event namespace. Ensure '?namespace=' query param or '%s' header is set.", NAMESPACE_HEADER))
}

func getHints(r *http.Request) []string {
	hints := []string{}
	for key, val := range r.Header {
		if _, exists := ignoreHeadersInHints[key]; !exists {
			hints = append(hints, fmt.Sprintf("%s:%s", key, val[0]))
		}
	}
	for key, val := range r.URL.Query() {
		hints = append(hints, fmt.Sprintf("%s=%s", key, val[0]))
	}
	return hints
}
