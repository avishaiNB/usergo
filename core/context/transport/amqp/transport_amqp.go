package amqp

import (
	"context"

	"github.com/thelotter-enterprise/usergo/core/context/transport"
)

type amqptransport struct {
}

// NewAMQPTransport will create a new AMQP transport
func NewAMQPTransport() transport.Transport {
	return amqptransport{}
}

func (amqptrans amqptransport) Read(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	return ctx, nil
}

func (amqptrans amqptransport) Write(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	return ctx, nil
}
