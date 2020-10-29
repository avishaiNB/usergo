package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	tlectxrabbit "github.com/thelotter-enterprise/usergo/core/context/transport/rabbitmq"
	"github.com/thelotter-enterprise/usergo/core/errors"
)

// Publisher is used to publish messages to rabbit
type Publisher interface {
	PublishEndpoint(context.Context, *Message, string, amqptransport.EncodeRequestFunc, ...amqptransport.PublisherOption) (endpoint.Endpoint, error)
	Close(context.Context) error
	Publish(context.Context, *Message, string, amqptransport.EncodeRequestFunc) error
}

type publisher struct {
	connectionManager *ConnectionManager
	ch                *amqp.Channel
	isConnected       bool
	topology          Topology
}

// NewPublisher will create a new publisher and will establish a connection to rabbit
func NewPublisher(conn *ConnectionManager) Publisher {
	p := publisher{
		connectionManager: conn,
		topology:          NewTopology(),
	}
	p.connect()
	return &p
}

// Publish will publish a message into the requested exchange
// if the exchange do not exist it will create it
func (p publisher) Publish(ctx context.Context, message *Message, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) error {
	if p.isConnected == false {
		return errors.NewApplicationErrorf("before publishing, you must connect to rabbitMQ")
	}

	p.buildExchange(exchangeName)
	ep, _ := p.PublishEndpoint(ctx, message, exchangeName, encodeFunc)
	_, err := ep(ctx, message)

	return err
}

func (p publisher) PublishEndpoint(ctx context.Context, message *Message, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc, options ...amqptransport.PublisherOption) (endpoint.Endpoint, error) {
	duration, _ := tlectx.GetTimeout(ctx)

	// building the publisher options
	before := amqptransport.PublisherBefore(
		// must come first since it creates the transport context!
		tlectxrabbit.WriteMessageRequestFunc(),
		amqptransport.SetPublishDeliveryMode(2),
		amqptransport.SetPublishExchange(exchangeName))

	ops := make([]amqptransport.PublisherOption, 0)
	ops = append(ops, options...)
	ops = append(ops, amqptransport.PublisherTimeout(duration), amqptransport.PublisherDeliverer(amqptransport.SendAndForgetDeliverer))
	ops = append(ops, before)

	publisher := amqptransport.NewPublisher(p.ch, &amqp.Queue{Name: ""}, encodeFunc, NoopResponseDecoder, ops...)

	return publisher.Endpoint(), nil
}

func (p publisher) buildExchange(exchanegName string) {
	conn := *p.connectionManager
	ch, err := conn.GetChannel()
	if err == nil {
		defer ch.Close()
		p.topology.BuildDurableExchange(ch, exchanegName)
	}
}

// NoopResponseDecoder is a no operation needed
// Used for One way messages
func NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		message := request.(*Message)
		body, err := json.Marshal(message)
		p.Body = body
		return err
	}
	return f
}

// Close will shutdown the client gracely
func (p *publisher) Close(ctx context.Context) error {
	var err error

	if p.isConnected && p.ch != nil {
		cherr := p.ch.Close()

		if cherr != nil {
			err = errors.NewApplicationError(err, cherr.Error())
		} else {
			p.isConnected = false
		}
	}

	return err
}

func (p *publisher) connect() error {
	connMgr := *p.connectionManager
	ch, err := connMgr.GetChannel()
	if err == nil {
		p.ch = ch
		//p.changeConnection(ctx, conn, ch)
		p.isConnected = true
	}
	return err
}
