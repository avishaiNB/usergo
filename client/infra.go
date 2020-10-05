package client

import (
	"context"
	"net/url"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/thelotter-enterprise/usergo/shared"
)

// ProxyEndpoint holds the information needed to build a go-kit Client
// A Client than can be constructed for a single remote method.
type ProxyEndpoint struct {
	method string
	tgt    *url.URL
	enc    http.EncodeRequestFunc
	dec    http.DecodeResponseFunc
}

// ProxyMiddleware holds the input and output data, which our middleware should have
type ProxyMiddleware struct {
	In  ProxyMiddlewareInput
	Out ProxyMiddlewareOutput
}

// ProxyMiddlewareOutput holds the return value when we make a middleware
type ProxyMiddlewareOutput struct {
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

// ProxyMiddlewareInput holds all the input data required to generate a middleware which supports
// endpoints, circuit breaker, rate limit and timeouts
type ProxyMiddlewareInput struct {
	Context            context.Context
	HystrixCommandName string
	HystrixConfig      hystrix.CommandConfig
	ProxyEndpoints     []ProxyEndpoint
	MaxQueryPerSecond  int
	RetryAttempts      int
	MaxTimeout         time.Duration
}

// MakeDefaultMiddlewareInput creates an opinonated instance of ProxyMiddlewareInput which is common to many simple endpoints
func MakeDefaultMiddlewareInput(ctx context.Context, commandName string, proxyEndpoints []ProxyEndpoint) ProxyMiddlewareInput {
	var (
		qps         = 100                    // beyond which we will return an error
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)

	config := shared.NewHystrixCommandConfig()

	hystrixConfig := hystrix.CommandConfig{
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,

		Timeout: config.Timeout,
	}

	return ProxyMiddlewareInput{
		Context:            ctx,
		HystrixCommandName: commandName,
		HystrixConfig:      hystrixConfig,
		ProxyEndpoints:     proxyEndpoints,
		MaxTimeout:         maxTime,
		RetryAttempts:      maxAttempts,
		MaxQueryPerSecond:  qps,
	}
}
