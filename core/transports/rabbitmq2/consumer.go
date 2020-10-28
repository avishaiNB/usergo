package rabbitmq2

import (
	"context"
)

// Consumer will listen to events received from an exchange and react to them.
type Consumer interface {
	exchangeName() string
	handler(ctx context.Context, message *Message) error
}
