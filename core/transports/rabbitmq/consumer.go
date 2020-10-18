package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// Consumer ...Consumer
type Consumer struct {
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
func NewConsumer(
	name string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc) Consumer {

	sub := newSubscriber(endpoint, exchangeName, dec)
	consumer := Consumer{
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
