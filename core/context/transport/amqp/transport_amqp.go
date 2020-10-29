package amqp

import (
	"context"

	kitamqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/context/transport"
)

type amqptransport struct {
}

// NewAMQPTransport will create a new AMQP transport
func NewAMQPTransport() transport.Transport {
	return amqptransport{}
}

func (amqptrans amqptransport) Read(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	pub := req.(*amqp.Publishing)
	corrid := pub.CorrelationId

	ctx = tlectx.SetCorrealtion(ctx, corrid)

	// need to read the deadline and duration and apply deadline
	// ctx.Deadline()

	return ctx, nil
}

func (amqptrans amqptransport) Write(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	// message := req.(*rabbitmq.Message)
	// newCtx, cancel := transport.CreateTransportContext(parentContext)
	// duration, deadline := tlectx.GetTimeout(newCtx)
	// corrid := tlectx.GetCorrelation(newCtx)
	// message.CorrelationID = payload.CorrelationID
	// message.Deadline = deadline
	// message.Duration = duration

	return ctx, nil
}

// ReadMessageRequestFunc will be executed once a message is consumed
// it will read from the delivery and will create a context
func ReadMessageRequestFunc() kitamqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		t := NewAMQPTransport()
		newCtx, _ := t.Read(ctx, pub)
		return newCtx
	}
}
