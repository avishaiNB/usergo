package core

import (
	"context"
	"errors"

	amqpkit "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// AMQPServer ...
type AMQPServer struct {
	Name     string
	Address  string
	Log      Log
	Tracer   Tracer
	RabbitMQ *RabbitMQ
}

// AMQPEndpoints ...
type AMQPEndpoints struct {
	ServerEndpoints []AMQPEndpoint
}

// NewAMQPEndpoints ...
func NewAMQPEndpoints() AMQPEndpoints {
	return AMQPEndpoints{
		ServerEndpoints: []AMQPEndpoint{},
	}
}

// AMQPEndpoint holds the information needed to build a server endpoint which client can call upon
type AMQPEndpoint struct {
	Endpoint     func(ctx context.Context, request interface{}) (interface{}, error)
	Dec          amqpkit.DecodeRequestFunc
	Enc          amqpkit.EncodeResponseFunc
	ExchangeName string
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

// Run will ...
func (server *AMQPServer) Run(endpoints *AMQPEndpoints) error {
	if endpoints == nil {
		return errors.New("no endpoints")
	}

	_, err := server.RabbitMQ.Connect()
	if err != nil {
		panic(err)
	}

	ch, err := server.RabbitMQ.Channel()
	if err != nil {
		panic(err)
	}

	var consumers = []consumer{}

	for _, endpoint := range endpoints.ServerEndpoints {
		//server.Log.Logger.Log("message", fmt.Sprintf("adding route http://%s/%s", server.Address, endpoint.Path))
		sub := server.RabbitMQ.NewSubscriber(endpoint.Endpoint, endpoint.ExchangeName, endpoint.Dec)
		f := sub.ServeDelivery(ch)
		f(&amqp.Delivery{})
		c := consumer{f: f, sub: sub}
		consumers = append(consumers, c)
	}

	//server.Handler = handlers.LoggingHandler(os.Stdout, server.Router)
	//server.Log.Logger.Log("message", fmt.Sprintf("http server started and listen on %s", server.Address))
	//http.ListenAndServe(server.Address, server.Handler)

	return nil
}

type consumer struct {
	f   func(*amqp.Delivery)
	sub *amqpkit.Subscriber
}
