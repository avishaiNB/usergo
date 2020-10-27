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
	Logger *tlelogger.Manager
	Tracer tletracer.Tracer
	Client *Client
}

// NewServer ...
func NewServer(logger *tlelogger.Manager, tracer tletracer.Tracer, rabbit *Client) Server {
	return Server{
		Client: rabbit,
		Logger: logger,
		Tracer: tracer,
	}
}

// Run will start all the listening on all the consumers
func (server *Server) Run(ctx context.Context) error {
	// cleaning up
	defer server.close(ctx)

	if err := server.open(ctx); err != nil {
		return err
	}

	forever := make(chan bool)
	server.consume(ctx)
	<-forever

	return nil
}

func (server *Server) consume(ctx context.Context) {
	logger := *server.Logger

	for _, sub := range *server.Client.Subscribers {
		ch, err := server.Client.NewChannel()
		messages, err := sub.Consume(ch)

		if err != nil {
			msg := fmt.Sprintf("failed to consume %s", sub.SubscriberName)
			logger.Error(ctx, msg)
		}

		if err == nil {
			go func() {
				for d := range messages {
					logger.Debug(ctx, "Received a message: %s", d.Body)
					sub.Sub.ServeDelivery(sub.Channel)
				}
			}()
		}
	}
}

func (server *Server) open(ctx context.Context) error {
	logger := *server.Logger

	logger.Debug(ctx, "opening rabbitmq connection")
	_, err := server.Client.OpenConnection()

	if err != nil {
		msg := "failed to open rabbitmq connection"
		logger.Error(ctx, msg)
		return tleerrors.NewApplicationError(err, msg)
	}
	return nil
}

func (server *Server) close(ctx context.Context) {
	logger := *server.Logger
	logger.Debug(ctx, "closing rabbitmq connection")
	server.Client.CloseConnection()

	for _, sub := range *server.Client.Subscribers {
		if sub.Channel != nil {
			err := sub.Channel.Close()
			if err != nil {
				msg := fmt.Sprintf("failed to close channel on consumer %s", sub.SubscriberName)
				logger.Error(ctx, msg)
			}
		}
	}
}
