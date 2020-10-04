package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	om "github.com/thelotter-enterprise/usergo/shared"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// NewUserByIDEndpoint ...
func NewUserByIDEndpoint(id int) (endpoint.Endpoint, error) {
	// TODO: how do we get the base URL? we need SD
	// TODO: how do we set different strategies for sd based on consul and coreDNS
	// TODO: how CB is handled?
	// TODO: how to retry?
	baseURL, err := url.Parse("http://localhost:8080/")

	if err != nil {
		return nil, err
	}
	endpoint := httptransport.NewClient(
		"GET",
		copyURL(baseURL, "/user/1"),
		encodeHTTPGenericRequest,
		decodeGetUserByIDResponse).Endpoint()

	return endpoint, nil
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp om.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
