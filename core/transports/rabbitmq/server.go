package rabbitmq

import (
	"errors"

	"github.com/thelotter-enterprise/usergo/core"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server ...
type Server struct {
	Name     string
	Address  string
	Log      core.Log
	Tracer   tletracer.Tracer
	RabbitMQ *RabbitMQ
}

// NewServer ...
func NewServer(log core.Log, tracer tletracer.Tracer, rabbit *RabbitMQ, serviceName string) Server {
	return Server{
		Name:     serviceName,
		RabbitMQ: rabbit,
		Log:      log,
		Tracer:   tracer,
	}
}

// Run will ...
func (server *Server) Run(endpoints *[]Consumer) error {
	if endpoints == nil {
		return errors.New("no endpoints")
	}

	server.RabbitMQ.OpenConnection()
	//consumers := make(map[string]chan, 1)

	for _, endpoint := range *endpoints {
		_, err := server.RabbitMQ.Consume(&endpoint)

		if err == nil {
			//consumers[endpoint.Consumer] = ch
		}
	}

	return nil
}
