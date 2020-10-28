package rabbitmq

import (
	"context"
	"fmt"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/thelotter-enterprise/usergo/core/errors"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
)

// Client is a rabbit contract to publish and consume messages
type Client interface {
	Consume(context.Context)
	Publish(context.Context, *Message, string, amqptransport.EncodeRequestFunc) error
	Close(context.Context) error
}

type client struct {
	connectionManager *ConnectionManager
	logger            *tlelogger.Manager
	subscribers       *[]Subscriber
	publisher         *Publisher
}

// NewClient will create a new instance of a client
// Best practice is to have a single one per application and reuse it
func NewClient(connMgr *ConnectionManager, logManager *tlelogger.Manager, publisher *Publisher, subscribers *[]Subscriber) Client {
	return &client{
		logger:            logManager,
		subscribers:       subscribers,
		publisher:         publisher,
		connectionManager: connMgr,
	}
}

// Consume will start consuming messages from all the subscribers
// For each subscriber it will create a new go routine and will wait on it for incoming messages
func (c *client) Consume(ctx context.Context) {
	conn := *c.connectionManager
	for _, sub := range *c.subscribers {
		ch, err := conn.GetChannel()
		messages, err := sub.Consume(ch)

		if err != nil {
			// msg := fmt.Sprintf("failed to consume %s", sub.SubscriberName)
			// logger.Error(ctx, msg)
		}

		if err == nil {
			go func() {
				for msg := range messages {
					// logger.Debug(ctx, "Received a message: %s", d.Body)
					fmt.Printf("Received raw message: %s", msg.Body)
					sub.KitSubscriber.ServeDelivery(sub.Channel)(&msg)
				}
			}()
		}
	}
}

// Publish will publish a message into the requested exchange
// if the exchange do not exist it will create it
func (c *client) Publish(ctx context.Context, message *Message, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) error {
	p := *c.publisher
	ep, _ := p.PublishEndpoint(ctx, exchangeName, encodeFunc)
	_, err := ep(ctx, message)

	return err
}

// Close will close the open connections and channels
// This must be called before the application terminate to prevent connection or channel leaks
func (c *client) Close(ctx context.Context) error {
	var err error

	// closing the publisher channel
	p := *c.publisher
	perr := p.Close(ctx)
	if perr != nil {
		err = errors.NewApplicationError(perr, "failed to close rabbit publisher")
	}

	// closing the subscribers channels
	subs := *c.subscribers
	if subs != nil && len(subs) > 0 {
		for _, sub := range subs {
			suberr := sub.Close(ctx)
			if suberr != nil {
				err = errors.Annotate(err, suberr.Error())
			}
		}
	}

	// closing the connection
	conn := *c.connectionManager
	cerr := conn.CloseConnection(ctx)
	if cerr != nil {
		err = errors.Annotate(err, cerr.Error())
	}
	return err
}
