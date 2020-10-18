package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/ctx"
)

// Message is the basic unit used to send and received data thought rabbitmq
// or more information about the URN http://masstransit-project.com/MassTransit/architecture/interoperability.html
type Message struct {
	URN            string                 `json:"-"`
	Data           interface{}            `json:"data,omitempty"`
	CorrelationID  string                 `json:"correlationId"`
	AdditionalData map[string]interface{} `json:"additionalData"`
}

// MessageWrapper is the MessageWrapper classes used by masstransit for every message
type MessageWrapper struct {
	MessageType []string `json:"messageType"`
	Message     *Message `json:"message"`
}

// MessageMarshaller is responsible of transforming the incoming and outcomming message to have the format needed.
type MessageMarshaller interface {
	Marshal(ctx context.Context, exchangeName string, data interface{}) (amqp.Publishing, error)
	Unmarshal(amqpMsg amqp.Delivery) (*Message, error)
}

// MessageMarshall concrete instance of IMarshaller.
type MessageMarshall struct{}

// Marshal will transform a message to be publish into rabbitmq
func (m *MessageMarshall) Marshal(ctx context.Context, exchangeName string, data interface{}) (amqp.Publishing, error) {
	urn := fmt.Sprintf("urn:message:%v", exchangeName)
	msg := Message{Data: data, URN: urn}
	msg.CorrelationID = tlectx.GetOrCreateCorrelationID(ctx)
	wrapper := MessageWrapper{MessageType: []string{urn}, Message: &msg}
	body, err := json.Marshal(wrapper)

	if err != nil {
		return amqp.Publishing{}, err
	}

	rabbitMessage := amqp.Publishing{Body: body}
	rabbitMessage.CorrelationId = msg.CorrelationID
	return rabbitMessage, nil
}

// Unmarshal will transform the message received from rabbitmq
func (m *MessageMarshall) Unmarshal(amqpMsg amqp.Delivery) (*Message, error) {
	wrapper := MessageWrapper{
		Message: &Message{},
	}
	err := json.Unmarshal(amqpMsg.Body, &wrapper)
	if err != nil {
		return &Message{}, err
	}
	return wrapper.Message, nil
}
