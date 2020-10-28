//Package rabbitmq2 provides utilities for working with rabbitmq message broker.
package rabbitmq2

import (
	"context"
)

// Client ..
type Client interface {
	RegisterPrivateConsumers(context.Context, ...Consumer) error
	RegisterCommandConsumers(context.Context, ...Consumer) error
	Publish(context.Context, string, interface{}) error
	Close(context.Context) error
	Run(context.Context) error
}

type client struct {
	config *Config

	publisher *publisher

	// privateSubscriber will listening to queue like xxx-uuid
	privateSubscriber *subscriber

	// commandSubscriber will be listening to queue like xxx-command
	commandSubscriber *subscriber
}

// NewClient initializes RabbitMQ publisher and subscriber.
// Example of connection string: 'amqp://guest:guest@tle-rabbitmq-headless:5672/thelotter'
func NewClient(serviceName string, connectionInfo ConnectionInfo) (Client, error) {
	config := NewConfig(connectionInfo.URL)

	commandSubscriber, err := newCommandSubscriber(config, serviceName)
	if err != nil {
		return nil, err
	}

	privateSubscriber, err := newPrivateSubscriber(config, serviceName)
	if err != nil {
		return nil, err
	}

	publisher := newPublisher(config)
	return &client{
		config:            config,
		commandSubscriber: commandSubscriber,
		privateSubscriber: privateSubscriber,
		publisher:         publisher,
	}, nil
}

// RegisterPrivateConsumers will initialize and run the consumers
func (c client) RegisterPrivateConsumers(ctx context.Context, consumers ...Consumer) error {
	err := c.privateSubscriber.registerConsumer(consumers...)
	if err != nil {
		//c.config.Logger.Error(ctx, err.Error())
		return err
	}
	// c.config.Logger.Info(ctx, "Private consumers registered", consumers)
	return nil
}

// RegisterCommandConsumers will initialize and run the consumers
func (c client) RegisterCommandConsumers(ctx context.Context, consumers ...Consumer) error {
	err := c.commandSubscriber.registerConsumer(consumers...)
	if err != nil {
		//c.config.Logger.Error(ctx, err.Error())
		return err
	}
	// c.config.Logger.Info(ctx, "Command consumers registered", consumers)
	return nil
}

// Publish will use the publisher to publish the message to the exchange
func (c client) Publish(ctx context.Context, exchangeName string, data interface{}) error {
	return c.publisher.Publish(ctx, exchangeName, data)
}

// Run will start all the consumers to consume incomming messages
func (c client) Run(ctx context.Context) error {
	err := c.commandSubscriber.Consume(ctx)
	if err != nil {
		return err
	}
	err = c.privateSubscriber.Consume(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Close will shutdown the client gracely
func (c client) Close(ctx context.Context) error {
	err := c.commandSubscriber.Close()
	if err != nil {
		//c.config.Logger.Error(ctx, "Error while closing command subscriber", err)
		return err
	}
	err = c.privateSubscriber.Close()
	if err != nil {
		//c.config.Logger.Error(ctx, "Error while closing private subscriber", err)
		return err
	}
	err = c.publisher.Close()
	if err != nil {
		//c.config.Logger.Error(ctx, "Error while closing command publisher", err)
		return err
	}
	return nil
}
