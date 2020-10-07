package core

import (
	"context"
	"net/url"
	"time"

	"github.com/afex/hystrix-go/hystrix"
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

// ProxyEndpoint holds the information needed to build a go-kit Client
// A Client than can be constructed for a single remote method.
type ProxyEndpoint struct {
	Method string
	Tgt    *url.URL
	Enc    httpkit.EncodeRequestFunc
	Dec    httpkit.DecodeResponseFunc
}

// ServerEndpoint holds the information needed to build a server endpoint which client can call upon
type ServerEndpoint struct {
	Method   string
	Endpoint func(ctx context.Context, request interface{}) (interface{}, error)
	Dec      httpkit.DecodeRequestFunc
	Enc      httpkit.EncodeResponseFunc
}

// ProxyMiddleware holds the return value when we make a middleware
type ProxyMiddleware struct {
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

// ProxyMiddlewareData holds all the input data required to generate a middleware which supports
// endpoints, circuit breaker, rate limit and timeouts
type ProxyMiddlewareData struct {
	Context            context.Context
	HystrixCommandName string
	HystrixConfig      hystrix.CommandConfig
	ProxyEndpoints     []ProxyEndpoint
	RetryAttempts      int
	MaxTimeout         time.Duration
}

// MakeProxyMiddlewareData creates an opinonated instance of ProxyMiddlewareInput which is common to many simple endpoints
func MakeProxyMiddlewareData(ctx context.Context, commandName string, proxyEndpoints []ProxyEndpoint) ProxyMiddlewareData {
	var (
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)

	config := NewHystrixCommandConfig()

	hystrixConfig := hystrix.CommandConfig{
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,

		Timeout: config.Timeout,
	}

	return ProxyMiddlewareData{
		Context:            ctx,
		HystrixCommandName: commandName,
		HystrixConfig:      hystrixConfig,
		ProxyEndpoints:     proxyEndpoints,
		MaxTimeout:         maxTime,
		RetryAttempts:      maxAttempts,
	}
}
