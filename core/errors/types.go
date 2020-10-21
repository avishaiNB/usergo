package errors

import (
	jujuerr "github.com/juju/errors"
)

// NewUnauthorizedError returns an error which wraps err and satisfies IsUnauthorized().
func NewUnauthorizedError(err error, msg string) error {
	return jujuerr.NewUnauthorized(err, msg)
}

// NewUnauthorizedErrorf returns an error which satisfies IsUnauthorized().
func NewUnauthorizedErrorf(format string, args ...interface{}) error {
	return jujuerr.Unauthorizedf(format, args)
}

// IsUnauthorized reports whether err was created with Unauthorizedf() or NewUnauthorized().
func IsUnauthorized(err error) bool {
	return jujuerr.IsUnauthorized(err)
}

// NewNotValidErrorf returns an error which satisfies IsNotValid().
func NewNotValidErrorf(format string, args ...interface{}) error {
	return jujuerr.NotValidf(format, args)
}

// NewNotValidError returns an error which wraps err and satisfies IsNotValid().
func NewNotValidError(err error, msg string) error {
	return jujuerr.NewNotValid(err, msg)
}

// IsNotValid reports whether the error was created with NotValidf() or NewNotValid().
func IsNotValid(err error) bool {
	return jujuerr.IsNotValid(err)
}

// NewNotSupportedErrorf returns an error which satisfies IsNotSupported().
func NewNotSupportedErrorf(format string, args ...interface{}) error {
	return jujuerr.NotSupportedf(format, args)
}

// NewNotSupportedError returns an error which wraps err and satisfies IsNotSupported().
func NewNotSupportedError(err error, msg string) error {
	return jujuerr.NewNotSupported(err, msg)
}

// IsNotSupported reports whether the error was created with NotSupportedf() or NewNotSupported().
func IsNotSupported(err error) bool {
	return jujuerr.IsNotSupported(err)
}

// NewBadRequestErrorf returns an error which satisfies IsBadRequest().
func NewBadRequestErrorf(format string, args ...interface{}) error {
	return jujuerr.BadRequestf(format, args)
}

// NewBadRequestError returns an error which wraps err that satisfies IsBadRequest().
func NewBadRequestError(err error, msg string) error {
	return jujuerr.NewBadRequest(err, msg)
}

// IsBadRequest reports whether err was created with BadRequestf() or NewBadRequest().
func IsBadRequest(err error) bool {
	return jujuerr.IsBadRequest(err)
}

// NewForbiddenErrorf returns an error which satistifes IsForbidden()
func NewForbiddenErrorf(format string, args ...interface{}) error {
	return jujuerr.Forbiddenf(format, args)
}

// NewForbiddenError returns an error which wraps err that satisfies IsForbidden().
func NewForbiddenError(err error, msg string) error {
	return jujuerr.NewForbidden(err, msg)
}

// IsForbidden reports whether err was created with Forbiddenf() or NewForbidden().
func IsForbidden(err error) bool {
	return jujuerr.IsForbidden(err)
}

// NewMethodNotAllowedErrorf returns an error which satisfies IsMethodNotAllowed().
func NewMethodNotAllowedErrorf(format string, args ...interface{}) error {
	return jujuerr.MethodNotAllowedf(format, args)
}

// NewMethodNotAllowedError returns an error which wraps err that satisfies IsMethodNotAllowed().
func NewMethodNotAllowedError(err error, msg string) error {
	return jujuerr.NewMethodNotAllowed(err, msg)
}

// IsMethodNotAllowed reports whether err was created with MethodNotAllowedf() or NewMethodNotAllowed().
func IsMethodNotAllowed(err error) bool {
	return jujuerr.IsMethodNotAllowed(err)
}

// NewNotFoundErrorf returns an error which satisfies IsNotFound().
func NewNotFoundErrorf(format string, args ...interface{}) error {
	return jujuerr.NotFoundf(format, args)
}

// NewNotFoundError returns an error which wraps err that satisfies IsNotFound().
func NewNotFoundError(err error, msg string) error {
	return jujuerr.NewNotFound(err, msg)
}

// IsNotFound reports whether err was created with NotFoundf() or NewNotFound().
func IsNotFound(err error) bool {
	return jujuerr.IsNotFound(err)
}

// NewTimeoutErrorf returns an error which satisfies IsTimeout().
func NewTimeoutErrorf(format string, args ...interface{}) error {
	return jujuerr.Timeoutf(format, args)
}

// NewTimeoutError returns an error which wraps err that satisfies IsTimeout().
func NewTimeoutError(err error, msg string) error {
	return jujuerr.NewTimeout(err, msg)
}

// IsTimeout reports whether err was created with Timeout() or NewTimeout().
func IsTimeout(err error) bool {
	return jujuerr.IsTimeout(err)
}

// NewNotImplementedErrorf returns an error which satisfies IsNotImplemented().
func NewNotImplementedErrorf(format string, args ...interface{}) error {
	return jujuerr.NotImplementedf(format, args)
}

// NewNotImplementedError returns an error which wraps err and satisfies IsNotImplemented().
func NewNotImplementedError(err error, msg string) error {
	return jujuerr.NewNotImplemented(err, msg)
}

// IsNotImplemented reports whether err was created with NotImplementedf() or NewNotImplemented().
func IsNotImplemented(err error) bool {
	return jujuerr.IsNotImplemented(err)
}

// ApplicationError ...
// type ApplicationError struct {
// 	Err           jujuerr.Err
// 	Message       string
// 	Context       context.Context
// 	CorrelationID string
// 	Deadline      time.Time
// 	Duration      time.Duration
// 	Args          []interface{}
// }

// // NewApplicationError return an application error
// func NewApplicationError(ctx context.Context, msg string, args ...interface{}) error {
// 	corrid := tlectx.GetCorrelationFromContext(ctx)
// 	duration, deadline := tlectx.CalcTimeoutFromContext(ctx)

// 	err := &ApplicationError{
// 		Err:           jujuerr.NewErr("%s application error", msg),
// 		Message:       msg,
// 		Context:       ctx,
// 		Args:          args,
// 		CorrelationID: corrid,
// 		Duration:      duration,
// 		Deadline:      deadline,
// 	}
// 	err.SetLocation(1)
// 	return err
//}
