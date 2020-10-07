package svc

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"
)

// HTTPServer ...
type HTTPServer struct {
	Name    string
	Address string
	Router  *mux.Router
	Handler http.Handler
	Logger  Logger
	Tracer  Tracer
}

// NewHTTPServer ...
func NewHTTPServer(logger Logger, tracer Tracer, serviceName string, hostAddress string) HTTPServer {

	return HTTPServer{
		Name:    serviceName,
		Address: hostAddress,
		Router:  mux.NewRouter(),
		Logger:  logger,
		Tracer:  tracer,
	}
}

// Run will create an instance handlers for incoming requests
// it allow to define for each route: handler, decoding requests and encoding responses
// decoding requests may be used for anti corruption layers
func (server HTTPServer) Run(endpoints *Endpoints) error {
	if endpoints == nil {
		return errors.New("no endpoints")
	}

	for _, endpoint := range endpoints.ServerEndpoints {
		getUserByIDHandler := httpkit.NewServer(endpoint.Endpoint, endpoint.Dec, endpoint.Enc)
		server.Router.Methods(endpoint.Method).Path(shared.UserByIDRoute).Handler(getUserByIDHandler)
	}

	server.Handler = handlers.LoggingHandler(os.Stdout, server.Router)
	fmt.Printf("Listernning on %s", server.Address)
	http.ListenAndServe(server.Address, server.Handler)

	return nil
}
