package http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events"
	"github.com/diggs/connectrix/events/event"
	"github.com/diggs/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const URL_ARG string = "URL"
const HEADERS string = "Headers"
const SELF_SIGNED_CERT_ARG = "Self Signed Cert"
const NAMESPACE_HEADER = "Connectrix-Namespace"

var ignoreHeadersInHints map[string]int

type HttpChannel struct {
}

func (HttpChannel) Name() string {
	return "HTTP"
}

func (HttpChannel) Description() string {
	return "The HTTP channel allows events to be published via HTTP requests."
}

func (HttpChannel) PubChannelArgs() []channels.Arg {
	// no args needed
	return nil
}

func (HttpChannel) ValidatePubChannelArgs(args map[string]string) error {
	// no validation needed as this channel doesn't need any source args
	return nil
}

func (HttpChannel) PubChannelInfo(args map[string]string) []channels.Info {
	// TODO we want to return the full url (including secret) that can be pushed to
	return nil
}

func (HttpChannel) SubChannelArgs() []channels.Arg {
	return []channels.Arg{
		channels.Arg{
			Name:        URL_ARG,
			Description: "The URL to publish events to.",
			Required:    true,
		},
		channels.Arg{
			Name:        HEADERS,
			Description: "Custom headers to send when draining the event. Format is a comma seperated string of header:value,header:value",
			Default:     "",
		},
		channels.Arg{
			Name:        SELF_SIGNED_CERT_ARG,
			Description: "Set top true if URL is using a self signed SSL cert.",
			Default:     "false",
		},
	}
}

func (HttpChannel) ValidateSubChannelArgs(args map[string]string) error {
	// ensure the url is valid
	url, err := url.Parse(args[URL_ARG])
	if err != nil {
		return err
	}
	if url.Host == "" {
		return errors.New("URL must be fully qualified")
	}

	// ensure this is a bool
	_, err = strconv.ParseBool(args[SELF_SIGNED_CERT_ARG])
	if err != nil {
		return err
	}

	return nil
}

func (HttpChannel) SubChannelInfo(map[string]string) []channels.Info {
	return nil
}

func (c HttpChannel) StartSubChannel(config map[string]string) error {
	// nothing to do
	return nil
}

func (c HttpChannel) StartPubChannel(config map[string]string, pubChannelArgs []map[string]string) error {

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
	http.HandleFunc("/events", c.handleWebRequest)
	glog.Infof("Starting HTTP channel on %s...", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), LogHandler(http.DefaultServeMux))
}

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Debugf("HTTP %s to %s", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (HttpChannel) handleWebRequest(w http.ResponseWriter, r *http.Request) {

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

	_, err = events.CreateEventFromChannel(HttpChannel{}.Name(), namespace, &body, getHints(r))
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

func getCustomHeaders(args map[string]string) map[string]string {
	customHeaders := make(map[string]string)
	if _, exists := args[HEADERS]; exists {
		headerStr := args[HEADERS]
		headers := strings.Split(headerStr, ",")
		for i := range headers {
			headerSplit := strings.Split(headers[i], ":")
			if len(headerSplit) == 2 {
				customHeaders[strings.Trim(headerSplit[0], " ")] = strings.Trim(headerSplit[1], " ")
			}
		}
	}
	return customHeaders
}

func (HttpChannel) Drain(args map[string]string, event *event.Event, content string) error {

	// args are validated via ValidateSinkArgs, assume they're correct here
	url := args[URL_ARG]
	selfSignedCert, _ := strconv.ParseBool(args[SELF_SIGNED_CERT_ARG])

	// use a custom transport so we can support accepting invalid certs (if enabled)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: selfSignedCert},
	}
	client := &http.Client{Transport: tr}

	// set up the request, including custom user-agent
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		return err
	}

	// add custom headers if any (format should be "header:value,header:value")
	for key, val := range getCustomHeaders(args) {
		req.Header.Set(key, val)
	}
	req.Header.Set(NAMESPACE_HEADER, event.Namespace)
	req.Header.Set("Content-Type", "application/json") // todo could support the type coming from the channel args
	req.Header.Set("User-Agent", "connectrix/http")

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if we don't get a 2XX response code then this was a failure
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("HTTP POST to %s failed with status code %s. Response: %s", url, resp.Status, string(body[:])))
	}

	return nil
}
