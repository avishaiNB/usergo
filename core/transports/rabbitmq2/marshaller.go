package rabbitmq2

import (
	"encoding/json"
	"fmt"

	"context"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"

	"github.com/streadway/amqp"
)

// Marshaller is responsible of transforming the incoming and outcomming message to have the format needed.
type Marshaller interface {
	Marshal(ctx context.Context, exchangeName string, data interface{}) (amqp.Publishing, error)
	Unmarshal(amqpMsg amqp.Delivery) (*Message, error)
}

// wrapper is the wrapper classes used by masstransit for every message
type wrapper struct {
	MessageType []string `json:"messageType"`
	Message     *Message `json:"message"`
}

// Marshall concrete instance of IMarshaller.
type Marshall struct{}

// Marshal will transform a message to be publish into rabbitmq
func (m *Marshall) Marshal(ctx context.Context, exchangeName string, data interface{}) (amqp.Publishing, error) {

	urn := fmt.Sprintf("urn:message:%v", exchangeName)
	body, err := json.Marshal(data)
	if err != nil {
		return amqp.Publishing{}, err
	}

	msg := Message{
		Data: body,
		URN:  urn,
	}

	msg.CorrelationID = tlectx.GetOrCreateCorrelation(ctx)

	wrapper := wrapper{
		MessageType: []string{
			urn,
		},
		Message: &msg,
	}

	body, err = json.Marshal(wrapper)
	if err != nil {
		return amqp.Publishing{}, err
	}
	return amqp.Publishing{
		Body: body,
	}, nil
}

// Unmarshal will ...
func (m *Marshall) Unmarshal(amqpMsg amqp.Delivery) (*Message, error) {
	wrapper := wrapper{
		Message: &Message{},
	}
	err := json.Unmarshal(amqpMsg.Body, &wrapper)
	if err != nil {
		return &Message{}, err
	}
	return wrapper.Message, nil
}
