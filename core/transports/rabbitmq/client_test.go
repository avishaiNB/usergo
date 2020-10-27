package rabbitmq_test

import (
	"context"
	"testing"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	"github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
)

type rabbitRequest struct {
	ID   int
	Name string
}

const (
	exchangeName string = "exchange1"
	username     string = "thelotter"
	pwd          string = "Dhvbuo1"
	host         string = "int-k8s1"
	vhost        string = "thelotter"
	port         int    = 32672
)

// Integration Test! Should not run on automated build
func TestPublisherEndpoint(t *testing.T) {
	ctx := context.Background()
	req := rabbitRequest{ID: 1, Name: "guy kolbis"}
	logManager := tlelogger.NewNopManager()
	conn := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)
	rabbit := rabbitmq.NewClient(&logManager, conn)

	rabbit.OpenConnection()
	err := rabbit.PublishOneWay(ctx, req, exchangeName, rabbit.DefaultRequestEncoder(exchangeName))

	if err != nil {
		t.Error(err)
	}
}

// Integration Test! Should not run on automated build
func TestSubscribers(t *testing.T) {
	ctx := context.Background()
	logManager := tlelogger.NewNopManager()
	conn := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)

	rabbit := rabbitmq.NewClient(&logManager, conn, nil, nil)
	defer rabbit.CloseConnection()
	rabbit.OpenConnection()
	//err := rabbit.PublishOneWay(ctx, req, exchangeName, rabbit.DefaultRequestEncoder(exchangeName))
}
