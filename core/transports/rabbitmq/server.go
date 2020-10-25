package rabbitmq

import (
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server ...
type Server struct {
	Logger   *tlelogger.Manager
	Tracer   tletracer.Tracer
	RabbitMQ *RabbitMQ
}

// NewServer ...
func NewServer(logger *tlelogger.Manager, tracer tletracer.Tracer, rabbit *RabbitMQ) Server {
	return Server{
		RabbitMQ: rabbit,
		Logger:   logger,
		Tracer:   tracer,
	}
}

// Run will ...
func (server *Server) Run(consumers *[]Consumer) error {
	defer server.RabbitMQ.CloseConnection()

	server.RabbitMQ.OpenConnection()
	//consumers := make(map[string]chan, 1)

	for _, endpoint := range *consumers {
		_, err := server.RabbitMQ.Consume(&endpoint)

		if err == nil {
			//consumers[endpoint.Consumer] = ch
		}
	}

	return nil
}
