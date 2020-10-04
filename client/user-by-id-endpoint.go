package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// NewUserByIDEndpoint ...
func NewUserByIDEndpoint(id int) (endpoint.Endpoint, error) {
	// TODO: how do we get the base URL? we need SD
	// TODO: how do we set different strategies for sd based on consul and coreDNS
	// TODO: how CB is handled?
	// TODO: how to retry?
	router := mux.NewRouter()
	u, err := router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))

	if err != nil {
		return nil, err
	}

	endpoint := httptransport.NewClient("GET", u, shared.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint()

	return endpoint, nil
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
