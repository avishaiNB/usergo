package rabbitmq

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
)

// EndpointMeta ...
type EndpointMeta struct {
	EP       endpoint.Endpoint
	Name     string
	Exchange string
	Queue    string
	Dec      amqptransport.DecodeRequestFunc
}

// MessagePayload is the basic unit used to send and received data thought rabbitmq
// or more information about the URN http://masstransit-project.com/MassTransit/architecture/interoperability.html
type MessagePayload struct {
	Data          interface{} `json:"data,omitempty"`
	CorrelationID string      `json:"correlationId"`
	// Timestamp      time.Time              `json:"timestamp"`
	AdditionalData map[string]interface{} `json:"additionalData"`
}

// Message is the Message classes used by masstransit for every message
type Message struct {
	MessageType   []string        `json:"messageType"`
	Payload       *MessagePayload `json:"message"`
	Deadline      time.Time       `json:"deadline"`
	Duration      time.Duration   `json:"duration"`
	CorrelationID string          `json:"correlationId"`
}
