package core

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
)

const (
	// DefaultMaxAttempts defines number of retry attempts per request, before giving up
	DefaultMaxAttempts = 5
)

// LoadBalancer ...
type LoadBalancer struct {
	FixedEndpointer sd.FixedEndpointer
}

// NewLoadBalancer ..
func NewLoadBalancer(fixedEndpointer sd.FixedEndpointer) LoadBalancer {
	return LoadBalancer{
		FixedEndpointer: fixedEndpointer,
	}
}

// DefaultRoundRobinWithRetryEndpoint ...
func (b *LoadBalancer) DefaultRoundRobinWithRetryEndpoint(ctx context.Context) endpoint.Endpoint {
	c := NewCtx(ctx)
	maxTime := c.CalcTimeout()
	return b.RoundRobinWithRetryEndpoint(DefaultMaxAttempts, maxTime)
}

// RoundRobinWithRetryEndpoint ..
func (b *LoadBalancer) RoundRobinWithRetryEndpoint(maxAttempts int, maxTime time.Duration) endpoint.Endpoint {
	balancer := lb.NewRoundRobin(b.FixedEndpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)
	return retry
}
