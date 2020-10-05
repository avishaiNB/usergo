package svc

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

// Server ...
type Server struct {
	Tracer  *zipkingo.Tracer
	Name    string
	Address string
	Router  *mux.Router
	Handler http.Handler
	Error   chan error
}

// NewServer ...
func NewServer(serviceName string, hostAddress string, zipkinURL string, errChan chan error) Server {

	zipkinTracer := makeZipkinEndpoint(serviceName, hostAddress, zipkinURL)

	return Server{
		Tracer:  zipkinTracer,
		Name:    serviceName,
		Address: hostAddress,
		Router:  mux.NewRouter(),
		Error:   errChan,
	}
}

// Run ...
func (s *Server) Run() {
	fmt.Printf("Listernning on %s", s.Address)
	s.Error <- http.ListenAndServe(s.Address, s.Handler)
}

// SetHandler sets the http handler (*mux.route)
func (s *Server) SetHandler(handler http.Handler) {
	s.Handler = handler
}

func makeZipkinEndpoint(serviceName string, hostAddress string, zipkinURL string) *zipkingo.Tracer {
	var zipkinTracer *zipkingo.Tracer
	{
		if zipkinURL != "" {
			var (
				err         error
				hostPort    = hostAddress //"localhost:80"
				serviceName = serviceName
				reporter    = zipkinhttp.NewReporter(zipkinURL)
			)
			defer reporter.Close()
			zEP, _ := zipkingo.NewEndpoint(serviceName, hostPort)
			zipkinTracer, err = zipkingo.NewTracer(reporter, zipkingo.WithLocalEndpoint(zEP))

			if err != nil {
				os.Exit(1)
			}
		}
	}
	return zipkinTracer
}
