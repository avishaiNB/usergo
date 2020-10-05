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

type ProxyEndpoint struct {
	method string
	tgt    *url.URL
	enc    http.EncodeRequestFunc
	dec    http.DecodeResponseFunc
}

type ProxyMiddleware struct {
	In  ProxyMiddlewareInput
	Out ProxyMiddlewareOutput
}

type ProxyMiddlewareOutput struct {
	Context context.Context
	Next    UserService
	This    endpoint.Endpoint
}

type ProxyMiddlewareInput struct {
	Context            context.Context
	HystrixCommandName string
	HystrixConfig      hystrix.CommandConfig
	ProxyEndpoints     []ProxyEndpoint
	MaxQueryPerSecond  int
	RetryAttempts      int
	MaxTimeout         time.Duration
}

func NewMiddlewareInput(ctx context.Context, commandName string, proxyEndpoints []ProxyEndpoint) ProxyMiddlewareInput {
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
