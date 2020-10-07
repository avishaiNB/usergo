package core

import (
	metrics "github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// this should be infra

// Instrumentor ..
type Instrumentor struct {
	ServiceName        string
	PromCounters       []metrics.Counter
	PromSummaryVectors []metrics.Histogram
	PromGauges         []metrics.Histogram
	PromHistograms     []metrics.Histogram
}

// NewInstrumentor ...
func NewInstrumentor(serviceName string) Instrumentor {
	return Instrumentor{
		ServiceName: serviceName,
	}
}

// InstrumentationInfo is a string for specifying the counter names
type InstrumentationInfo struct {
	Name string
	Help string
}

// RequestCount is "request_count"
var RequestCount = InstrumentationInfo{
	Name: "request_count",
	Help: "Number of requests received",
}

// LatencyInMili is "request_latency_microseconds"
var LatencyInMili = InstrumentationInfo{
	Name: "request_latency_microseconds",
	Help: "Total duration of requests in microseconds",
}

// AddPromCounter will add a new counter to prometheous
// Namespace, Subsystem, and Name are components of the fully-qualified name of the Metric (created by joining these components with "_").
// Only Name is mandatory, the others merely help structuring the name.
// promLabels are labels to differentiate the characteristics of the thing that is being measured
// read here for more information: https://prometheus.io/docs/practices/naming/
func (inst *Instrumentor) AddPromCounter(namespace string, subsystem string, info InstrumentationInfo, promLabels []string) metrics.Counter {
	opts := stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      info.Name,
		Help:      info.Help,
	}
	var counter metrics.Counter
	counter = kitprometheus.NewCounterFrom(opts, promLabels)
	inst.PromCounters = append(inst.PromCounters, counter)
	return counter
}

// AddPromSummary ...
func (inst *Instrumentor) AddPromSummary(namespace string, subsystem string, info InstrumentationInfo, promLabels []string) metrics.Histogram {
	opts := stdprometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      info.Name,
		Help:      info.Help,
	}

	var summary metrics.Histogram
	summary = kitprometheus.NewSummaryFrom(opts, promLabels)
	inst.PromSummaryVectors = append(inst.PromSummaryVectors, summary)
	return summary
}

// AddPromGauge ...
func (inst *Instrumentor) AddPromGauge(namespace string, subsystem string, info InstrumentationInfo, promLabels []string) {
	// TBD
}

// AddPromHistogram ...
func (inst *Instrumentor) AddPromHistogram(namespace string, subsystem string, info InstrumentationInfo, promLabels []string) {
	// TBD
}
