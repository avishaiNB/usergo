package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/client"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlesd "github.com/thelotter-enterprise/usergo/core/servicediscovery"
)

func TestClientIntegration(t *testing.T) {
	serviceName := "test"
	logger := tlelogger.NewNopManager()
	ctx := context.Background()
	id := 1

	serviceDiscoverator := makeServiceDiscovery(logger)
	c := client.NewServiceClientWithDefaults(&logger, serviceDiscoverator, serviceName)

	response := c.GetUserByID(ctx, id)

	if response.Data == nil {
		t.Fail()
	}
}

func makeLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", ":8080", "caller", log.DefaultCaller)

	return logger
}

func makeServiceDiscovery(logger tlelogger.Manager) *tlesd.ServiceDiscovery {
	consulAddress := "localhost:8500"
	sd := tlesd.NewConsulServiceDiscovery(logger, consulAddress)
	return &sd
}
