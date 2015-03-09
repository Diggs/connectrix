package http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/diggs/connectrix/channels"
	"github.com/diggs/connectrix/events/event"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (*HttpChannel) SubChannelArgs() []*channels.Arg {
	return []*channels.Arg{
		&channels.Arg{
			Name:        URL_ARG,
			Description: "The URL to publish events to.",
			Required:    true,
		},
		&channels.Arg{
			Name:        HEADERS,
			Description: "Custom headers to send when draining the event. Format is a comma seperated string of header:value,header:value",
			Default:     "",
		},
		&channels.Arg{
			Name:        SELF_SIGNED_CERT_ARG,
			Description: "Set top true if URL is using a self signed SSL cert.",
			Default:     "false",
		},
	}
}

func (*HttpChannel) ValidateSubChannelArgs(args map[string]string) error {
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

func (*HttpChannel) SubChannelInfo(map[string]string) []*channels.Info {
	return nil
}

func (*HttpChannel) StartSubChannel(config map[string]string) error {
	return nil
}

func (*HttpChannel) Drain(args map[string]string, event *event.Event, content string) error {

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
