package rabbitmq

import (
	"context"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/thelotter-enterprise/usergo/core/errors"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
)

// Client contains data required to make a connection to the rabbitMQ instance
type Client struct {
	ConnectionManager *ConnectionManager

	LogManager *tlelogger.Manager

	Subscribers *[]Subscriber

	Publisher *Publisher
}

// NewClient will create a new instance of empty RabbitMQ
func NewClient(connMgr *ConnectionManager, logManager *tlelogger.Manager, publisher *Publisher, subscribers *[]Subscriber) *Client {
	return &Client{
		LogManager:        logManager,
		Subscribers:       subscribers,
		Publisher:         publisher,
		ConnectionManager: connMgr,
	}
}

// Publish ...
func (c *Client) Publish(ctx context.Context, message *Message, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) error {
	p := *c.Publisher
	ep, _ := p.PublishOneWay(ctx, exchangeName, encodeFunc)
	_, err := ep(ctx, message)

	return err
}

// Close will close the open connection attached to the RabbitMQ instance
func (c *Client) Close(ctx context.Context) error {
	var err error

	// closing the publisher channel
	p := *c.Publisher
	perr := p.Close(ctx)
	if perr != nil {
		err = errors.NewApplicationError(perr, "failed to close rabbit publisher")
	}

	// closing the subscribers channels
	subs := *c.Subscribers
	if subs != nil && len(subs) > 0 {
		for _, sub := range subs {
			suberr := sub.Close(ctx)
			if suberr != nil {
				err = errors.Annotate(err, suberr.Error())
			}
		}
	}

	// closing the connection
	conn := *c.ConnectionManager
	cerr := conn.CloseConnection(ctx)
	if cerr != nil {
		err = errors.Annotate(err, cerr.Error())
	}
	return err
}
