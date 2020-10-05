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

func proxyingMiddleware(ctx context.Context, proxyEndpoints []ProxyEndpoint) ServiceMiddleware {
	commandName := "get-user-by-id"
	config := shared.NewHystrixCommandConfig()
	hystrixConfig := hystrix.CommandConfig{
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,
		Timeout:                config.Timeout,
	}

	hystrix.ConfigureCommand(commandName, hystrixConfig)
	breaker := circuitbreaker.Hystrix(commandName)

	var (
		qps         = 100                    // beyond which we will return an error
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)

	var endpointer sd.FixedEndpointer

	for _, proxyEndpoint := range proxyEndpoints {
		var e endpoint.Endpoint
		e = httptransport.NewClient(proxyEndpoint.method, proxyEndpoint.tgt, proxyEndpoint.enc, proxyEndpoint.dec).Endpoint()
		e = breaker(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// And finally, return the ServiceMiddleware, implemented by proxymw.
	return func(next UserService) UserService {
		return proxymw{ctx, next, retry}
	}
}

type proxymw struct {
	ctx         context.Context
	next        UserService       // Serve most requests via this service...
	getUserByID endpoint.Endpoint // ...except GetUserByID, which gets served by this endpoint
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

func (mw proxymw) GetUserByID(id int) shared.HTTPResponse {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = mw.getUserByID(mw.ctx, id); err != nil {
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
