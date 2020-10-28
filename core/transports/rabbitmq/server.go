package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server ...
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

// NewServer ...
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
	s.consume(ctx)
	<-forever

	return nil
}

func (s server) consume(ctx context.Context) {
	logger := *s.logger
	conn := *s.connectionManager
	for _, sub := range *s.client.Subscribers {
		ch, err := conn.GetChannel()
		messages, err := sub.Consume(ch)

		if err != nil {
			msg := fmt.Sprintf("failed to consume %s", sub.SubscriberName)
			logger.Error(ctx, msg)
		}

		if err == nil {
			go func() {
				for d := range messages {
					// logger.Debug(ctx, "Received a message: %s", d.Body)
					fmt.Printf("Received a message: %s", d.Body)
					sub.KitSubscriber.ServeDelivery(sub.Channel)(&amqp.Delivery{})
				}
			}()
		}
	}
}

// Shutdown will close the server and call client to close resources
func (s server) Shutdown(ctx context.Context) {
	s.client.Close(ctx)
	fmt.Print("Shutdown amqp server")
}
