package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
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

// ProxyMiddlewareData holds the return value when we make a middleware
type ProxyMiddlewareData struct {
	// Context holds the context
	Context context.Context

	// Next is a the service instance
	// We need to use Next, since it is used to satisfy the middleware pattern
	// Each middleware is responbsible for a single API, yet, due to the service interface,
	// it need to implement all the service interface APIs. To support it, we use Next to obstract the implementation
	Next interface{}

	// This is the current API which we plan to support in the service interface contract
	This endpoint.Endpoint
}
