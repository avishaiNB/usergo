package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	om "github.com/thelotter-enterprise/usergo/usershared"
)

// Endpoints holds all the endpoints which are supported by the service
type Endpoints struct {
	GetUserByID endpoint.Endpoint
}

// MakeEndpoints creates an instance of Endpoints
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserByID: makeUserByIDEndpoint(s),
	}
}

func makeUserByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(om.ByIDRequest)
		user, err := s.GetUserByID(ctx, req.ID)
		return om.NewUserResponse(user), err
	}
}

//NewServer will create an instance handlers for incoming requests
// it allow to define for each route: handler, decoding requests and encoding responses
// decoding requests may be used for anti corruption layers
func NewServer(ctx context.Context, endpoints Endpoints) http.Handler {
	router := mux.NewRouter()
	getUserByIDHandler := httptransport.NewServer(endpoints.GetUserByID, decodeUserByIDRequest, encodeReponseToJSON)
	router.Methods("GET").Path("/user/{id}").Handler(getUserByIDHandler)

	return handlers.LoggingHandler(os.Stdout, router)
}

// encoding the response into json
// e.g. GetUsetByIDResponse --> json
func encodeReponseToJSON(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// decoding request into object (acting as anti corruption layer)
// e.g. url --> GetUserByIDRequest
func decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	req := om.NewByIDRequest(id)
	fmt.Println(req)
	return req, nil
}
