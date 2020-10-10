package svc

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/transport"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"
)

// HTTPServer ...
type HTTPServer struct {
	Name    string
	Address string
	Router  *mux.Router
	Handler http.Handler
	Log     core.Log
	Tracer  Tracer
}

// NewHTTPServer ...
func NewHTTPServer(log core.Log, tracer Tracer, serviceName string, hostAddress string) HTTPServer {

	return HTTPServer{
		Name:    serviceName,
		Address: hostAddress,
		Router:  mux.NewRouter(),
		Log:     log,
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

	options := []httpkit.ServerOption{
		httpkit.ServerErrorHandler(transport.NewLogErrorHandler(server.Log.Logger)),
		core.ReadCtxBefore(),
	}

	for _, endpoint := range endpoints.ServerEndpoints {
		getUserByIDHandler := httpkit.NewServer(endpoint.Endpoint, endpoint.Dec, endpoint.Enc, options...)
		server.Router.Methods(endpoint.Method).Path(shared.UserByIDRoute).Handler(getUserByIDHandler)
	}

	server.Handler = handlers.LoggingHandler(os.Stdout, server.Router)
	fmt.Printf("Listernning on %s", server.Address)
	http.ListenAndServe(server.Address, server.Handler)

	return nil
}
