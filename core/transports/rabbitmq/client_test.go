package rabbitmq_test

import (
	"context"
	"testing"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	"github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	"github.com/thelotter-enterprise/usergo/shared"
)

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
	req := shared.LoggedInCommandData{ID: 1, Name: "guy kolbis"}
	message, newCtx, _ := rabbitmq.NewMessage(ctx, req, "thelotter.userloggedin")

	logManager := tlelogger.NewNopManager()
	connInfo := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)
	conn := rabbitmq.NewConnectionManager(connInfo)
	publisher := rabbitmq.NewPublisher(&conn)
	client := rabbitmq.NewClient(&conn, &logManager, &publisher, nil)
	err := client.Publish(newCtx, &message, exchangeName, rabbitmq.DefaultRequestEncoder(exchangeName))

	if err != nil {
		t.Error(err)
	}
}
