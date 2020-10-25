package circuitbreaker

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
)

// NewDefaultHystrixCommandMiddleware ...
func NewDefaultHystrixCommandMiddleware(commandName string) endpoint.Middleware {
	config := NewHystrixCommandConfig()
	return NewHystrixCommandMiddleware(commandName, config)
}

// NewHystrixCommandMiddleware ...
func NewHystrixCommandMiddleware(commandName string, config HystrixCommandConfig) endpoint.Middleware {
	hystrixConfig := hystrix.CommandConfig{
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,
		Timeout:                config.Timeout,
	}
	hystrix.ConfigureCommand(commandName, hystrixConfig)
	breaker := circuitbreaker.Hystrix(commandName)
	return breaker
}
