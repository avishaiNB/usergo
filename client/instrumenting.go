package client

import (
	"fmt"
	"time"

	metrics "github.com/go-kit/kit/metrics"
	"github.com/thelotter-enterprise/usergo/core"
)

func makeInstrumentingMiddleware(inst core.Instrumentor, serviceName string, api string) UserServiceMiddleware {
	counter := inst.AddPromCounter(serviceName, "getuserbyid", core.RequestCount, []string{"method", "error"})
	requestLatency := inst.AddPromSummary(serviceName, "getuserbyid", core.LatencyInMili, []string{"method", "error"})

	return func(next UserService) UserService {
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
	next           UserService
}

func (mw instrumentingMiddleware) GetUserByID(id int) (response core.HTTPResponse) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserByID", "error", fmt.Sprint(response.Error != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response = mw.next.GetUserByID(id)
	return response
}

func (mw instrumentingMiddleware) GetUserByEmail(email string) (response core.HTTPResponse) {
	return mw.next.GetUserByEmail(email)
}
