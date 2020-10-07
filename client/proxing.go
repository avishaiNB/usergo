package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"
	"golang.org/x/time/rate"
)

func makeProxyMiddleware(in core.ProxyMiddlewareData) UserServiceMiddleware {
	hystrix.ConfigureCommand(in.HystrixCommandName, in.HystrixConfig)
	breaker := circuitbreaker.Hystrix(in.HystrixCommandName)
	var endpointer sd.FixedEndpointer

	for _, proxyEndpoint := range in.ProxyEndpoints {
		var e endpoint.Endpoint
		e = httptransport.NewClient(proxyEndpoint.Method, proxyEndpoint.Tgt, proxyEndpoint.Enc, proxyEndpoint.Dec).Endpoint()
		e = breaker(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), in.MaxQueryPerSecond))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(in.RetryAttempts, in.MaxTimeout, balancer)

	return func(next UserService) UserService {
		out := core.ProxyMiddleware{
			Context: in.Context,
			Next:    next,
			This:    retry,
		}

		return userByIDProxyMiddleware{
			mw: out,
		}
	}
}

func makeEndpoints(id int) []core.ProxyEndpoint {
	var endpoints []core.ProxyEndpoint
	router := mux.NewRouter()
	tgt, _ := router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))

	endpoint1 := core.ProxyEndpoint{
		Method: "GET",
		Tgt:    tgt,
		Enc:    core.EncodeRequestToJSON,
		Dec:    decodeGetUserByIDResponse,
	}

	endpoints = append(endpoints, endpoint1)
	return endpoints
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

type userByIDProxyMiddleware struct {
	mw core.ProxyMiddleware
}

// GetUserByID will execute the endpoint using the middleware and will constract an shared.HTTPResponse
func (proxymw userByIDProxyMiddleware) GetUserByID(id int) core.HTTPResponse {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = proxymw.mw.This(proxymw.mw.Context, id); err != nil {
		// TODO: need a refactor to analyze the response
		circuitOpen = true
		statusCode = 500
	}

	return core.HTTPResponse{
		Result:        res,
		Error:         err,
		CircuitOpened: circuitOpen,
		StatusCode:    statusCode,
	}
}

// GetUserByEmail will proxy the implementation to the responsible middleware
// We do this to satisfy the service interface
func (proxymw userByIDProxyMiddleware) GetUserByEmail(email string) core.HTTPResponse {
	svc := proxymw.mw.Next.(UserService)
	return svc.GetUserByEmail(email)
}
