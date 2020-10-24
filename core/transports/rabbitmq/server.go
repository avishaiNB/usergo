package rabbitmq

import (
	"errors"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server ...
type Server struct {
	Name          string
	Address       string
	LoggerManager tlelogger.Manager
	Tracer        tletracer.Tracer
	RabbitMQ      *RabbitMQ
}

// NewServer ...
func NewServer(log tlelogger.Manager, tracer tletracer.Tracer, rabbit *RabbitMQ, serviceName string) Server {
	return Server{
		Name:          serviceName,
		RabbitMQ:      rabbit,
		LoggerManager: log,
		Tracer:        tracer,
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
