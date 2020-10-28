package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
)

// Publisher ...
type Publisher interface {
	PublishOneWay(context.Context, *Message, string, amqptransport.EncodeRequestFunc) endpoint.Endpoint
}

type publisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewPublisher ...
func NewPublisher(connection *amqp.Connection) Publisher {
	ch, _ := connection.Channel()

	return &publisher{
		Connection: connection,
		Channel:    ch,
	}
}

func (p publisher) PublishOneWay(ctx context.Context, message *Message, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	corrid := tlectx.GetCorrelation(ctx)
	duration, _ := tlectx.GetTimeout(ctx)
	var queue *amqp.Queue

	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		p.Channel,
		queue,
		encodeFunc,
		NoopResponseDecoder,
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
func NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}
