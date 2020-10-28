package rabbitmq

import (
	"context"
	"fmt"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server responsibility to initiate the ability to consume messages
type Server interface {
	Run(context.Context) error
	Shutdown(context.Context)
}

type server struct {
	logger            *tlelogger.Manager
	tracer            tletracer.Tracer
	connectionManager *ConnectionManager
	client            *Client
}

// NewServer will create a new instance of Server which can be executed to start and recieving messages
func NewServer(logger *tlelogger.Manager, tracer tletracer.Tracer, rabbit *Client, conn *ConnectionManager) Server {
	return &server{
		client:            rabbit,
		logger:            logger,
		tracer:            tracer,
		connectionManager: conn,
	}
}

// Run will start all the listening on all the consumers
func (s server) Run(ctx context.Context) error {
	// cleaning up
	defer s.Shutdown(ctx)

	forever := make(chan bool)
	c := *s.client
	c.Consume(ctx)
	<-forever

	return nil
}

// Shutdown will close the server and call client to close resources
func (s server) Shutdown(ctx context.Context) {
	c := *s.client
	c.Close(ctx)
	fmt.Print("Shutdown amqp server")
}
