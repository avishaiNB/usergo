package svc

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/thelotter-enterprise/usergo/core"
	tlamqp "github.com/thelotter-enterprise/usergo/core/amqp"

	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
)

// UserAMQPConsumerEndpoints ...
type UserAMQPConsumerEndpoints struct {
	Service   Service
	Log       core.Log
	Tracer    tletracer.Tracer
	Consumers *[]tlamqp.RabbitMQConsumer
	RabbitMQ  *tlamqp.RabbitMQ
}

// NewUserAMQPConsumerEndpoints will create all the AMQP endpopints
func NewUserAMQPConsumerEndpoints(log core.Log, tracer tletracer.Tracer, service Service, rabbitMQ *tlamqp.RabbitMQ) *UserAMQPConsumerEndpoints {
	userEndpoints := UserAMQPConsumerEndpoints{
		Log:      log,
		Tracer:   tracer,
		Service:  service,
		RabbitMQ: rabbitMQ,
	}

	userEndpoints.Consumers = userEndpoints.makeConsumerEndpoints()

	return &userEndpoints
}

func (a UserAMQPConsumerEndpoints) makeConsumerEndpoints() *[]tlamqp.RabbitMQConsumer {
	consumers := []tlamqp.RabbitMQConsumer{}
	ep := newLoginEndpoint(a.Service)
	consumer := a.RabbitMQ.NewConsumer(ep.Name, ep.Exchange, ep.Queue, ep.EP, ep.Dec)
	consumers = append(consumers, consumer)
	return &consumers
}

func newLoginEndpoint(service Service) tlamqp.AMQPEndpoint {
	return tlamqp.AMQPEndpoint{
		EP: func(ctx context.Context, request interface{}) (interface{}, error) {
			err := service.ConsumeLoginCommand(ctx, 1)
			return true, err
		},
		Queue: "queue2",
		Dec: func(_ context.Context, d *amqp.Delivery) (interface{}, error) {
			return nil, nil
		},
		Exchange: "exchange2",
		Name:     "user_login_consumer",
	}
}
