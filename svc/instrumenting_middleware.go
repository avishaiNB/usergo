package svc

import (
	"context"

	metrics "github.com/go-kit/kit/metrics"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlemetrics "github.com/thelotter-enterprise/usergo/core/metrics"
	"github.com/thelotter-enterprise/usergo/shared"
)

// NewInstrumentingMiddleware ..
func NewInstrumentingMiddleware(logger tlelogger.Manager, inst tlemetrics.PrometheusInstrumentor) ServiceMiddleware {
	counter := inst.AddPromCounter("user", "getuserbyid", tlemetrics.RequestCount, []string{"method", "error"})
	requestLatency := inst.AddPromSummary("user", "getuserbyid", tlemetrics.LatencyInMili, []string{"method", "error"})

	return func(next Service) Service {
		mw := instrumentingMiddleware{
			next:           next,
			requestCount:   counter,
			requestLatency: requestLatency,
			logger:         logger,
		}
		return mw
	}
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
	logger         tlelogger.Manager
}

func (mw instrumentingMiddleware) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	return mw.next.GetUserByID(ctx, userID)
}

func (mw instrumentingMiddleware) ConsumeLoginCommand(ctx context.Context, userID int) error {
	return mw.next.ConsumeLoginCommand(ctx, userID)
}
