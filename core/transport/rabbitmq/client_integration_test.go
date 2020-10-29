// +build integration

package rabbitmq_test

import (
	"testing"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	"github.com/thelotter-enterprise/usergo/core/transport/rabbitmq"
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

func TestPublishMessage(t *testing.T) {
	ctx := tlectx.Root()
	req := shared.LoggedInCommandData{ID: 1, Name: "guy kolbis"}
	message := rabbitmq.NewMessage(req, "thelotter.userloggedin")

	logManager := tlelogger.NewNopManager()
	connInfo := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)
	conn := rabbitmq.NewConnectionManager(connInfo)
	publisher := rabbitmq.NewPublisher(&conn)
	client := rabbitmq.NewClient(&conn, &logManager, &publisher, nil)
	err := client.Publish(ctx, &message, exchangeName, rabbitmq.DefaultRequestEncoder(exchangeName))

	if err != nil {
		t.Error(err)
	}
}
