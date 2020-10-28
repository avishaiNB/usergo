package rabbitmq

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	tlectxamqp "github.com/thelotter-enterprise/usergo/core/context/transport/amqp"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// SubscribeManager ...
type SubscribeManager interface {
	NewPrivateSubscriber(subscriberName string, exchangeName string, queueName string, endpoint endpoint.Endpoint, dec amqptransport.DecodeRequestFunc, enc amqptransport.EncodeResponseFunc, options ...amqptransport.SubscriberOption) Subscriber
	NewCommandSubscriber(subscriberName string, exchangeName string, queueName string, endpoint endpoint.Endpoint, dec amqptransport.DecodeRequestFunc, enc amqptransport.EncodeResponseFunc, options ...amqptransport.SubscriberOption) Subscriber
}

type submgr struct {
	connMgr  *ConnectionManager
	topology Topology
}

// NewSubscriberManager ...
func NewSubscriberManager(connMgr *ConnectionManager) SubscribeManager {
	s := submgr{
		connMgr:  connMgr,
		topology: NewTopology(),
	}

	return s
}

// NewCommandSubscriber will create a new rabbitMQ consumer
func (s submgr) NewCommandSubscriber(
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption,
) Subscriber {

	queueName = queueName + "-command"
	sub := newKitSubscriber(endpoint, exchangeName, dec, enc, options...)
	return Subscriber{
		ConnectionManager:     s.connMgr,
		KitSubscriber:         sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		SubscriberName:        subscriberName,
		BuildQueueTopology:    s.topology.BuildDurableQueue,
		BuildExchangeTopology: s.topology.BuildDurableExchange,
		BindQueueTopology:     s.topology.QueueBind,
		ConsumeTopology:       s.topology.Consume,
		QosTopology:           s.topology.Qos,
	}
}

// NewPrivateSubscriber will create a new rabbitMQ consumer
func (s submgr) NewPrivateSubscriber(
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption,
) Subscriber {

	queueName = fmt.Sprintf("%s-private-%s", queueName, utils.NewUUID())
	sub := newKitSubscriber(endpoint, exchangeName, dec, enc, options...)
	return Subscriber{
		ConnectionManager:     s.connMgr,
		KitSubscriber:         sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		SubscriberName:        subscriberName,
		BuildQueueTopology:    s.topology.BuildNonDurableQueue,
		BuildExchangeTopology: s.topology.BuildNonDurableExchange,
		BindQueueTopology:     s.topology.QueueBind,
		ConsumeTopology:       s.topology.Consume,
		QosTopology:           s.topology.Qos,
	}
}

func newKitSubscriber(
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
			tlectxamqp.ReadMessageRequestFunc(),
			amqptransport.SetPublishDeliveryMode(2),
		))

	sub := amqptransport.NewSubscriber(endpoint, dec, enc, ops...)

	return sub
}
