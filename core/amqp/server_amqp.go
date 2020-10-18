package amqp

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
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

// Endpoint ...
type Endpoint struct {
	EP       endpoint.Endpoint
	Name     string
	Exchange string
	Queue    string
	Dec      amqptransport.DecodeRequestFunc
}

// Run will ...
func (server *Server) Run(endpoints *[]RabbitMQConsumer) error {
	if endpoints == nil {
		return errors.New("no endpoints")
	}

	//consumers := make(map[string]chan, 1)

	for _, endpoint := range *endpoints {
		_, err := server.RabbitMQ.Consume(&endpoint)

		if err == nil {
			//consumers[endpoint.Consumer] = ch
		}
	}

	return nil
}
