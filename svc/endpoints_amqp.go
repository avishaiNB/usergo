package svc

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/thelotter-enterprise/usergo/core"
)

// UserAMQPConsumerEndpoints ...
type UserAMQPConsumerEndpoints struct {
	Service   Service
	Log       core.Log
	Tracer    core.Tracer
	Consumers *[]core.RabbitMQConsumer
	RabbitMQ  *core.RabbitMQ
}

// NewUserAMQPConsumerEndpoints will create all the AMQP endpopints
func NewUserAMQPConsumerEndpoints(log core.Log, tracer core.Tracer, service Service, rabbitMQ *core.RabbitMQ) *UserAMQPConsumerEndpoints {
	userEndpoints := UserAMQPConsumerEndpoints{
		Log:      log,
		Tracer:   tracer,
		Service:  service,
		RabbitMQ: rabbitMQ,
	}

	userEndpoints.Consumers = userEndpoints.makeConsumerEndpoints()

	return &userEndpoints
}

func (a UserAMQPConsumerEndpoints) makeConsumerEndpoints() *[]core.RabbitMQConsumer {
	consumers := []core.RabbitMQConsumer{}
	ep := newLoginEndpoint()
	consumer := a.RabbitMQ.NewConsumer(ep.Name, ep.Exchange, ep.Queue, ep.EP, ep.Dec)
	consumers = append(consumers, consumer)
	return &consumers
}

func newLoginEndpoint() core.AMQPEndpoint {
	return core.AMQPEndpoint{
		EP: func(_ context.Context, request interface{}) (interface{}, error) {
			return true, nil
		},
		Queue: "queue1",
		Dec: func(_ context.Context, d *amqp.Delivery) (interface{}, error) {
			return nil, nil
		},
		Exchange: "exchange1",
		Name:     "user_login_consumer",
	}
}
