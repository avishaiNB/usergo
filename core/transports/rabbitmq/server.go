package rabbitmq

import (
	"context"
	"fmt"

	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server ...
type Server struct {
	Logger   *tlelogger.Manager
	Tracer   tletracer.Tracer
	RabbitMQ *Client
}

// NewServer ...
func NewServer(logger *tlelogger.Manager, tracer tletracer.Tracer, rabbit *Client) Server {
	return Server{
		RabbitMQ: rabbit,
		Logger:   logger,
		Tracer:   tracer,
	}
}

// Run will start all the listening on all the consumers
func (server *Server) Run(ctx context.Context, consumers *[]Subscriber) error {
	// cleaning up
	defer server.close(ctx, consumers)

	if err := server.open(ctx); err != nil {
		return err
	}

	forever := make(chan bool)

	server.consume(ctx, consumers)
	<-forever

	return nil
}

func (server *Server) consume(ctx context.Context, consumers *[]Subscriber) {
	logger := *server.Logger

	for _, consumer := range *consumers {
		messages, err := server.RabbitMQ.Consume(&consumer)

		if err != nil {
			msg := fmt.Sprintf("failed to consume %s", consumer.SubscriberName)
			logger.Error(ctx, msg)
		}

		if err == nil {
			go func() {
				for d := range messages {
					// TODO: how to consume the messages?
					logger.Debug(ctx, "Received a message: %s", d.Body)
				}
			}()
		}
	}
}

func (server *Server) open(ctx context.Context) error {
	logger := *server.Logger

	logger.Debug(ctx, "opening rabbitmq connection")
	_, err := server.RabbitMQ.OpenConnection()

	if err != nil {
		msg := "failed to open rabbitmq connection"
		logger.Error(ctx, msg)
		return tleerrors.NewApplicationError(err, msg)
	}
	return nil
}

func (server *Server) close(ctx context.Context, consumers *[]Subscriber) {
	logger := *server.Logger
	logger.Debug(ctx, "closing rabbitmq connection")
	server.RabbitMQ.CloseConnection()

	for _, consumer := range *consumers {
		if consumer.Channel != nil {
			err := consumer.Channel.Close()
			if err != nil {
				msg := fmt.Sprintf("failed to close channel on consumer %s", consumer.SubscriberName)
				logger.Error(ctx, msg)
			}
		}
	}
}
