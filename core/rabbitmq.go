package core

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// RabbitMQ contains data required to make a connection to the rabbitMQ instance
type RabbitMQ struct {
	// URL like amqp://guest:guest@localhost:5672/
	URL string

	// Usewrname to connect to RabbitMQ
	Username string

	// Pwd to connect to RabbitMQ
	Pwd string

	// VirtualHost to connect to RabbitMQ
	VirtualHost string

	// Port to connect to RabbitMQ
	Port int

	// Host to connect to RabbitMQ
	Host string

	// Connection to rabbitMQ. Will be nil until Connect will be called
	Connection *amqp.Connection

	Log Log
}

// RabbitMQConsumer ...RabbitMQConsumer
type RabbitMQConsumer struct {
	Sub       *amqpkit.Subscriber
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWail    bool
	Args      amqp.Table
}

// NewRabbitMQ will create a new instance of empty RabbitMQ
func NewRabbitMQ(log Log, host string, port int, username string, password string, vhost string) RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
	return RabbitMQ{
		URL:         url,
		Log:         log,
		Host:        host,
		VirtualHost: vhost,
		Pwd:         password,
		Username:    username,
		Port:        port,
	}
}

// Consume ...
func (a *RabbitMQ) Consume(consumer *RabbitMQConsumer) (<-chan amqp.Delivery, error) {
	_, err := a.Connect()
	if err != nil {
		panic(err)
	}

	channel, err := a.Channel()
	if err != nil {
		panic(err)
	}

	c, err := channel.Consume(
		consumer.Queue,
		consumer.Consumer,
		consumer.AutoAck,
		consumer.Exclusive,
		consumer.NoLocal,
		consumer.NoWail,
		consumer.Args)

	return c, err
}

// NewConsumer will create a new rabbitMQ consumer
func (a *RabbitMQ) NewConsumer(
	name string,
	exchangeName string,
	queue string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc) RabbitMQConsumer {

	sub := a.NewSubscriber(endpoint, exchangeName, dec)
	consumer := RabbitMQConsumer{
		Sub:       sub,
		Queue:     queue,
		Consumer:  name,
		Args:      nil,
		Exclusive: true,
		AutoAck:   true,
		NoLocal:   false,
		NoWail:    false,
	}

	return consumer
}

// Connect will create a new connection to RabbitMQ based on the input entered when created the RabbitMQ instance
// Connection will be returned BUT also stored in the RabbitMQ instance
func (a *RabbitMQ) Connect() (*amqp.Connection, error) {
	if a.Connection != nil {
		return a.Connection, nil
	}
	conn, err := amqp.Dial(a.URL)
	if err != nil {
		// TODO: better logging here
		a.Log.Logger.Log(err)
		conn = nil
	} else {
		a.Connection = conn
	}
	return conn, err
}

// Close will close the open connection attached to the RabbitMQ instance
func (a *RabbitMQ) Close() error {
	var err error
	if a.Connection != nil {
		err = a.Connection.Close()

		if err != nil {
			// TODO: better logging here
			a.Log.Logger.Log(err)
		}
	}
	return err
}

// Channel ..
func (a *RabbitMQ) Channel() (*amqp.Channel, error) {
	var err error
	var ch *amqp.Channel
	if a.Connection == nil {
		err = NewApplicationError("Connect to rabbit before tring to get a channel", nil)
		// TODO: better logging here
		a.Log.Logger.Log(err)
	} else {
		ch, err = a.Connection.Channel()
		if err != nil {
			// TODO: better logging here
			a.Log.Logger.Log(err)
		}
	}

	return ch, err
}

// NewQueue will create a new queue
func (a *RabbitMQ) NewQueue(name string, durable bool, autoDelete bool) (amqp.Queue, error) {
	var err error
	var channel *amqp.Channel
	var queue amqp.Queue
	channel, err = a.Channel()

	if err == nil {
		queue, err = channel.QueueDeclare(
			name,
			durable,
			autoDelete,
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)

		if err != nil {
			// TODO: better logging here
			a.Log.Logger.Log(err)
		}
	}
	return queue, err
}

// NewExchange will create a new exchange
func (a *RabbitMQ) NewExchange(name string, t string, durable bool, autoDelete bool) error {

	var err error
	var channel *amqp.Channel
	channel, err = a.Channel()

	if err == nil {
		err = channel.ExchangeDeclare(
			name,       // name
			t,          // type
			durable,    // durable
			autoDelete, // auto-deleted
			false,      // internal
			false,      // no-wait
			nil,        // arguments
		)

		if err != nil {
			// TODO: better logging here
			a.Log.Logger.Log(err)
		}
	}

	return err
}

// OneWayPublisherEndpoint will create a 'send and forget' publisher endpoint
func (a *RabbitMQ) OneWayPublisherEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	c := NewCtx()
	corrid := c.GetCorrelationFromContext(ctx)
	duration, _ := c.GetTimeoutFromContext(ctx)
	var channel amqptransport.Channel
	var queue *amqp.Queue
	_, _ = a.Connect()
	channel, _ = a.Channel()
	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		channel,
		queue,
		encodeFunc,
		a.NoopResponseDecoder,
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
func (a *RabbitMQ) NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func (a *RabbitMQ) DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		var err error
		marshall := MessageMarshall{}
		*p, err = marshall.Marshal(ctx, exchangeName, request)
		return err
	}
	return f
}

// NewSubscriber ...
func (a *RabbitMQ) NewSubscriber(endpoint endpoint.Endpoint, exchangeName string, dec amqptransport.DecodeRequestFunc) *amqptransport.Subscriber {

	// todo: cache it?

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
