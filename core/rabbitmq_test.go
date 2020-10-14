package core_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/streadway/amqp"
	"github.com/thelotter-enterprise/usergo/core"
)

func TestNewRabbitMQ(t *testing.T) {
	username := "user"
	pwd := "pwd"
	host := "localhost"
	vhost := "thelotter"
	port := 5672
	log := core.NewLog(nil)
	r := core.NewRabbitMQ(log, host, port, username, pwd, vhost)

	want := "amqp://user:pwd@localhost:5672/thelotter"
	is := r.URL
	if is != want {
		t.Fail()
	}
}

type rabbitRequest struct {
	ID   int
	Name string
}

func TestPublisherEndpoint(t *testing.T) {
	username := "thelotter"
	pwd := "Dhvbuo1"
	host := "int-k8s1"
	vhost := "thelotter"
	port := 32672
	exchangeName := "exchange1"
	ctx := context.Background()
	//req := core.Request{}.Wrap(ctx, rabbitRequest{ID: 1, Name: "guy kolbis"})
	marshall := core.MessageMarshall{}
	req, err := marshall.Marshal(ctx, exchangeName, rabbitRequest{ID: 1, Name: "guy kolbis"})
	log := core.NewLogWithDefaults()
	r := core.NewRabbitMQ(log, host, port, username, pwd, vhost)

	ep := r.OneWayPublisherEndpoint(
		ctx,
		exchangeName,
		func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
			//req, _ := request.(core.Request)
			b, err := json.Marshal(req)
			if err != nil {
				return err
			}
			p.Body = b
			return nil
		},
		func(_ context.Context, d *amqp.Delivery) (response interface{}, err error) {
			return struct{}{}, nil
		},
	)

	_, err = ep(ctx, req)

	if err != nil {
		t.Error(err)
	}
}
