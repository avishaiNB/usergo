package amqp

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
