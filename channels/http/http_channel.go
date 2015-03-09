package http

const (
	URL_ARG              string = "URL"
	HEADERS              string = "Headers"
	SELF_SIGNED_CERT_ARG string = "Self Signed Cert"
	NAMESPACE_HEADER     string = "Connectrix-Namespace"
)

type HttpChannel struct {
}

func (*HttpChannel) Name() string {
	return "http"
}

func (*HttpChannel) Description() string {
	return "The HTTP channel allows events to be sent and received via HTTP requests."
}
