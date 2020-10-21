package svc

import (
	"context"

	"github.com/streadway/amqp"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
)

// UserAMQPConsumerEndpoints ...
type UserAMQPConsumerEndpoints struct {
	Service   Service
	Log       tlelogger.Log
	Tracer    tletracer.Tracer
	Consumers *[]tlerabbitmq.Consumer
	RabbitMQ  *tlerabbitmq.RabbitMQ
}

// NewUserAMQPConsumerEndpoints will create all the AMQP endpopints
func NewUserAMQPConsumerEndpoints(log tlelogger.Log, tracer tletracer.Tracer, service Service, rabbitMQ *tlerabbitmq.RabbitMQ) *UserAMQPConsumerEndpoints {
	userEndpoints := UserAMQPConsumerEndpoints{
		Log:      log,
		Tracer:   tracer,
		Service:  service,
		RabbitMQ: rabbitMQ,
	}

	userEndpoints.Consumers = userEndpoints.makeConsumerEndpoints()

	return &userEndpoints
}

func (a UserAMQPConsumerEndpoints) makeConsumerEndpoints() *[]tlerabbitmq.Consumer {
	consumers := []tlerabbitmq.Consumer{}
	ep := newLoginEndpoint(a.Service)
	consumer := tlerabbitmq.NewConsumer(ep.Name, ep.Exchange, ep.Queue, ep.EP, ep.Dec)
	consumers = append(consumers, consumer)
	return &consumers
}

func newLoginEndpoint(service Service) tlerabbitmq.EndpointMeta {
	return tlerabbitmq.EndpointMeta{
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
