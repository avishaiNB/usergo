package core

import (
	"context"

	httpkit "github.com/go-kit/kit/transport/http"
)

// HTTPResponse is a base response
type HTTPResponse struct {
	Error         error
	Result        interface{}
	CircuitOpened bool
	Context       context.Context
	StatusCode    int
	// CorrelactionID
	// timeout
}

// HTTPRequest is a base response
// TODO: use it as the request wrapper
type HTTPRequest struct {
	Request interface{}
	Context context.Context
	// CorrelactionID
	// timeout
}

// ServerEndpoint holds the information needed to build a server endpoint which client can call upon
type ServerEndpoint struct {
	Method   string
	Endpoint func(ctx context.Context, request interface{}) (interface{}, error)
	Dec      httpkit.DecodeRequestFunc
	Enc      httpkit.EncodeResponseFunc
}
