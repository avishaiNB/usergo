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

// Endpoints holds all the endpoints which are supported by the service
type Endpoints struct {
	GetUserByID endpoint.Endpoint
}

// MakeEndpoints creates an instance of Endpoints
func MakeEndpoints(s Service) Endpoints {
	userbyidEndpoint := makeUserByIDEndpoint(s)
	//userbyidEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Sum")(userbyidEndpoint)
	//shared.ProxyEndpoint
	return Endpoints{
		GetUserByID: userbyidEndpoint,
	}
}

func makeUserByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(shared.ByIDRequest)
		user, err := s.GetUserByID(ctx, req.ID)
		return shared.NewUserResponse(user), err
	}
}

// MakeServer will create an instance handlers for incoming requests
// it allow to define for each route: handler, decoding requests and encoding responses
// decoding requests may be used for anti corruption layers
func MakeServer(serviceName, hostAdress, zipkinURL string, endpoints Endpoints, errChan chan error) Server {
	server := NewServer(serviceName, hostAdress, zipkinURL, errChan)
	getUserByIDHandler := httptransport.NewServer(endpoints.GetUserByID, decodeUserByIDRequest, shared.EncodeReponseToJSON)
	server.Router.Methods("GET").Path(shared.UserByIDRoute).Handler(getUserByIDHandler)
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
