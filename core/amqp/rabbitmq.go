package amqp

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"github.com/thelotter-enterprise/usergo/core"
	tlectx "github.com/thelotter-enterprise/usergo/core/ctx"
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

	Log core.Log
}

// RabbitMQConsumer ...RabbitMQConsumer
type RabbitMQConsumer struct {
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

// NewRabbitMQ will create a new instance of empty RabbitMQ
func NewRabbitMQ(log core.Log, host string, port int, username string, password string, vhost string) RabbitMQ {
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

	ch, err := a.Channel()
	if err != nil {
		panic(err)
	}

	consumer.Channel = ch

	consumer.NewExchange(consumer.ExchangeName, consumer.ExchangeDurable, consumer.ExchangeAutoDelete)
	consumer.NewQueue(consumer.QueueName, consumer.QueueDurable, consumer.QueueAutoDelete)
	consumer.Bind(consumer.QueueName, consumer.ExchangeName)

	c, err := consumer.Channel.Consume(
		consumer.QueueName,
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
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc) RabbitMQConsumer {

	sub := a.NewSubscriber(endpoint, exchangeName, dec)
	consumer := RabbitMQConsumer{
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

// Channel will create a new rabbitMQ channel
func (a *RabbitMQ) Channel() (*amqp.Channel, error) {
	var err error
	var ch *amqp.Channel
	if a.Connection == nil {
		err = core.NewApplicationError("Connect to rabbit before tring to get a channel", nil)
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
func (c *RabbitMQConsumer) NewQueue(name string, durable bool, autoDelete bool) (amqp.Queue, error) {
	var err error
	var queue amqp.Queue

	queue, err = c.Channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return queue, err
}

// NewExchange will create a new exchange
func (c *RabbitMQConsumer) NewExchange(name string, durable bool, autoDelete bool) error {

	err := c.Channel.ExchangeDeclare(
		name,       // name
		"fanout",   // type
		durable,    // durable
		autoDelete, // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)

	return err
}

// Bind will bind the rabbitMQ queue and exchange together
func (c *RabbitMQConsumer) Bind(queueName string, exchangeName string) error {

	err := c.Channel.QueueBind(
		queueName,
		"", // bindingKey
		exchangeName,
		false, // noWait
		nil,   // arguments
	)

	return err
}

// OneWayPublisherEndpoint will create a 'send and forget' publisher endpoint
func (a *RabbitMQ) OneWayPublisherEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc) endpoint.Endpoint {
	corrid := tlectx.GetCorrelationFromContext(ctx)
	duration, _ := tlectx.GetTimeoutFromContext(ctx)
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
