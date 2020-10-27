package rabbitmq

import (
	"context"
)

// Consumer will listen to events received from an exchange and react to them.
type Consumer interface {
	// Name gets the name of the consumer
	// It will be used as a key in a consumer map
	Name() string

	// MessageURNs gets the message URN list which the consumer wants to be notified about
	// When a message enters the queue, we need to route it to the relevant consumer
	// This is done using ghe URNs
	MessageURNs() []string

	// Handler will be called when a the consumer was matched and the message need to be processed by it
	Handler(ctx context.Context, message *Message) error
}
