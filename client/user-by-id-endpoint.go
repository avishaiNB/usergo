package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type userByIDEndpoint struct {
	Router  *mux.Router
	EP      endpoint.Endpoint
	URL     *url.URL
	ID      int
	Context context.Context
	Result  shared.ByIDResponse
	Err     error
}

func newUserByIDEndpoint(ctx context.Context, id int, router *mux.Router) userByIDEndpoint {
	return userByIDEndpoint{
		Context: ctx,
		ID:      id,
		Router:  router,
	}
}

// MakeUserByIDEndpoint ...
func (ep *userByIDEndpoint) build() {
	// TODO: how do we get the base URL? we need SD
	// TODO: how do we set different strategies for sd based on consul and coreDNS
	// TODO: how CB is handled?
	// TODO: how to retry?
	ep.URL, _ = ep.Router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(ep.ID))
	ep.EP = httptransport.NewClient("GET", ep.URL, shared.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint()

	return
}

func (ep *userByIDEndpoint) exec() {
	res, err := ep.EP(ep.Context, ep.ID)
	response := res.(shared.ByIDResponse)
	ep.Result = response
	ep.Err = err
}

func (ep *userByIDEndpoint) result() (shared.ByIDResponse, error) {
	return ep.Result, ep.Err
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
