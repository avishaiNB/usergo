package amqp

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// RabbitMQConsumer ...RabbitMQConsumer
type RabbitMQConsumer struct {
	Sub                *amqpkit.Subscriber
	Channel            *amqp.Channel
	QueueName          string
	QueueDurable       bool
	QueueAutoDelete    bool
	ExchangeDurable    bool
	ExchangeAutoDelete bool
	ExchangeName       string
	Consumer           string
	AutoAck            bool
	Exclusive          bool
	NoLocal            bool
	NoWail             bool
	Args               amqp.Table
}

// NewConsumer will create a new rabbitMQ consumer
func (a *RabbitMQ) NewConsumer(
	name string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc) RabbitMQConsumer {

	sub := newSubscriber(endpoint, exchangeName, dec)
	consumer := RabbitMQConsumer{
		Sub:                sub,
		QueueName:          queueName,
		ExchangeName:       exchangeName,
		Consumer:           name,
		Args:               nil,
		Exclusive:          true,
		AutoAck:            true,
		NoLocal:            false,
		NoWail:             false,
		ExchangeAutoDelete: false,
		ExchangeDurable:    true,
		QueueAutoDelete:    false,
		QueueDurable:       true,
	}

	return consumer
}

// SetConsumerChannel ...
func (a *RabbitMQ) SetConsumerChannel(consumer *RabbitMQConsumer) {
	var err error
	var ch *amqp.Channel
	ch, err = a.NewChannel()

	if err == nil {
		consumer.Channel = ch
	}
}

// NewQueue will create a new queue
func (c *RabbitMQConsumer) NewQueue(name string, durable bool, autoDelete bool) (amqp.Queue, error) {
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

// NewExchange will create a new exchange
func (c *RabbitMQConsumer) NewExchange(name string, durable bool, autoDelete bool) error {

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

// Bind will bind the rabbitMQ queue and exchange together
func (c *RabbitMQConsumer) Bind(queueName string, exchangeName string) error {

	err := c.Channel.QueueBind(
		queueName,
		"", // bindingKey
		exchangeName,
		false, // noWait
		nil,   // arguments
	)

	return err
}

// NewSubscriber ...
func newSubscriber(endpoint endpoint.Endpoint, exchangeName string, dec amqptransport.DecodeRequestFunc) *amqptransport.Subscriber {
	sub := amqptransport.NewSubscriber(
		endpoint,
		dec,
		amqptransport.EncodeJSONResponse,
		amqptransport.SubscriberResponsePublisher(amqptransport.NopResponsePublisher),
		amqptransport.SubscriberErrorEncoder(amqptransport.ReplyErrorEncoder),
		amqptransport.SubscriberBefore(
			amqptransport.SetPublishExchange(exchangeName),
			readMessageIntoContext(),
			//amqptransport.SetPublishKey(key),
			amqptransport.SetPublishDeliveryMode(2),
		),
	)

	return sub
}

// TODO: need to read into the context the correaltion ID and etc.
func readMessageIntoContext() amqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		return ctx
	}
}
