package client

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
)

// ServiceClient is a facade for all APIs exposed by the service
type ServiceClient struct {
	Logger      log.Logger
	SD          *core.ServiceDiscovery
	ServiceName string
	CB          core.CircuitBreaker
	Limiter     core.RateLimiter
	Inst        core.Instrumentor
	Router      *mux.Router
}

// NewServiceClientWithDefaults with defaults
func NewServiceClientWithDefaults(logger log.Logger, sd *core.ServiceDiscovery, serviceName string) ServiceClient {

	return NewServiceClient(
		logger,
		sd,
		core.NewCircuitBreakerator(),
		core.NewRateLimitator(),
		core.NewInstrumentor(serviceName),
		mux.NewRouter(),
		serviceName,
	)
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient(logger log.Logger, sd *core.ServiceDiscovery, cb core.CircuitBreaker, limiter core.RateLimiter, inst core.Instrumentor, router *mux.Router, serviceName string) ServiceClient {
	client := ServiceClient{
		Logger:      logger,
		SD:          sd,
		ServiceName: serviceName,
		CB:          cb,
		Limiter:     limiter,
		Inst:        inst,
		Router:      router,
	}
	return client
}

// GetUserByID , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client ServiceClient) GetUserByID(ctx context.Context, id int) core.Response {
	var service UserService
	proxy := NewProxy(client.CB, client.Limiter, client.SD, client.Logger, client.Router)
	instMiddleware := makeInstrumentingMiddleware(client.Inst)
	logMiddleware := makeLoggingMiddleware(client.Logger)
	proxyMiddleware := proxy.UserByIDMiddleware(ctx, id)

	service = proxyMiddleware(service)
	service = logMiddleware(service)
	service = instMiddleware(service)

	res := service.GetUserByID(id)
	return res
}

// GetUserByEmail , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client ServiceClient) GetUserByEmail(ctx context.Context, email string) core.Response {
	return core.Response{}
}
