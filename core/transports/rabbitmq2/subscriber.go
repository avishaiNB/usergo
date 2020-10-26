package rabbitmq2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

var (
	errConsumerAlreadyRegisterToExchange = errors.New("Exchange already have an consumer registered")
	errSubscriberAlreadyClosed           = errors.New("Subscriber is already closed")
)

type subscriber struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	queueName   string
	config      *Config
	consumers   map[string]Consumer
	buildQueue  func(channel *amqp.Channel, queueName string) (amqp.Queue, error)
	alive       bool
	isConnected bool
	done        chan os.Signal
	notifyClose chan *amqp.Error
}

// newCommandSubscriber will create a subscriber like TheLotter.XXX.Service-command
func newCommandSubscriber(config *Config, queueBaseName string, consumers ...Consumer) (*subscriber, error) {
	subscriber := subscriber{
		config:     config,
		queueName:  queueBaseName + "-command",
		consumers:  make(map[string]Consumer),
		buildQueue: config.Topology.BuildDurableQueue,
		alive:      true,
	}

	ctx := context.Background()
	err := subscriber.registerConsumer(consumers...)
	if err != nil {
		return nil, err
	}

	go subscriber.handleReconnect(ctx)

	return &subscriber, nil
}

// newPrivateSubscriber will create a subscriber like TheLotter.XXX.Service-private-uuid
func newPrivateSubscriber(config *Config, queueBaseName string, consumers ...Consumer) (*subscriber, error) {
	subscriber := subscriber{
		config:     config,
		queueName:  queueBaseName + "-private-" + uuid.NewV4().String(),
		consumers:  make(map[string]Consumer),
		buildQueue: config.Topology.BuildNonDurableQueue,
		alive:      true,
	}

	ctx := context.Background()
	err := subscriber.registerConsumer(consumers...)
	if err != nil {
		return nil, err
	}

	go subscriber.handleReconnect(ctx)

	return &subscriber, nil
}

// handleReconnect will wait for a connection error on
// notifyClose, and then continuously attempt to reconnect.
func (s *subscriber) handleReconnect(ctx context.Context) {
	for s.alive {
		s.isConnected = false
		//s.config.Logger.Info(ctx, fmt.Sprintf("Attempting to connect subscriber %s to rabbitMQ: %s\n", s.queueName, s.config.AmqpURI))
		var retryCount int
		for !s.connect(ctx, s.config) {
			//s.config.Logger.Error(ctx, "disconnected from rabbitMQ and failed to connect")
			time.Sleep(reconnectDelay + time.Duration(retryCount)*time.Second)
			retryCount++
		}
		select {
		case <-s.done:
			return
		case <-s.notifyClose:
		}
	}
}

// connect will make a single attempt to connect to
// RabbitMq. It returns the success of the attempt.
func (s *subscriber) connect(ctx context.Context, config *Config) bool {
	conn, err := amqp.Dial(config.ConnectionInfo.URL)
	if err != nil {
		//s.config.Logger.Error(ctx, "failed to dial rabbitMQ server: %v", err)
		return false
	}
	ch, err := conn.Channel()
	if err != nil {
		//s.config.Logger.Error(ctx, "failed connecting to channel: %v", err)
		return false
	}

	s.changeConnection(conn, ch)
	s.isConnected = true
	return true
}

// changeConnection takes a new connection to the subscriber,
// and updates the channel listeners to reflect this.
func (s *subscriber) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	s.connection = connection
	s.channel = channel
	s.notifyClose = make(chan *amqp.Error)
	s.channel.NotifyClose(s.notifyClose)
}

// registerConsumer will register the consumers
// If two consumers register targeting the same exchange, an error will be raised
func (s *subscriber) registerConsumer(consumers ...Consumer) error {
	for _, consumer := range consumers {
		_, found := s.consumers[consumer.exchangeName()]
		if found {
			err := fmt.Errorf("registering to the exchange %v: %w", consumer.exchangeName(), errConsumerAlreadyRegisterToExchange)
			return err
		}
		s.consumers[consumer.exchangeName()] = consumer
	}
	return nil
}

// Consume will bind to all the target exchanges and if case of not existing, create them
// Then will be running in loop processing all the incoming messages
func (s *subscriber) Consume(ctx context.Context) error {
	for {
		if s.isConnected {
			break
		}
		//s.config.Logger.Info(ctx, "Waiting to connect with RabbitMQ")
		time.Sleep(1 * time.Second)
	}

	err := s.channel.Qos(s.config.PrefetchCount, 0, false)
	if err != nil {
		return err
	}

	q, err := s.buildQueue(s.channel, s.queueName)
	if err != nil {
		return err
	}

	for _, consumer := range s.consumers {
		err = s.config.Topology.BuildExchange(s.channel, consumer.exchangeName())
		if err != nil {
			return err
		}

		err = s.config.Topology.QueueBind(s.channel, q.Name, consumer.exchangeName())
		if err != nil {
			return err
		}
	}

	msgs, err := s.config.Topology.Consume(s.channel, q.Name)
	if err != nil {
		return err
	}
	go s.processMessages(s.channel, msgs)
	return nil
}

// processMessages will process all the incomming messages redirecting them to the correct consumer
func (s *subscriber) processMessages(ch *amqp.Channel, msgs <-chan amqp.Delivery) {
Consuming:
	for {
		select {
		case msg := <-msgs:

			consumer, found := s.consumers[msg.Exchange]
			if !found {
				msg.Ack(false)
				continue Consuming
			}

			message, err := s.config.Marshaller.Unmarshal(msg)
			if err != nil {
				s.useErrorHandler(context.Background(), ch, msg, err)
				msg.Ack(false)
				continue Consuming
			}

			// TODO: we need to read it from transport!!!
			ctx := context.Background()

			// ctx := context.Context{
			// 	CorrelationID: message.CorrelationID,
			// }
			err = consumer.handler(ctx, message)
			if err != nil {
				s.useErrorHandler(ctx, ch, msg, err)
				msg.Nack(false, true)
				continue Consuming
			}
			msg.Ack(false)
			//s.config.Logger.Info(ctx, "Received a message:", string(msg.Body))
			continue Consuming
		case <-s.notifyClose:
			return
		}
	}
}

// useErrorHandler will run in case there is any exception while processing the messages and
// according to the exceptions handlers configurated execute them
func (s *subscriber) useErrorHandler(ctx context.Context, ch *amqp.Channel, msg amqp.Delivery, err error) {
	for _, handler := range s.config.ErrorHandlers {
		err := handler.Handler(ch, s.config.Topology, s.queueName, msg, err)
		if err == nil {
			break
		}
		//s.config.Logger.Error(ctx, err.Error())
	}
}

// Close will shutdown the subscriber gracely
func (s *subscriber) Close() error {
	if !s.isConnected {
		return errSubscriberAlreadyClosed
	}
	s.alive = false
	// ctx := context.Background()

	// s.config.Logger.Info(ctx, "Closing consumer")
	err := s.channel.Cancel("", false)
	if err != nil {
		//s.config.Logger.Error(ctx, "Error while closing consumer", err)
	}

	err = s.channel.Close()
	if err != nil {
		//s.config.Logger.Error(ctx, "Error closing the channel", err)
	}
	err = s.connection.Close()
	if err != nil {
		//s.config.Logger.Error(ctx, "Error closing the connection", err)
	}
	return nil
}
