package amqp

import (
	"context"
)

type amqptransport struct {
}

// NewAMQPTransport will create a new AMQP transport
func NewAMQPTransport() Transport {
	return amqptransport{}
}

func (amqptransport httptransport) Read(ctx context.Context, req interface{}) context.Context {
	return ctx
}

func (amqptransport httptransport) Write(ctx context.Context, req interface{}) context.Context {
	return ctx
}
