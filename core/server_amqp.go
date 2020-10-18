package core

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
)

// AMQPServer ...
type AMQPServer struct {
	Name     string
	Address  string
	Log      Log
	Tracer   Tracer
	RabbitMQ *RabbitMQ
}

// NewAMQPServer ...
func NewAMQPServer(log Log, tracer Tracer, rabbit *RabbitMQ, serviceName string) AMQPServer {
	return AMQPServer{
		Name:     serviceName,
		RabbitMQ: rabbit,
		Log:      log,
		Tracer:   tracer,
	}
}

// AMQPEndpoint ...
type AMQPEndpoint struct {
	EP       endpoint.Endpoint
	Name     string
	Exchange string
	Queue    string
	Dec      amqptransport.DecodeRequestFunc
}

// Run will ...
func (server *AMQPServer) Run(endpoints *[]RabbitMQConsumer) error {
	if endpoints == nil {
		return errors.New("no endpoints")
	}

	//consumers := make(map[string]chan, 1)

	for _, endpoint := range *endpoints {
		_, err := server.RabbitMQ.Consume(&endpoint)

		if err == nil {
			//consumers[endpoint.Consumer] = ch
		}
	}

	return nil
}
