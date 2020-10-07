package core

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
)

// CircuitBreakerator ...
type CircuitBreakerator struct {
}

// NewCircuitBreakerator ....
func NewCircuitBreakerator() CircuitBreakerator {
	return CircuitBreakerator{}
}

// NewDefaultHystrixCommandMiddleware ...
func (cb CircuitBreakerator) NewDefaultHystrixCommandMiddleware(commandName string) endpoint.Middleware {
	config := NewHystrixCommandConfig()
	return cb.NewHystrixCommandMiddleware(commandName, config)
}

// NewHystrixCommandMiddleware ...
func (cb CircuitBreakerator) NewHystrixCommandMiddleware(commandName string, config HystrixCommandConfig) endpoint.Middleware {
	hystrixConfig := hystrix.CommandConfig{
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,

		Timeout: config.Timeout,
	}
	hystrix.ConfigureCommand(commandName, hystrixConfig)
	breaker := circuitbreaker.Hystrix(commandName)
	return breaker
}
