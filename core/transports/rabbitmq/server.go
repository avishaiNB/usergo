package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// Server ...
type Server struct {
	Logger            *tlelogger.Manager
	Tracer            tletracer.Tracer
	ConnectionManager *ConnectionManager
	Client            *Client
}

// NewServer ...
func NewServer(logger *tlelogger.Manager, tracer tletracer.Tracer, rabbit *Client, conn *ConnectionManager) Server {
	return Server{
		Client:            rabbit,
		Logger:            logger,
		Tracer:            tracer,
		ConnectionManager: conn,
	}
}

// Run will start all the listening on all the consumers
func (server *Server) Run(ctx context.Context) error {
	// cleaning up
	defer server.close(ctx)

	forever := make(chan bool)
	server.consume(ctx)
	<-forever

	return nil
}

func (server *Server) consume(ctx context.Context) {
	logger := *server.Logger
	conn := *server.ConnectionManager
	for _, sub := range *server.Client.Subscribers {
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
					sub.Sub.ServeDelivery(sub.Channel)(&amqp.Delivery{})
				}
			}()
		}
	}
}

func (server *Server) close(ctx context.Context) {
	// logger := *server.Logger
	// logger.Debug(ctx, "closing rabbitmq connection")
	// server.Client.CloseConnection()

	// for _, sub := range *server.Client.Subscribers {
	// 	if sub.Channel != nil {
	// 		err := sub.Channel.Close()
	// 		if err != nil {
	// 			msg := fmt.Sprintf("failed to close channel on consumer %s", sub.SubscriberName)
	// 			logger.Error(ctx, msg)
	// 		}
	// 	}
	// }
}
