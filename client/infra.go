package client

import (
	"net/url"

	"github.com/go-kit/kit/transport/http"
)

// should be in common infra
type ProxyEndpoint struct {
	method string
	tgt    *url.URL
	enc    http.EncodeRequestFunc
	dec    http.DecodeResponseFunc
}
