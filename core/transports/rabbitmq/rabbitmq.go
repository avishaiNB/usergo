package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
)

// RabbitMQ contains data required to make a connection to the rabbitMQ instance
type RabbitMQ struct {
	ConnectionMeta ConnectionInfo

	// AMQPConnection to rabbitMQ. Will be nil until Connect will be called
	AMQPConnection *amqp.Connection

	LogManager *tlelogger.Manager
}

// NewRabbitMQ will create a new instance of empty RabbitMQ
func NewRabbitMQ(logManager *tlelogger.Manager, connection ConnectionInfo) RabbitMQ {
	return RabbitMQ{
		ConnectionMeta: connection,
		LogManager:     logManager,
	}
}

// OpenConnection will create a new connection to RabbitMQ based on the input entered when created the RabbitMQ instance
// Connection will be returned BUT also stored in the RabbitMQ instance
func (rabbit *RabbitMQ) OpenConnection() (*amqp.Connection, error) {
	if rabbit.AMQPConnection != nil {
		return rabbit.AMQPConnection, nil
	}
	conn, err := amqp.Dial(rabbit.ConnectionMeta.URL)
	if err == nil {
		rabbit.AMQPConnection = conn
	}
	return conn, err
}

// CloseConnection will close the open connection attached to the RabbitMQ instance
func (rabbit *RabbitMQ) CloseConnection() error {
	var err error
	if rabbit.AMQPConnection != nil {
		err = rabbit.AMQPConnection.Close()
	}
	return err
}

// Consume ...
func (rabbit *RabbitMQ) Consume(consumer *Consumer) (<-chan amqp.Delivery, error) {
	rabbit.newConsumerChannel(consumer)
	consumer.newExchange(consumer.ExchangeName, consumer.ExchangeDurable, consumer.ExchangeAutoDelete)
	consumer.newQueue(consumer.QueueName, consumer.QueueDurable, consumer.QueueAutoDelete)
	consumer.bind(consumer.QueueName, consumer.ExchangeName)

	c, err := consumer.Channel.Consume(
		consumer.QueueName,
		consumer.ConsumerName,
		consumer.AutoAck,
		consumer.Exclusive,
		consumer.NoLocal,
		consumer.NoWail,
		consumer.Args)

	return c, err
}

// PublishOneWay will 'send and forget' a message to the given exchange
func (rabbit *RabbitMQ) PublishOneWay(ctx context.Context, request interface{}, tgtExchangeName string, encodeFunc amqptransport.EncodeRequestFunc) error {
	e := rabbit.oneWayPublisherEndpoint(ctx, tgtExchangeName, encodeFunc)
	_, err := e(ctx, request)
	return err
}

// OneWayPublisherEndpoint will create a 'send and forget' publisher endpoint
func (rabbit *RabbitMQ) oneWayPublisherEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	m := tlectx.NewManager()
	corrid := m.GetCorrelationID(ctx)
	duration, _ := m.GetTimeout(ctx)
	var channel amqptransport.Channel
	var queue *amqp.Queue
	channel, _ = rabbit.NewChannel()
	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		channel,
		queue,
		encodeFunc,
		rabbit.NoopResponseDecoder,
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
func (rabbit *RabbitMQ) NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func (rabbit *RabbitMQ) DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		var err error
		marshall := MessageMarshall{}
		*p, err = marshall.Marshal(ctx, exchangeName, request)
		return err
	}
	return f
}

func (rabbit *RabbitMQ) newConsumerChannel(consumer *Consumer) {
	if consumer.Channel != nil {
		return
	}

	var err error
	var ch *amqp.Channel
	ch, err = rabbit.NewChannel()

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
