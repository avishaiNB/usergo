package rabbitmq2

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/streadway/amqp"
)

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second
)

var (
	errPublisherAlreadyClosed = errors.New("Publisher connection is already closed")
)

type publisher struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	config        *Config
	alive         bool
	isConnected   bool
	done          chan os.Signal
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
}

func newPublisher(config *Config) *publisher {
	publisher := publisher{
		config: config,
		alive:  true,
	}
	ctx := context.Background()
	publisher.connect(ctx, config)
	go publisher.handleReconnect(ctx, config)
	return &publisher
}

func (p *publisher) handleReconnect(ctx context.Context, config *Config) {
	for p.alive {
		//p.config.Logger.Info(ctx, "Attempting to connect publisher to rabbitMQ: %s\n", config.AmqpURI)
		var retryCount int
		for !p.connect(ctx, config) {
			//p.config.Logger.Error(ctx, "disconnected from rabbitMQ and failed to connect")
			time.Sleep(reconnectDelay + time.Duration(retryCount)*time.Second)
			retryCount++
		}
		select {
		case <-p.done:
			return
		case <-p.notifyClose:
			p.isConnected = false
		}
	}
}

func (p *publisher) connect(ctx context.Context, config *Config) bool {
	conn, err := amqp.Dial(config.ConnectionInfo.URL)
	if err != nil {
		//p.config.Logger.Error(ctx, "failed to dial rabbitMQ server: %v", err)
		return false
	}
	ch, err := conn.Channel()
	if err != nil {
		//p.config.Logger.Error(ctx, "failed connecting to channel: %v", err)
		return false
	}

	p.changeConnection(ctx, conn, ch)
	p.isConnected = true
	return true
}

func (p *publisher) changeConnection(ctx context.Context, connection *amqp.Connection, channel *amqp.Channel) {
	p.connection = connection
	p.channel = channel
	p.notifyClose = make(chan *amqp.Error)
	p.notifyConfirm = make(chan amqp.Confirmation)
	p.channel.NotifyClose(p.notifyClose)

	err := p.channel.Confirm(false)
	if err != nil {
		//p.config.Logger.Error(ctx, "publisher confirms not supported")
		close(p.notifyConfirm)
	} else {
		p.channel.NotifyPublish(p.notifyConfirm)
	}
}

// Publish will publish the message to rabbitMQ and wait until message has been confirm
// In case the broker is not reachable, will retry every second
func (p *publisher) Publish(ctx context.Context, exchangeName string, data interface{}) error {
	if !p.isConnected {
		return errors.New("failed to push push: not connected")
	}
	for {
		err := p.unsafePush(ctx, exchangeName, data)
		if err != nil {
			continue
		}
		select {
		case confirm := <-p.notifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-time.After(1 * time.Second):
		}
	}
}

// unsafePush will publish the message to rabbitMQ, but won't deail with errors
func (p *publisher) unsafePush(ctx context.Context, exchangeName string, data interface{}) error {
	err := p.config.Topology.BuildExchange(p.channel, exchangeName)
	if err != nil {
		return err
	}

	msg, err := p.config.Marshaller.Marshal(ctx, exchangeName, data)
	if err != nil {
		return err
	}

	err = p.config.Topology.Publish(p.channel, exchangeName, "", msg)
	if err != nil {
		return err
	}

	//p.config.Logger.Info(ctx, "Message published", string(msg.Body))
	return nil
}

// Close will shutdown the publisher gracely
func (p *publisher) Close() error {
	if !p.isConnected {
		return errPublisherAlreadyClosed
	}
	// ctx := context.Background()
	p.alive = false
	err := p.channel.Close()
	if err != nil {
		//p.config.Logger.Error(ctx, "Error closing the channel", err)
	}
	err = p.connection.Close()
	if err != nil {
		//p.config.Logger.Error(ctx, "Error closing the connection", err)
	}
	return nil
}
