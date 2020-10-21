package ratelimit

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

// RateLimiter ...
type RateLimiter struct {
	DefaultMaxQueriesPerSecond int
	DefaultLimit               rate.Limit
}

// NewRateLimitator ...
func NewRateLimitator() RateLimiter {
	return RateLimiter{
		DefaultMaxQueriesPerSecond: 100,
		DefaultLimit:               rate.Every(time.Second),
	}
}

// NewDefaultErrorLimitterMiddleware returns an endpoint.Middleware that acts as a rate limiter.
// Requests that would exceed the maximum request rate are simply rejected with an error.
// NewDefaultLimitterMiddleware will create a middleware which retry every second and allows 100 concurrent requests per second.
func (limiter RateLimiter) NewDefaultErrorLimitterMiddleware() endpoint.Middleware {
	return limiter.NewErrorLimitterMiddleware(limiter.DefaultLimit, limiter.DefaultMaxQueriesPerSecond)
}

// NewErrorLimitterMiddleware returns an endpoint.Middleware that acts as a rate limiter.
// Requests that would exceed the maximum request rate are simply rejected with an error.
func (limiter RateLimiter) NewErrorLimitterMiddleware(limit rate.Limit, burst int) endpoint.Middleware {
	l := rate.NewLimiter(limit, burst)
	middleware := ratelimit.NewErroringLimiter(l)
	return middleware
}

// NewDefaultDelayingLimitterMiddleware returns an endpoint.Middleware that acts as a request throttler.
// Requests that would exceed the maximum request rate are delayed via the Waiter function
// NewDefaultLimitterMiddleware will create a middleware which retry every second and allows 100 concurrent requests per second.
func (limiter RateLimiter) NewDefaultDelayingLimitterMiddleware() endpoint.Middleware {
	return limiter.NewDelayingLimitterMiddleware(limiter.DefaultLimit, limiter.DefaultMaxQueriesPerSecond)
}

// NewDelayingLimitterMiddleware returns an endpoint.Middleware that acts as a request throttler.
// Requests that would exceed the maximum request rate are delayed via the Waiter function
func (limiter RateLimiter) NewDelayingLimitterMiddleware(limit rate.Limit, burst int) endpoint.Middleware {
	l := rate.NewLimiter(limit, burst)
	middleware := ratelimit.NewDelayingLimiter(l)
	return middleware
}
