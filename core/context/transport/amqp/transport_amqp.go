package amqp

import (
	"context"

	kitamqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
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

// ReadMessageRequestFunc ...
func ReadMessageRequestFunc() kitamqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		t := NewAMQPTransport()
		newCtx, _ := t.Read(ctx, pub)
		return newCtx
	}
}
