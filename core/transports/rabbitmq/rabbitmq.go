package rabbitmq

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
	ConnectionMeta ConnectionMeta

	// AMQPConnection to rabbitMQ. Will be nil until Connect will be called
	AMQPConnection *amqp.Connection

	Log core.Log
}

// NewRabbitMQ will create a new instance of empty RabbitMQ
func NewRabbitMQ(log core.Log, connection ConnectionMeta) RabbitMQ {
	return RabbitMQ{
		ConnectionMeta: connection,
		Log:            log,
	}
}

// Consume ...
func (a *RabbitMQ) Consume(consumer *Consumer) (<-chan amqp.Delivery, error) {
	a.setConsumerChannel(consumer)
	consumer.newExchange(consumer.ExchangeName, consumer.ExchangeDurable, consumer.ExchangeAutoDelete)
	consumer.newQueue(consumer.QueueName, consumer.QueueDurable, consumer.QueueAutoDelete)
	consumer.bind(consumer.QueueName, consumer.ExchangeName)

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

// OpenConnection will create a new connection to RabbitMQ based on the input entered when created the RabbitMQ instance
// Connection will be returned BUT also stored in the RabbitMQ instance
func (a *RabbitMQ) OpenConnection() (*amqp.Connection, error) {
	if a.AMQPConnection != nil {
		return a.AMQPConnection, nil
	}
	conn, err := amqp.Dial(a.ConnectionMeta.URL)
	if err == nil {
		a.AMQPConnection = conn
	}
	return conn, err
}

// CloseConnection will close the open connection attached to the RabbitMQ instance
func (a *RabbitMQ) CloseConnection() error {
	var err error
	if a.AMQPConnection != nil {
		err = a.AMQPConnection.Close()
	}
	return err
}

// OneWayPublisherEndpoint will create a 'send and forget' publisher endpoint
func (a *RabbitMQ) OneWayPublisherEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	corrid := tlectx.GetCorrelationFromContext(ctx)
	duration, _ := tlectx.GetTimeoutFromContext(ctx)
	var channel amqptransport.Channel
	var queue *amqp.Queue
	_, _ = a.OpenConnection()
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

// setConsumerChannel ...
func (a *RabbitMQ) setConsumerChannel(consumer *Consumer) {
	var err error
	var ch *amqp.Channel
	ch, err = a.NewChannel()

	if err == nil {
		consumer.Channel = ch
	}
}

// newQueue will create a new queue
func (c *Consumer) newQueue(name string, durable bool, autoDelete bool) (amqp.Queue, error) {
	var err error
	var queue amqp.Queue

	queue, err = c.Channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return queue, err
}

func (c *Consumer) newExchange(name string, durable bool, autoDelete bool) error {

	err := c.Channel.ExchangeDeclare(
		name,       // name
		"fanout",   // type
		durable,    // durable
		autoDelete, // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)

	return err
}

// bind will bind the rabbitMQ queue and exchange together
func (c *Consumer) bind(queueName string, exchangeName string) error {
	err := c.Channel.QueueBind(
		queueName,
		"", // bindingKey
		exchangeName,
		false, // noWait
		nil,   // arguments
	)

	return err
}
