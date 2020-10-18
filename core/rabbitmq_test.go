package core_test

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/core"
)

func TestNewRabbitMQ(t *testing.T) {
	username := "user"
	pwd := "pwd"
	host := "localhost"
	vhost := "thelotter"
	port := 5672
	logger := log.NewNopLogger()
	log := core.NewLog(logger, 2)
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

// Integration Test! Should not run on automated build
func TestPublisherEndpoint(t *testing.T) {
	username := "thelotter"
	pwd := "Dhvbuo1"
	host := "int-k8s1"
	vhost := "thelotter"
	port := 32672
	exchangeName := "exchange1"
	ctx := context.Background()
	req := rabbitRequest{ID: 1, Name: "guy kolbis"}
	log := core.NewLogWithDefaults()
	r := core.NewRabbitMQ(log, host, port, username, pwd, vhost)
	ep := r.OneWayPublisherEndpoint(ctx, exchangeName, r.DefaultRequestEncoder(exchangeName))
	_, err := ep(ctx, req)

	if err != nil {
		t.Error(err)
	}
}
