package userclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	om "github.com/thelotter-enterprise/usergo/usershared"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// NewUserByIDEndpoint ...
func NewUserByIDEndpoint(id int) (endpoint.Endpoint, error) {

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

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp om.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
