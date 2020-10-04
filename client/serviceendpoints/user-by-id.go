package serviceendpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

// UserByIDServiceClient ...
type UserByIDServiceClient struct {
	Router               *mux.Router
	Endpoints            []endpoint.Endpoint
	HystrixCommandConfig hystrix.CommandConfig
	HystrixCommandName   string
	Params               map[string]interface{}
	Context              context.Context
	// TODO we need a result and we need to include if circuit was opened
	Result shared.ByIDResponse
	Err    error
}

// NewUserByIDServiceClient ...
func NewUserByIDServiceClient(router *mux.Router) UserByIDServiceClient {
	return UserByIDServiceClient{
		Router: router,
	}
}

// WithParams ...
func (serviceClient *UserByIDServiceClient) WithParams(params map[string]interface{}) {
	serviceClient.Params = params
}

// WithContext ...
func (serviceClient *UserByIDServiceClient) WithContext(ctx context.Context) {
	serviceClient.Context = ctx
}

// WithCircuitBreaker ...
func (serviceClient *UserByIDServiceClient) WithCircuitBreaker(commandName string, config shared.HystrixCommandConfig) {
	c := hystrix.CommandConfig{
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,
		Timeout:                config.Timeout,
	}
	serviceClient.HystrixCommandConfig = c
	serviceClient.HystrixCommandName = commandName
}

// BuildEndpoints will look for all the service endpoints
func (serviceClient *UserByIDServiceClient) BuildEndpoints() {
	id := serviceClient.Params["ID"].(int)

	hystrix.ConfigureCommand(serviceClient.HystrixCommandName, serviceClient.HystrixCommandConfig)
	breaker := circuitbreaker.Hystrix(serviceClient.HystrixCommandName)

	url1, _ := serviceClient.Router.Schemes("http").Host("localhost:8082").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))
	ep1 := breaker(httptransport.NewClient("GET", url1, shared.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint())

	// This is a non existing URL
	url2, _ := serviceClient.Router.Schemes("http").Host("localhost:8081").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))
	ep2 := breaker(httptransport.NewClient("GET", url2, shared.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint())

	serviceClient.Endpoints = []endpoint.Endpoint{
		ep2,
		ep1,
	}
}

// Exec ...
func (serviceClient *UserByIDServiceClient) Exec() {
	id := serviceClient.Params["ID"].(int)
	var res interface{}
	var err error
	endpointer := sd.FixedEndpointer(serviceClient.Endpoints)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(1000, 10000*time.Millisecond, balancer)

	if res, err = retry(serviceClient.Context, id); err != nil {
		panic(err)
	}

	response := res.(shared.ByIDResponse)
	serviceClient.Result = response
	serviceClient.Err = err

	return
}

// GetResult ...GetResult
func (serviceClient *UserByIDServiceClient) GetResult() (shared.ByIDResponse, error) {
	return serviceClient.Result, serviceClient.Err
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
