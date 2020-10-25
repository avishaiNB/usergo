package context

// CorrelationIDHeaderKey ..
type CorrelationIDHeaderKey string

// CorrelationIDRootHeaderKey ..
type CorrelationIDRootHeaderKey string

// DurationHeaderKey ..
type DurationHeaderKey string

// DeadlineHeaderKey ..
type DeadlineHeaderKey string

const (
	// CorrelationIDKey ...
	CorrelationIDKey CorrelationIDHeaderKey = "correlation_id"

	// CorrelationIDRootKey ..
	CorrelationIDRootKey CorrelationIDRootHeaderKey = "root_correlation_id"

	// DurationKey ..
	DurationKey DurationHeaderKey = "duration"

	// DeadlineKey ...
	DeadlineKey DeadlineHeaderKey = "deadline"
)
