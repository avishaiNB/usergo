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

func proxyingMiddleware(in ProxyMiddlewareInput) ServiceMiddleware {
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

	return func(next UserService) UserService {
		out := ProxyMiddlewareOutput{in.Context, next, retry}
		return ProxyMiddleware{
			Out: out,
			In:  in,
		}
	}
}

func makeProxyEndpoint(id int) ProxyEndpoint {
	router := mux.NewRouter()
	tgt, _ := router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))

	return ProxyEndpoint{
		method: "GET",
		tgt:    tgt,
		enc:    shared.EncodeRequestToJSON,
		dec:    decodeGetUserByIDResponse,
	}
}

func (mw ProxyMiddleware) GetUserByID(id int) shared.HTTPResponse {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = mw.Out.This(mw.In.Context, id); err != nil {
		circuitOpen = true
		statusCode = 500
	}

	return shared.HTTPResponse{
		// shared.ByIDResponse,
		Result:        res,
		Error:         err,
		CircuitOpened: circuitOpen,
		StatusCode:    statusCode,
	}
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
