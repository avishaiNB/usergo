package amqp

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"github.com/thelotter-enterprise/usergo/core"
	tlectx "github.com/thelotter-enterprise/usergo/core/ctx"
)

// RabbitMQ contains data required to make a connection to the rabbitMQ instance
type RabbitMQ struct {
	Connection Connection

	// AMQPConnection to rabbitMQ. Will be nil until Connect will be called
	AMQPConnection *amqp.Connection

	Log core.Log
}

// NewRabbitMQ will create a new instance of empty RabbitMQ
func NewRabbitMQ(log core.Log, connection Connection) RabbitMQ {
	return RabbitMQ{
		Connection: connection,
		Log:        log,
	}
}

// Consume ...
func (a *RabbitMQ) Consume(consumer *RabbitMQConsumer) (<-chan amqp.Delivery, error) {
	a.SetConsumerChannel(consumer)
	consumer.NewExchange(consumer.ExchangeName, consumer.ExchangeDurable, consumer.ExchangeAutoDelete)
	consumer.NewQueue(consumer.QueueName, consumer.QueueDurable, consumer.QueueAutoDelete)
	consumer.Bind(consumer.QueueName, consumer.ExchangeName)

	c, err := consumer.Channel.Consume(
		consumer.QueueName,
		consumer.Consumer,
		consumer.AutoAck,
		consumer.Exclusive,
		consumer.NoLocal,
		consumer.NoWail,
		consumer.Args)

	return c, err
}

// OneWayPublisherEndpoint will create a 'send and forget' publisher endpoint
func (a *RabbitMQ) OneWayPublisherEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	corrid := tlectx.GetCorrelationFromContext(ctx)
	duration, _ := tlectx.GetTimeoutFromContext(ctx)
	var channel amqptransport.Channel
	var queue *amqp.Queue
	_, _ = a.Connect()
	channel, _ = a.NewChannel()
	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		channel,
		queue,
		encodeFunc,
		a.NoopResponseDecoder,
		amqptransport.PublisherBefore(
			amqptransport.SetCorrelationID(corrid),
			amqptransport.SetPublishDeliveryMode(2), // queue implementation use - non-persistent (1) or persistent (2)
			amqptransport.SetPublishExchange(exchangeName)),
		amqptransport.PublisherTimeout(duration),
		amqptransport.PublisherDeliverer(amqptransport.SendAndForgetDeliverer),
	)

	return publisher.Endpoint()
}

// NoopResponseDecoder is a no operation needed
// Used for One way messages
func (a *RabbitMQ) NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func (a *RabbitMQ) DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		var err error
		marshall := MessageMarshall{}
		*p, err = marshall.Marshal(ctx, exchangeName, request)
		return err
	}
	return f
}
