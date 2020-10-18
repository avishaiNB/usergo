package rabbitmq_test

import (
	"context"
	"testing"

	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
)

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
	conn := rabbitmq.NewConnectionMeta(host, port, username, pwd, vhost)
	r := rabbitmq.NewRabbitMQ(log, conn)
	ep := r.OneWayPublisherEndpoint(ctx, exchangeName, r.DefaultRequestEncoder(exchangeName))
	_, err := ep(ctx, req)

	if err != nil {
		t.Error(err)
	}
}
