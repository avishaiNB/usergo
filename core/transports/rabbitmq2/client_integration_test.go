// +build integration

package rabbitmq2

import (
	"context"
	"encoding/json"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const (
	amqpURI = "amqp://guest:guest@localhost:5672/thelotter"
)

type PrivateHandler struct {
	targetExchange string
	done           chan bool
	result         PrivateEvent
	err            error
}

func (p *PrivateHandler) exchangeName() string {
	return p.targetExchange
}

func (p *PrivateHandler) handler(ctx context.Context, message *message) error {
	p.err = json.Unmarshal(message.Data, &p.result)
	p.done <- true
	return p.err
}

type PrivateEvent struct {
	Numbers []int `json:"numbers"`
}

func setup(t *testing.T) (*client, error) {
	config := DefaultRabbitMQConfig(amqpURI, &FakeLogger{})
	serviceName := uuid.NewV4().String()
	client, err := newClient(config, serviceName)

	t.Cleanup(func() {
		cleanup(t, client)
	})

	return client, err

}

func cleanup(t *testing.T, client *client) {

	ch, err := client.commandSubscriber.connection.Channel()
	defer ch.Close()

	for _, consumer := range client.commandSubscriber.consumers {
		err := ch.ExchangeDelete(consumer.exchangeName(), false, false)
		assert.NoError(t, err)
	}

	messagesPurged, err := ch.QueueDelete(client.commandSubscriber.queueName, false, false, false)
	assert.Equal(t, 0, messagesPurged)
	assert.NoError(t, err)

	ch, err = client.privateSubscriber.connection.Channel()
	defer ch.Close()

	for _, consumer := range client.privateSubscriber.consumers {
		err := ch.ExchangeDelete(consumer.exchangeName(), false, false)
		assert.NoError(t, err)
	}

	messagesPurged, err = ch.QueueDelete(client.privateSubscriber.queueName, false, false, false)
	assert.Equal(t, 0, messagesPurged)
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Close(ctx)
	assert.NoError(t, err)

}

func TestSimplePublisherSubscribe(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)

	consumer := &PrivateHandler{
		targetExchange: uuid.NewV4().String(),
		done:           make(chan bool),
		err:            nil,
		result:         PrivateEvent{},
	}

	ctx := context.Background()
	err = client.RegisterPrivateConsumers(ctx, consumer)
	assert.NoError(t, err)

	err = client.Run(ctx)
	assert.NoError(t, err)

	privateEvent := PrivateEvent{
		Numbers: []int{1, 2, 3, 4},
	}

	err = client.Publish(ctx, consumer.targetExchange, privateEvent)
	assert.NoError(t, err)

	<-consumer.done
	assert.NoError(t, consumer.err)
	assert.Equal(t, privateEvent, consumer.result)
}

type NonImportantHandler struct {
	done   chan bool
	result interface{}
	err    error
}

func (p *NonImportantHandler) exchangeName() string {
	return "non-important-exchange"
}

func (p *NonImportantHandler) handler(ctx context.Context, message *message) error {

	err := json.Unmarshal(message.Data, &p.result)
	p.done <- true
	p.err = err

	return err
}

func TestEventsNotMix(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)

	targetConsumer := &PrivateHandler{
		targetExchange: uuid.NewV4().String(),
		done:           make(chan bool),
		err:            nil,
		result:         PrivateEvent{},
	}

	nonImportantHandler := &NonImportantHandler{
		err:    nil,
		result: nil,
	}

	ctx := context.Background()
	err = client.RegisterPrivateConsumers(ctx, targetConsumer, nonImportantHandler)
	assert.NoError(t, err)

	err = client.Run(ctx)
	assert.NoError(t, err)

	privateEvent := PrivateEvent{
		Numbers: []int{1, 2, 3, 4},
	}

	err = client.Publish(ctx, targetConsumer.targetExchange, privateEvent)
	assert.NoError(t, err)

	<-targetConsumer.done
	assert.NoError(t, targetConsumer.err)
	assert.Equal(t, privateEvent, targetConsumer.result)

	assert.NoError(t, nonImportantHandler.err)
	assert.Equal(t, nil, nonImportantHandler.result)
}

func getQueueInfo(subscriber *subscriber) (amqp.Queue, error) {
	ch, err := subscriber.connection.Channel()
	defer ch.Close()
	if err != nil {
		return amqp.Queue{}, err
	}

	return ch.QueueInspect(subscriber.queueName)
}

func TestNoConsumerForMessage(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Run(ctx)
	assert.NoError(t, err)

	privateEvent := PrivateEvent{
		Numbers: []int{1, 2, 3, 4},
	}

	randomExchange := uuid.NewV4().String()
	err = client.Publish(ctx, randomExchange, privateEvent)

	queue, err := getQueueInfo(client.privateSubscriber)
	assert.NoError(t, err)
	assert.Equal(t, 1, queue.Consumers)
	assert.Equal(t, 0, queue.Messages)

	queue, err = getQueueInfo(client.commandSubscriber)
	assert.NoError(t, err)
	assert.Equal(t, 1, queue.Consumers)
	assert.Equal(t, 0, queue.Messages)
}
