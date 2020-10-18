package transports

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	tlectx "github.com/thelotter-enterprise/usergo/core/ctx"
)

const (
	// DefaultMaxAttempts defines number of retry attempts per request, before giving up
	DefaultMaxAttempts = 5
)

// LoadBalancer ...
type LoadBalancer struct {
	FixedEndpointer sd.FixedEndpointer
	Endpointer      *sd.DefaultEndpointer
}

// NewLoadBalancer ..
func NewLoadBalancer(fixedEndpointer sd.FixedEndpointer, endpointer *sd.DefaultEndpointer) LoadBalancer {
	return LoadBalancer{
		FixedEndpointer: fixedEndpointer,
		Endpointer:      endpointer,
	}
}

// DefaultRoundRobinWithRetryEndpoint ...
func (b *LoadBalancer) DefaultRoundRobinWithRetryEndpoint(ctx context.Context) endpoint.Endpoint {
	maxTime, _ := tlectx.CalcTimeoutFromContext(ctx)
	return b.RoundRobinWithRetryEndpoint(DefaultMaxAttempts, maxTime)
}

// RoundRobinWithRetryEndpoint ..
func (b *LoadBalancer) RoundRobinWithRetryEndpoint(maxAttempts int, maxTime time.Duration) endpoint.Endpoint {
	var balancer lb.Balancer
	var endpointer sd.Endpointer
	if b.FixedEndpointer != nil {
		endpointer = b.FixedEndpointer
	} else {
		endpointer = b.Endpointer
	}
	balancer = lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)
	return retry
}
