package client

import (
	"context"

	"github.com/gorilla/mux"
	tlecb "github.com/thelotter-enterprise/usergo/core/circuitbreaker"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlemetrics "github.com/thelotter-enterprise/usergo/core/metrics"
	tleratelimit "github.com/thelotter-enterprise/usergo/core/ratelimit"
	tlesd "github.com/thelotter-enterprise/usergo/core/servicediscovery"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
)

// ServiceClient is a facade for all APIs exposed by the service
type ServiceClient struct {
	Logger      *tlelogger.Manager
	SD          *tlesd.ServiceDiscovery
	ServiceName string
	CB          tlecb.CircuitBreaker
	Limiter     tleratelimit.RateLimiter
	Inst        tlemetrics.PrometheusInstrumentor
	Router      *mux.Router
}

// NewServiceClientWithDefaults with defaults
func NewServiceClientWithDefaults(logger *tlelogger.Manager, sd *tlesd.ServiceDiscovery, serviceName string) ServiceClient {
	return NewServiceClient(
		logger,
		sd,
		tlecb.NewCircuitBreakerator(),
		tleratelimit.NewRateLimitator(),
		tlemetrics.NewPrometheusInstrumentor(serviceName),
		mux.NewRouter(),
		serviceName,
	)
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient(logger *tlelogger.Manager, sd *tlesd.ServiceDiscovery, cb tlecb.CircuitBreaker, limiter tleratelimit.RateLimiter, inst tlemetrics.PrometheusInstrumentor, router *mux.Router, serviceName string) ServiceClient {
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
func (client ServiceClient) GetUserByID(ctx context.Context, id int) tlehttp.Response {
	var service Service
	proxy := NewProxy(client.CB, client.Limiter, client.SD, client.Logger, client.Router)
	instMiddleware := NewInstrumentingMiddleware(client.Inst)
	logMiddleware := NewLoggingMiddleware(client.Logger)
	proxyMiddleware := proxy.UserByIDMiddleware(ctx, id)

	service = proxyMiddleware(service)
	service = logMiddleware(service)
	service = instMiddleware(service)

	res := service.GetUserByID(id)
	return res
}

// GetUserByEmail , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client ServiceClient) GetUserByEmail(ctx context.Context, email string) tlehttp.Response {
	return tlehttp.Response{}
}
