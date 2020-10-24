package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	tlecb "github.com/thelotter-enterprise/usergo/core/circuitbreaker"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	tleloadbalancer "github.com/thelotter-enterprise/usergo/core/loadbalancer"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tleratelimit "github.com/thelotter-enterprise/usergo/core/ratelimit"
	tlesd "github.com/thelotter-enterprise/usergo/core/servicediscovery"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"github.com/thelotter-enterprise/usergo/shared"
)

// Proxy ...
type Proxy struct {
	cb         tlecb.CircuitBreaker
	limmitermw endpoint.Middleware
	router     *mux.Router
	limiter    tleratelimit.RateLimiter
	sd         tlesd.ServiceDiscovery
	logger     tlelogger.Manager
}

type userByIDProxyMiddleware struct {
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

// NewProxy ..
func NewProxy(cb tlecb.CircuitBreaker, limiter tleratelimit.RateLimiter, sd *tlesd.ServiceDiscovery, logger tlelogger.Manager, router *mux.Router) Proxy {
	return Proxy{
		cb:      cb,
		limiter: limiter,
		router:  router,
		sd:      *sd,
		logger:  logger,
	}
}

// UserByIDMiddleware ..
func (proxy Proxy) UserByIDMiddleware(ctx context.Context, id int) ServiceMiddleware {
	consulInstancer, _ := proxy.sd.ConsulInstance("user", []string{}, true)
	//consulInstancer := proxy.sd.DNSInstance("user")
	endpointer := sd.NewEndpointer(consulInstancer, proxy.factoryForGetUserByID(ctx, id), proxy.logger.(kitlog.Logger))
	//TODO: refactor. dont like the nil. consider New().With()
	lb := tleloadbalancer.NewLoadBalancer(nil, endpointer)
	retry := lb.DefaultRoundRobinWithRetryEndpoint(ctx)

	return func(next Service) Service {
		return userByIDProxyMiddleware{Context: ctx, Next: next, This: retry}
	}
}

func (proxy Proxy) factoryForGetUserByID(ctx context.Context, id int) sd.Factory {
	// TODO: is this the way to handle path replacement?
	path := fmt.Sprintf(shared.UserByIDClientRoute, id)

	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		breakermw := proxy.cb.NewDefaultHystrixCommandMiddleware("get_user_by_id")
		limitermw := proxy.limiter.NewDefaultErrorLimitterMiddleware()

		tgt, _ := url.Parse(instance) // e.g. parse http://localhost:8080"
		tgt.Path = path

		endpoint := httptransport.NewClient("GET", tgt, encodeGetUserByIDRequest, decodeGetUserByIDResponse, tlectx.WriteBefore()).Endpoint()
		endpoint = breakermw(endpoint)
		endpoint = limitermw(endpoint)

		return endpoint, nil, nil
	}
}

func encodeGetUserByIDRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := shared.NewByIDRequest(ctx, request.(int))
	enc := utils.EncodeRequestToJSON(ctx, r, req)
	return enc
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponseData
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// GetUserByID will execute the endpoint using the middleware and will constract an shared.HTTPResponse
func (proxymw userByIDProxyMiddleware) GetUserByID(id int) tlehttp.Response {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = proxymw.This(proxymw.Context, id); err != nil {
		// TODO: need a refactor to analyze the response
		circuitOpen = true
		statusCode = 500
	}

	return tlehttp.Response{
		Data:          res,
		Error:         err,
		CircuitOpened: circuitOpen,
		StatusCode:    statusCode,
	}
}

// GetUserByEmail will proxy the implementation to the responsible middleware
// We do this to satisfy the service interface
func (proxymw userByIDProxyMiddleware) GetUserByEmail(email string) tlehttp.Response {
	svc := proxymw.Next.(Service)
	return svc.GetUserByEmail(email)
}
