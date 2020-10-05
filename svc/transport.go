package svc

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"
)

// MakeEndpoints creates an instance of Endpoints
func MakeEndpoints(s Service) []shared.ServerEndpoint {
	var serverEndpoints []shared.ServerEndpoint

	userbyid := shared.ServerEndpoint{
		Endpoint: makeUserByIDEndpoint(s),
		Enc:      shared.EncodeReponseToJSON,
		Dec:      decodeUserByIDRequest,
		Method:   "GET",
	}

	serverEndpoints = append(serverEndpoints, userbyid)
	return serverEndpoints
}

func makeUserByIDEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(shared.ByIDRequest)
		user, err := service.GetUserByID(ctx, req.ID)
		return shared.NewUserResponse(user), err
	}
}

// decoding request into object (acting as anti corruption layer)
// e.g. url --> GetUserByIDRequest
func decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	req := shared.NewByIDRequest(id)
	return req, nil
}
