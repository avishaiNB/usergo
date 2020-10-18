package amqp

import (
	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
)

// Endpoint ...
type Endpoint struct {
	EP       endpoint.Endpoint
	Name     string
	Exchange string
	Queue    string
	Dec      amqptransport.DecodeRequestFunc
}
