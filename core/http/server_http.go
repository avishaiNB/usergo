package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/transport"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// HTTPServer ...
type HTTPServer struct {
	Name    string
	Address string
	Router  *mux.Router
	Handler http.Handler
	Log     core.Log
	Tracer  tletracer.Tracer
}

// HTTPEndpoints ...
type HTTPEndpoints struct {
	ServerEndpoints []HTTPEndpoint
}

// HTTPEndpoint holds the information needed to build a server endpoint which client can call upon
type HTTPEndpoint struct {
	Method   string
	Endpoint func(ctx context.Context, request interface{}) (interface{}, error)
	Dec      httpkit.DecodeRequestFunc
	Enc      httpkit.EncodeResponseFunc
	Path     string
}

// NewHTTPServer ...
func NewHTTPServer(log core.Log, tracer tletracer.Tracer, serviceName string, hostAddress string) HTTPServer {
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
func (server *HTTPServer) Run(endpoints *HTTPEndpoints) error {
	if endpoints == nil {
		return errors.New("no endpoints")
	}

	c := core.NewCtx()

	options := []httpkit.ServerOption{
		httpkit.ServerErrorHandler(transport.NewLogErrorHandler(server.Log.Logger)),
		c.ReadBefore(),
	}

	for _, endpoint := range endpoints.ServerEndpoints {
		server.Log.Logger.Log("message", fmt.Sprintf("adding route http://%s/%s", server.Address, endpoint.Path))
		getUserByIDHandler := httpkit.NewServer(endpoint.Endpoint, endpoint.Dec, endpoint.Enc, options...)
		server.Router.Methods(endpoint.Method).Path(endpoint.Path).Handler(getUserByIDHandler)
	}

	server.Handler = handlers.LoggingHandler(os.Stdout, server.Router)
	server.Log.Logger.Log("message", fmt.Sprintf("http server started and listen on %s", server.Address))
	http.ListenAndServe(server.Address, server.Handler)

	return nil
}
