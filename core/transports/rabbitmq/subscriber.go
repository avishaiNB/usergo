package rabbitmq

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectxamqp "github.com/thelotter-enterprise/usergo/core/context/transport/amqp"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// Subscriber ...
type Subscriber struct {
	Sub                   *amqpkit.Subscriber
	Channel               *amqp.Channel
	Consumers             map[string]Consumer
	ExchangeName          string
	QueueName             string
	SubscriberName        string
	BuildQueueTopology    func(channel *amqp.Channel, queueName string) (amqp.Queue, error)
	BuildExchangeTopology func(*amqp.Channel, string) error
	BindQueueTopology     func(*amqp.Channel, string, string) error
	ConsumeTopology       func(*amqp.Channel, string) (<-chan amqp.Delivery, error)
	QosTopology           func(ch *amqp.Channel) error
	// TBD: IsConnected
}

// NewPrivateSubscriber will create a new rabbitMQ consumer
func NewPrivateSubscriber(
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options []amqptransport.SubscriberOption,
	consumers ...Consumer,
) Subscriber {

	queueName = fmt.Sprintf("%s-private-%s", queueName, utils.NewUUID())
	topology := NewTopology()
	sub := newSubscriber(endpoint, exchangeName, dec, enc, options)
	s := Subscriber{
		Sub:                   sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		SubscriberName:        subscriberName,
		Consumers:             make(map[string]Consumer),
		BuildQueueTopology:    topology.BuildNonDurableQueue,
		BuildExchangeTopology: topology.BuildNonDurableExchange,
		BindQueueTopology:     topology.QueueBind,
		ConsumeTopology:       topology.Consume,
		QosTopology:           topology.Qos,
	}

	s.registerConsumers(consumers...)

	return s
}

// NewCommandSubscriber will create a new rabbitMQ consumer
func NewCommandSubscriber(
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options []amqptransport.SubscriberOption,
	consumers ...Consumer,
) Subscriber {

	queueName = queueName + "-command"
	topology := NewTopology()
	sub := newSubscriber(endpoint, exchangeName, dec, enc, options)
	s := Subscriber{
		Sub:                   sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		SubscriberName:        subscriberName,
		Consumers:             make(map[string]Consumer),
		BuildQueueTopology:    topology.BuildDurableQueue,
		BuildExchangeTopology: topology.BuildDurableExchange,
		BindQueueTopology:     topology.QueueBind,
		ConsumeTopology:       topology.Consume,
		QosTopology:           topology.Qos,
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
