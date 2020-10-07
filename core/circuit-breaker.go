package core

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
)

// CircuitBreaker ...
type CircuitBreaker struct {
}

// NewCircuitBreakerator ....
func NewCircuitBreakerator() CircuitBreaker {
	return CircuitBreaker{}
}

// NewDefaultHystrixCommandMiddleware ...
func (cb CircuitBreaker) NewDefaultHystrixCommandMiddleware(commandName string) endpoint.Middleware {
	config := NewHystrixCommandConfig()
	return cb.NewHystrixCommandMiddleware(commandName, config)
}

// NewHystrixCommandMiddleware ...
func (cb CircuitBreaker) NewHystrixCommandMiddleware(commandName string, config HystrixCommandConfig) endpoint.Middleware {
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
