package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/client"
	"github.com/thelotter-enterprise/usergo/core"
)

func TestClientIntegration(t *testing.T) {
	logger := makeLogger()
	sd := makeServiceDiscovery(logger)
	c := client.NewServiceClient(logger, sd, "test")

	ctx := context.Background()
	id := 1

	response := c.GetUserByID(ctx, id)

	if response.Result == nil {
		t.Fail()
	}
}

func makeLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", ":8080", "caller", log.DefaultCaller)

	return logger
}

func makeServiceDiscovery(logger log.Logger) *core.ServiceDiscoverator {
	consulAddress := "localhost:8080"
	sd, _ := core.NewServiceDiscovery(logger, consulAddress)
	return &sd
}
