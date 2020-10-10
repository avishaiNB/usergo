package svc

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"
)

// Endpoints ...
type Endpoints struct {
	Log     core.Log
	Tracer  Tracer
	Service Service

	ServerEndpoints []core.ServerEndpoint
}

// NewEndpoints ...
func NewEndpoints(log core.Log, tracer Tracer, service Service) Endpoints {
	endpoints := Endpoints{
		Log:     log,
		Tracer:  tracer,
		Service: service,
	}

	endpoints.AddEndpoints()

	return endpoints
}

// AddEndpoints ...
func (endpoints *Endpoints) AddEndpoints() {
	var serverEndpoints []core.ServerEndpoint

	userbyid := core.ServerEndpoint{
		Endpoint: makeUserByIDEndpoint(endpoints.Service),
		Enc:      core.EncodeReponseToJSON,
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
