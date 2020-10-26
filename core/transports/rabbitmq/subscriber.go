package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// Subscriber ...
type Subscriber struct {
	Sub                *amqpkit.Subscriber
	Channel            *amqp.Channel
	QueueName          string
	QueueDurable       bool
	QueueAutoDelete    bool
	ExchangeDurable    bool
	ExchangeAutoDelete bool
	ExchangeName       string
	ConsumerName       string
	AutoAck            bool
	Exclusive          bool
	NoLocal            bool
	NoWail             bool
	Args               amqp.Table
}

// NewSubscriber will create a new rabbitMQ consumer
func NewSubscriber(
	consumerName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption) Subscriber {

	sub := newSubscriber(endpoint, exchangeName, dec, enc, options...)
	consumer := Subscriber{
		Sub:                sub,
		QueueName:          queueName,
		ExchangeName:       exchangeName,
		ConsumerName:       consumerName,
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

// NewSubscriber ...
func newSubscriber(
	endpoint endpoint.Endpoint,
	exchangeName string,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption) *amqptransport.Subscriber {

	ops := make([]amqpkit.SubscriberOption, 0)
	ops = append(ops, options...)
	ops = append(ops, amqptransport.SubscriberResponsePublisher(amqptransport.NopResponsePublisher))
	ops = append(ops, amqptransport.SubscriberErrorEncoder(amqptransport.ReplyErrorEncoder))
	ops = append(
		ops,
		amqptransport.SubscriberBefore(
			amqptransport.SetPublishExchange(exchangeName),
			readMessageIntoContext(),
			amqptransport.SetPublishDeliveryMode(2),
		))

	sub := amqptransport.NewSubscriber(endpoint, dec, enc, ops...)

	return sub
}

// TODO: need to read into the context the correaltion ID and etc.
func readMessageIntoContext() amqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		return ctx
	}
}
