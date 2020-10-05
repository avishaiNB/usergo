package client

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/thelotter-enterprise/usergo/shared"
)

func makeInstrumentingMiddleware(service string, api string) UserServiceClientMiddleware {

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: service,
		Subsystem: api,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: service,
		Subsystem: api,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: service,
		Subsystem: api,
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	return func(next UserServiceClient) UserServiceClient {
		return instrumentingMiddleware{requestCount, requestLatency, countResult, next}
	}
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	UserServiceClient
}

func (mw instrumentingMiddleware) GetUserByID(id int) (response shared.HTTPResponse) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserByID", "error", fmt.Sprint(response.Error != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response = mw.UserServiceClient.GetUserByID(id)
	return response
}

func (mw instrumentingMiddleware) GetUserByEmail(email string) (response shared.HTTPResponse) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserByEmail", "error", fmt.Sprint(response.Error != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response = mw.UserServiceClient.GetUserByEmail(email)
	return response
}
