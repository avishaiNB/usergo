package context

import (
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// CorrelationID used to create correlation ID
type CorrelationID interface {
	New() string
}

type correlationid struct{}

// NewCorrelationID creates a new instance of the correlationID interface
func NewCorrelationID() CorrelationID {
	return correlationid{}
}

// NewCorrelation will return a new correlation ID as string
func (correlationid) New() string {
	return utils.NewUUID()
}
