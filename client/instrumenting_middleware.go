package client

import (
	"fmt"
	"time"

	metrics "github.com/go-kit/kit/metrics"
	tleinst "github.com/thelotter-enterprise/usergo/core/metrics"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transport/http"
)

// NewInstrumentingMiddleware ...
func NewInstrumentingMiddleware(inst tleinst.PrometheusInstrumentor) ServiceMiddleware {
	counter := inst.AddPromCounter("user", "getuserbyid", tleinst.RequestCount, []string{"method", "error"})
	requestLatency := inst.AddPromSummary("user", "getuserbyid", tleinst.LatencyInMili, []string{"method", "error"})

	return func(next Service) Service {
		mw := instrumentingMiddleware{
			next:           next,
			requestCount:   counter,
			requestLatency: requestLatency,
		}
		return mw
	}
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (mw instrumentingMiddleware) GetUserByID(id int) (response tlehttp.Response) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserByID", "error", fmt.Sprint(response.Error != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response = mw.next.GetUserByID(id)
	return response
}

func (mw instrumentingMiddleware) GetUserByEmail(email string) (response tlehttp.Response) {
	return mw.next.GetUserByEmail(email)
}
