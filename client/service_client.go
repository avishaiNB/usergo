package client

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	tlecb "github.com/thelotter-enterprise/usergo/core/cb"
	tleinst "github.com/thelotter-enterprise/usergo/core/inst"
	tlesd "github.com/thelotter-enterprise/usergo/core/sd"
	tletrans "github.com/thelotter-enterprise/usergo/core/transports"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
)

// ServiceClient is a facade for all APIs exposed by the service
type ServiceClient struct {
	Logger      log.Logger
	SD          *tlesd.ServiceDiscovery
	ServiceName string
	CB          tlecb.CircuitBreaker
	Limiter     tletrans.RateLimiter
	Inst        tleinst.Instrumentor
	Router      *mux.Router
}

// NewServiceClientWithDefaults with defaults
func NewServiceClientWithDefaults(logger log.Logger, sd *tlesd.ServiceDiscovery, serviceName string) ServiceClient {

	return NewServiceClient(
		logger,
		sd,
		tlecb.NewCircuitBreakerator(),
		tletrans.NewRateLimitator(),
		tleinst.NewInstrumentor(serviceName),
		mux.NewRouter(),
		serviceName,
	)
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient(logger log.Logger, sd *tlesd.ServiceDiscovery, cb tlecb.CircuitBreaker, limiter tletrans.RateLimiter, inst tleinst.Instrumentor, router *mux.Router, serviceName string) ServiceClient {
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
func (client ServiceClient) GetUserByEmail(ctx context.Context, email string) tlehttp.Response {
	return tlehttp.Response{}
}
