package svc

import (
	"fmt"
	"net/http"
	"os"

	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/shared"
)

// Server ...
type Server struct {
	Name    string
	Address string
	Router  *mux.Router
	Handler http.Handler
	Error   chan error
	Logger  *Logger
	Tracer  *Tracer
}

// NewServer ...
func NewServer(logger *Logger, tracer *Tracer, serviceName string, hostAddress string, errChan chan error) Server {

	return Server{
		Name:    serviceName,
		Address: hostAddress,
		Router:  mux.NewRouter(),
		Error:   errChan,
		Logger:  logger,
		Tracer:  tracer,
	}
}

// Run will create an instance handlers for incoming requests
// it allow to define for each route: handler, decoding requests and encoding responses
// decoding requests may be used for anti corruption layers
func (server Server) Run(endpoints *Endpoints) {

	for _, endpoint := range endpoints.ServerEndpoints {
		getUserByIDHandler := httpkit.NewServer(endpoint.Endpoint, endpoint.Dec, endpoint.Enc)
		server.Router.Methods(endpoint.Method).Path(shared.UserByIDRoute).Handler(getUserByIDHandler)
	}

	server.Handler = handlers.LoggingHandler(os.Stdout, server.Router)
	fmt.Printf("Listernning on %s", server.Address)
	server.Error <- http.ListenAndServe(server.Address, server.Handler)
}
