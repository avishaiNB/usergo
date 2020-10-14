package core

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
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

func (a *RabbitMQ) channel() (*amqp.Channel, error) {
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
	channel, err = a.channel()

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
	channel, err = a.channel()

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
func (a *RabbitMQ) OneWayPublisherEndpoint(
	ctx context.Context,
	exchangeName string,
	encodeFunc amqptransport.EncodeRequestFunc,
	decodeFunc amqptransport.DecodeResponseFunc,
) endpoint.Endpoint {
	c := NewCtx()
	corrid := c.GetCorrelationFromContext(ctx)
	duration, _ := c.GetTimeoutFromContext(ctx)
	var channel amqptransport.Channel
	var queue *amqp.Queue
	a.Connection, _ = a.Connect()
	channel, _ = a.channel()
	// queue name is not important for one way. So as long as it is not nil, it should be fine.
	queue = &amqp.Queue{Name: ""}

	publisher := amqptransport.NewPublisher(
		channel,
		queue,
		encodeFunc,
		decodeFunc,
		amqptransport.PublisherBefore(
			amqptransport.SetCorrelationID(corrid),
			// TODO: need to configure the below:
			// TODO: we need to append headers: correlation ID, deadline and duration
			// amqptransport.SetPublishDeliveryMode()
			// amqptransport.SetContentEncoding()
			// amqptransport.SetPublishKey()
			// amqptransport.SetContentType(),
			amqptransport.SetPublishExchange(exchangeName)),
		amqptransport.PublisherTimeout(duration),
		amqptransport.PublisherDeliverer(amqptransport.SendAndForgetDeliverer),
	)

	return publisher.Endpoint()
}
