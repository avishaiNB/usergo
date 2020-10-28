package rabbitmq2

import (
	"github.com/streadway/amqp"
)

// ErrorHandler will decide what to do when an error it raise from consumers.
type ErrorHandler interface {
	Handler(ch *amqp.Channel, builder Topology, queueName string, msg amqp.Delivery, stackTrace error) error
}

// ErrorQueueHandler will redirect message to an error queue for future audit.
type ErrorQueueHandler struct{}

// Handler ...
func (e ErrorQueueHandler) Handler(ch *amqp.Channel, topology Topology, queueName string, msg amqp.Delivery, stackTrace error) error {
	errorQueue := queueName + "_error"
	q, err := topology.BuildDurableQueue(
		ch,
		errorQueue,
	)
	if err != nil {
		return err
	}

	headers := msg.Headers
	if headers == nil {
		headers = amqp.Table{}
	}
	headers["Error"] = stackTrace.Error()

	err = topology.Publish(
		ch,
		"",
		q.Name,
		amqp.Publishing{
			Headers: headers,
			Body:    msg.Body,
		},
	)

	return err
}
