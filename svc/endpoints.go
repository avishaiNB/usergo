package svc

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"
)

// Endpoints ...
type Endpoints struct {
	Logger  Logger
	Tracer  Tracer
	Service Service

	ServerEndpoints []shared.ServerEndpoint
}

// NewEndpoints ...
func NewEndpoints(logger Logger, tracer Tracer, service Service) Endpoints {
	endpoints := Endpoints{
		Logger:  logger,
		Tracer:  tracer,
		Service: service,
	}

	endpoints.AddEndpoints()

	return endpoints
}

// AddEndpoints ...
func (endpoints *Endpoints) AddEndpoints() {
	var serverEndpoints []shared.ServerEndpoint

	userbyid := shared.ServerEndpoint{
		Endpoint: makeUserByIDEndpoint(endpoints.Service),
		Enc:      shared.EncodeReponseToJSON,
		Dec:      decodeUserByIDRequest,
		Method:   "GET",
	}

	endpoints.ServerEndpoints = append(serverEndpoints, userbyid)
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
