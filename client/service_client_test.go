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
	serviceName := "test"
	logger := makeLogger()
	serviceDiscoverator := makeServiceDiscovery(logger)
	breakerator := makeCircuitBreakerator()
	limitator := makeRateLimitator()
	inst := makeInstrumentator(serviceName)
	c := client.NewServiceClient(logger, serviceDiscoverator, breakerator, limitator, inst, serviceName)

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
	sd := core.NewServiceDiscovery(logger)
	sd.WithConsul(consulAddress)
	return &sd
}

func makeCircuitBreakerator() core.CircuitBreakerator {
	return core.NewCircuitBreakerator()
}

func makeRateLimitator() core.RateLimitator {
	return core.NewRateLimitator()
}

func makeInstrumentator(serviceName string) core.Instrumentor {
	return core.NewInstrumentor(serviceName)
}
