package errors

import (
	"context"
	"fmt"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
)

// ErrorHelper holds helper functions for working with errors
type ErrorHelper struct{}

// ApplicationError ...
func (ErrorHelper) ApplicationError(ctx context.Context, err error, msg string, format string, args ...interface{}) error {
	corrid := tlectx.GetCorrelationFromContext(ctx)
	appErr := NewApplicationError(err, msg)
	appErr = Annotatef(appErr, format, args...)
	appErr = Annotate(appErr, fmt.Sprintf("correlationid: %s", corrid))
	return appErr
}
