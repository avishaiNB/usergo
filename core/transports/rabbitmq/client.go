package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/errors"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
)

// Client contains data required to make a connection to the rabbitMQ instance
type Client struct {

	// ConnInfo contains the connection information to connect to be able to dail to rabbitMQ
	ConnInfo ConnectionInfo

	// AMQPConnection to rabbitMQ. Will be nil until Connect will be called
	AMQPConnection *amqp.Connection

	// IsConnected indiates if a successful dail into tabbit was already established
	IsConnected bool

	LogManager *tlelogger.Manager

	Subscribers *[]Subscriber
}

// NewClient will create a new instance of empty RabbitMQ
func NewClient(logManager *tlelogger.Manager, connection ConnectionInfo, subscribers *[]Subscriber) *Client {
	return &Client{
		ConnInfo:    connection,
		LogManager:  logManager,
		IsConnected: false,
		Subscribers: subscribers,
	}
}

// OpenConnection will create a new connection to RabbitMQ based on the input entered when created the RabbitMQ instance
// Connection will be returned BUT also stored in the RabbitMQ instance
func (rabbit *Client) OpenConnection() (*amqp.Connection, error) {
	if rabbit.AMQPConnection != nil {
		return rabbit.AMQPConnection, nil
	}
	conn, err := amqp.Dial(rabbit.ConnInfo.URL)
	if err == nil {
		rabbit.AMQPConnection = conn
		rabbit.IsConnected = true
	}
	return conn, err
}

// CloseConnection will close the open connection attached to the RabbitMQ instance
func (rabbit *Client) CloseConnection() error {
	var err error
	if rabbit.AMQPConnection != nil && rabbit.IsConnected {
		err = rabbit.AMQPConnection.Close()
	}
	return err
}

// PublishOneWay will 'send and forget' a message to the given exchange
func (rabbit *Client) PublishOneWay(ctx context.Context, request interface{}, tgtExchangeName string, encodeFunc amqptransport.EncodeRequestFunc) error {
	e := rabbit.oneWayPublisherEndpoint(ctx, tgtExchangeName, encodeFunc)
	_, err := e(ctx, request)
	return err
}

// OneWayPublisherEndpoint will create a 'send and forget' publisher endpoint
func (rabbit *Client) oneWayPublisherEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	corrid := tlectx.GetCorrelation(ctx)
	duration, _ := tlectx.GetTimeout(ctx)
	var channel amqptransport.Channel
	var queue *amqp.Queue
	channel, _ = rabbit.NewChannel()
	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		channel,
		queue,
		encodeFunc,
		rabbit.NoopResponseDecoder,
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
func (rabbit *Client) NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func (rabbit *Client) DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		var err error
		marshall := MessageMarshall{}
		*p, err = marshall.Marshal(ctx, exchangeName, request)
		return err
	}
	return f
}

func (rabbit *Client) newSubscriberChannel(sub *Subscriber) {
	if sub.Channel != nil {
		return
	}

	var err error
	var ch *amqp.Channel
	ch, err = rabbit.NewChannel()

	if err == nil {
		sub.Channel = ch
	}
}

// NewChannel will create a new rabbitMQ channel
func (rabbit *Client) NewChannel() (*amqp.Channel, error) {
	var err error
	var ch *amqp.Channel
	if rabbit.AMQPConnection == nil {
		err = errors.NewApplicationErrorf("Connect to rabbit before tring to get a channel")
	} else {
		ch, err = rabbit.AMQPConnection.Channel()
	}

	return ch, err
}
