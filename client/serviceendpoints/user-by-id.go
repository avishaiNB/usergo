package serviceendpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

// UserByIDServiceClient ...
type UserByIDServiceClient struct {
	Router    *mux.Router
	Endpoints []endpoint.Endpoint
	ID        int
	Context   context.Context
	Result    shared.ByIDResponse
	Err       error
}

// NewUserByIDServiceClient ...
func NewUserByIDServiceClient(ctx context.Context, id int, router *mux.Router) UserByIDServiceClient {
	return UserByIDServiceClient{
		Context: ctx,
		ID:      id,
		Router:  router,
	}
}

// BuildEndpoints will look for all the service endpoints
func (serviceClient *UserByIDServiceClient) BuildEndpoints() {
	url1, _ := serviceClient.Router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(serviceClient.ID))
	ep1 := httptransport.NewClient("GET", url1, shared.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint()

	// This is a non existing URL
	url2, _ := serviceClient.Router.Schemes("http").Host("localhost:8081").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(serviceClient.ID))
	ep2 := httptransport.NewClient("GET", url2, shared.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint()

	serviceClient.Endpoints = []endpoint.Endpoint{
		ep2,
		ep1,
	}
}

// Build ...
func (serviceClient *UserByIDServiceClient) BuildCircuitBreaker() {
	// Set some parameters for our client.
	// var (
	// 	maxAttempts = 3                      // per request, before giving up
	// 	maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	// )

	// var endpointer sd.FixedEndpointer
	// var e endpoint.Endpoint = ep.EP
	// endpointer = append(endpointer, e)

	// // Now, build a single, retrying, load-balancing endpoint out of all of
	// // those individual endpoints.
	// balancer := lb.NewRoundRobin(endpointer)
	// retry := lb.Retry(maxAttempts, maxTime, balancer)

	return
}

// Exec ...
func (serviceClient *UserByIDServiceClient) Exec() {
	endpointer := sd.FixedEndpointer(serviceClient.Endpoints)
	balancer := lb.NewRoundRobin(endpointer)

	var res interface{}
	var err error
	for i := 0; i < len(serviceClient.Endpoints); i++ {
		ep, _ := balancer.Endpoint()
		res, err = ep(serviceClient.Context, serviceClient.ID)

		if err == nil {
			break
		}
	}

	response := res.(shared.ByIDResponse)
	serviceClient.Result = response
	serviceClient.Err = err
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
