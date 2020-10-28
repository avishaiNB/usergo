package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
)

// MessageMarshall concrete instance of IMarshaller.
type MessageMarshall struct{}

// MessageMarshaller is responsible of transforming the incoming and outcomming message to have the format needed.
type MessageMarshaller interface {
	Marshal(ctx context.Context, exchangeName string, data interface{}) (amqp.Publishing, error)
	Unmarshal(amqpMsg amqp.Delivery) (*MessagePayload, error)
}

// Unmarshal will transform the message received from rabbitmq
func (m *MessageMarshall) Unmarshal(amqpMsg amqp.Delivery) (*MessagePayload, error) {
	wrapper := Message{
		Payload: &MessagePayload{},
	}
	err := json.Unmarshal(amqpMsg.Body, &wrapper)
	if err != nil {
		return &MessagePayload{}, err
	}
	return wrapper.Payload, nil
}
