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
	"github.com/thelotter-enterprise/usergo/shared"
	"golang.org/x/time/rate"
)

func makeUserByIDMiddleware(in ProxyMiddlewareInput) UserServiceClientMiddleware {
	hystrix.ConfigureCommand(in.HystrixCommandName, in.HystrixConfig)
	breaker := circuitbreaker.Hystrix(in.HystrixCommandName)
	var endpointer sd.FixedEndpointer

	for _, proxyEndpoint := range in.ProxyEndpoints {
		var e endpoint.Endpoint
		e = httptransport.NewClient(proxyEndpoint.method, proxyEndpoint.tgt, proxyEndpoint.enc, proxyEndpoint.dec).Endpoint()
		e = breaker(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), in.MaxQueryPerSecond))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(in.RetryAttempts, in.MaxTimeout, balancer)

	return func(next UserServiceClient) UserServiceClient {
		out := ProxyMiddlewareOutput{in.Context, next, retry}
		return ProxyMiddleware{
			Out: out,
			In:  in,
		}
	}
}

func makeUserByIDEndpoints(id int) []ProxyEndpoint {
	var endpoints []ProxyEndpoint
	router := mux.NewRouter()
	tgt, _ := router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))

	endpoint1 := ProxyEndpoint{
		method: "GET",
		tgt:    tgt,
		enc:    shared.EncodeRequestToJSON,
		dec:    decodeGetUserByIDResponse,
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

// GetUserByID will execute the endpoint using the middleware and will constract an shared.HTTPResponse
func (mw ProxyMiddleware) GetUserByID(id int) shared.HTTPResponse {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = mw.Out.This(mw.In.Context, id); err != nil {
		// TODO: need a refactor to analyze the response
		circuitOpen = true
		statusCode = 500
	}

	return shared.HTTPResponse{
		Result:        res,
		Error:         err,
		CircuitOpened: circuitOpen,
		StatusCode:    statusCode,
	}
}

// GetUserByEmail will proxy the implementation to the responsible middleware
// We do this to satisfy the service interface
func (mw ProxyMiddleware) GetUserByEmail(email string) shared.HTTPResponse {
	return mw.Out.Next.GetUserByEmail(email)
}
