package rabbitmq

import (
	"context"
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
	ConnectionManager     *ConnectionManager
	Channel               *amqp.Channel
	IsConnected           bool
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
	connMgr *ConnectionManager,
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption,
) Subscriber {

	queueName = fmt.Sprintf("%s-private-%s", queueName, utils.NewUUID())
	topology := NewTopology()
	sub := newSubscriber(endpoint, exchangeName, dec, enc, options...)
	s := Subscriber{
		ConnectionManager:     connMgr,
		Sub:                   sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		SubscriberName:        subscriberName,
		BuildQueueTopology:    topology.BuildNonDurableQueue,
		BuildExchangeTopology: topology.BuildNonDurableExchange,
		BindQueueTopology:     topology.QueueBind,
		ConsumeTopology:       topology.Consume,
		QosTopology:           topology.Qos,
	}

	return s
}

// NewCommandSubscriber will create a new rabbitMQ consumer
func NewCommandSubscriber(
	connMgr *ConnectionManager,
	subscriberName string,
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption,
) Subscriber {

	queueName = queueName + "-command"
	topology := NewTopology()
	sub := newSubscriber(endpoint, exchangeName, dec, enc, options...)
	s := Subscriber{
		ConnectionManager:     connMgr,
		Sub:                   sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		SubscriberName:        subscriberName,
		BuildQueueTopology:    topology.BuildDurableQueue,
		BuildExchangeTopology: topology.BuildDurableExchange,
		BindQueueTopology:     topology.QueueBind,
		ConsumeTopology:       topology.Consume,
		QosTopology:           topology.Qos,
	}

	return s
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
			tlectxamqp.ReadMessageRequestFunc(),
			amqptransport.SetPublishDeliveryMode(2),
		))

	sub := amqptransport.NewSubscriber(endpoint, dec, enc, ops...)

	return sub
}

// Consume ...
func (sub *Subscriber) Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	sub.Channel = ch
	sub.QosTopology(sub.Channel)
	sub.BuildQueueTopology(sub.Channel, sub.QueueName)
	sub.BuildExchangeTopology(sub.Channel, sub.ExchangeName)
	sub.BindQueueTopology(sub.Channel, sub.QueueName, sub.ExchangeName)
	c, err := sub.ConsumeTopology(sub.Channel, sub.QueueName)

	return c, err
}

// Close will shutdown the client gracely
func (sub *Subscriber) Close(ctx context.Context) error {
	conn := *sub.ConnectionManager
	err := conn.CloseChannel(ctx, sub.Channel)
	if err == nil {
		sub.IsConnected = false
	}
	return err
}

func (sub *Subscriber) connect() error {
	conn := *sub.ConnectionManager
	ch, err := conn.GetChannel()
	if err == nil {
		sub.Channel = ch
		sub.IsConnected = true
	}
	//p.changeConnection(ctx, conn, ch)
	return err
}
