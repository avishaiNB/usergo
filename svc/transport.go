package svc

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
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

// MakeServer will create an instance handlers for incoming requests
// it allow to define for each route: handler, decoding requests and encoding responses
// decoding requests may be used for anti corruption layers
func MakeServer(serviceName, hostAdress, zipkinURL string, endpoints []shared.ServerEndpoint, errChan chan error) Server {
	server := NewServer(serviceName, hostAdress, zipkinURL, errChan)

	for _, endpoint := range endpoints {
		getUserByIDHandler := httptransport.NewServer(endpoint.Endpoint, endpoint.Dec, endpoint.Enc)
		server.Router.Methods("GET").Path(shared.UserByIDRoute).Handler(getUserByIDHandler)
	}

	server.SetHandler(handlers.LoggingHandler(os.Stdout, server.Router))
	return server
}

// decoding request into object (acting as anti corruption layer)
// e.g. url --> GetUserByIDRequest
func decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	req := shared.NewByIDRequest(id)
	return req, nil
}
