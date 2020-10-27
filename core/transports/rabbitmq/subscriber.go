package rabbitmq

import (
	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectxamqp "github.com/thelotter-enterprise/usergo/core/context/transport/amqp"
)

// Subscriber ...
type Subscriber struct {
	Sub                *amqpkit.Subscriber
	Channel            *amqp.Channel
	Consumers          map[string]Consumer
	QueueName          string
	QueueDurable       bool
	QueueAutoDelete    bool
	ExchangeDurable    bool
	ExchangeAutoDelete bool
	ExchangeName       string
	SubscriberName     string
	AutoAck            bool
	Exclusive          bool
	NoLocal            bool
	NoWail             bool
	Args               amqp.Table
}

// NewSubscriber will create a new rabbitMQ consumer
func NewSubscriber(
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options []amqptransport.SubscriberOption,
	consumers ...Consumer,
) Subscriber {

	sub := newSubscriber(endpoint, exchangeName, dec, enc, options)
	s := Subscriber{
		Sub:                sub,
		QueueName:          queueName,
		ExchangeName:       exchangeName,
		SubscriberName:     subscriberName,
		Args:               nil,
		Exclusive:          true,
		AutoAck:            true,
		NoLocal:            false,
		NoWail:             false,
		ExchangeAutoDelete: false,
		ExchangeDurable:    true,
		QueueAutoDelete:    false,
		QueueDurable:       true,
		Consumers:          make(map[string]Consumer),
	}

	s.registerConsumers(consumers...)

	return s
}

// NewSubscriber ...
func newSubscriber(
	endpoint endpoint.Endpoint,
	exchangeName string,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options []amqptransport.SubscriberOption) *amqptransport.Subscriber {

	ops := make([]amqpkit.SubscriberOption, 0)
	ops = append(ops, options...)
	ops = append(ops, amqptransport.SubscriberResponsePublisher(amqptransport.NopResponsePublisher))
	ops = append(ops, amqptransport.SubscriberErrorEncoder(amqptransport.ReplyErrorEncoder))
	ops = append(
		ops,
		amqptransport.SubscriberBefore(
			amqptransport.SetPublishExchange(exchangeName),
			tlectxamqp.ReadMessageRequestFunc(),
			amqptransport.SetPublishDeliveryMode(2),
		))

	sub := amqptransport.NewSubscriber(endpoint, dec, enc, ops...)

	return sub
}

// registerConsumer will register the consumers
// If two consumers register targeting the same exchange, an error will be raised
func (s *Subscriber) registerConsumers(consumers ...Consumer) error {
	for _, consumer := range consumers {
		s.Consumers[consumer.Name()] = consumer
	}
	return nil
}
