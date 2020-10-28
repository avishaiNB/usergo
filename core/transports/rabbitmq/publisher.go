package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/errors"
)

// Publisher ...
type Publisher interface {
	PublishOneWay(context.Context, string, amqptransport.EncodeRequestFunc) (endpoint.Endpoint, error)
	Close(context.Context) error
}

type publisher struct {
	connection     *amqp.Connection
	ConnectionInfo ConnectionInfo
	ch             *amqp.Channel
	isConnected    bool
}

// NewPublisher will create a new publisher and will establish a connection to rabbit
func NewPublisher(connInfo ConnectionInfo) Publisher {
	p := publisher{
		ConnectionInfo: connInfo,
	}
	p.connect()
	return &p
}

func (p publisher) PublishOneWay(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) (endpoint.Endpoint, error) {
	if p.isConnected == false {
		return nil, errors.NewApplicationErrorf("before publishing, you must connect to rabbitMQ")
	}

	corrid := tlectx.GetCorrelation(ctx)
	duration, _ := tlectx.GetTimeout(ctx)
	var queue *amqp.Queue

	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		p.ch,
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

	return publisher.Endpoint(), nil
}

// NoopResponseDecoder is a no operation needed
// Used for One way messages
func NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		var err error
		marshall := MessageMarshall{}
		*p, err = marshall.Marshal(ctx, exchangeName, request)
		return err
	}
	return f
}

// Close will shutdown the client gracely
func (p *publisher) Close(ctx context.Context) error {
	var err error

	if p.isConnected {
		connerr := p.connection.Close()

		if connerr != nil {
			err = errors.NewApplicationErrorf("failed to close rabbit connection %s", connerr.Error())
		} else {
			if p.ch != nil {
				cherr := p.ch.Close()

				if cherr != nil {
					err = errors.Annotate(err, cherr.Error())
				}
			}
		}
	}

	if err == nil {
		p.isConnected = false
	}

	return err
}

func (p *publisher) connect() error {
	conn, err := amqp.Dial(p.ConnectionInfo.URL)
	if err != nil {
		return errors.NewApplicationError(err, "failed to connect to rabbit")
	}
	ch, err := conn.Channel()
	if err != nil {
		return errors.NewApplicationError(err, "failed to create channel")
	}
	p.ch = ch
	p.connection = conn
	//p.changeConnection(ctx, conn, ch)
	p.isConnected = true
	return nil
}
