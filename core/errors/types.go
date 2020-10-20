package errors

import (
	jujuerr "github.com/juju/errors"
)

// NewUnauthorized returns an error which wraps err and satisfies IsUnauthorized().
func NewUnauthorizedError(err error, msg string) error {
	return jujuerr.NewUnauthorized(err, msg)
}

// NewUnauthorizedf returns an error which satisfies IsUnauthorized().
func NewUnauthorizedErrorf(format string, args ...interface{}) error {
	return jujuerr.Unauthorizedf(format, args)
}

func IsUnauthorized(err error) bool {
	return jujuerr.IsUnauthorized(err)
}

// NewNotValidf returns an error which satisfies IsNotValid().
func NewNotValidErrorf(format string, args ...interface{}) error {
	return jujuerr.NotValidf(format, args)
}

// NewNotValid returns an error which wraps err and satisfies IsNotValid().
func NewNotValidError(err error, msg string) error {
	return jujuerr.NewNotValid(err, msg)
}

func IsNotValid(err error) bool {
	return jujuerr.IsNotValid(err)
}

func NewNotSupportedErrorf(format string, args ...interface{}) error {
	return jujuerr.NotSupportedf(format, args)
}

func NewNotSupportedError(err error, msg string) error {
	return jujuerr.NewNotSupported(err, msg)
}

func IsNotSupported(err error) bool {
	return jujuerr.IsNotSupported(err)
}

func NewBadRequestErrorf(format string, args ...interface{}) error {
	return jujuerr.BadRequestf(format, args)
}

func NewBadRequestError(err error, msg string) error {
	return jujuerr.NewBadRequest(err, msg)
}

func IsBadRequest(err error) bool {
	return jujuerr.IsBadRequest(err)
}

func NewForbiddenErrorf(format string, args ...interface{}) error {
	return jujuerr.Forbiddenf(format, args)
}

func NewForbiddenError(err error, msg string) error {
	return jujuerr.NewForbidden(err, msg)
}

func IsForbidden(err error) bool {
	return jujuerr.IsForbidden(err)
}

func NewMethodNotAllowedErrorf(format string, args ...interface{}) error {
	return jujuerr.MethodNotAllowedf(format, args)
}

func NewMethodNotAllowedError(err error, msg string) error {
	return jujuerr.NewMethodNotAllowed(err, msg)
}

func IsMethodNotAllowed(err error) bool {
	return jujuerr.IsMethodNotAllowed(err)
}

func NewNotFoundErrorf(format string, args ...interface{}) error {
	return jujuerr.NotFoundf(format, args)
}

func NewNotFoundError(err error, msg string) error {
	return jujuerr.NewNotFound(err, msg)
}

func IsNotFound(err error) bool {
	return jujuerr.IsNotFound(err)
}

func NewTimeoutErrorf(format string, args ...interface{}) error {
	return jujuerr.Timeoutf(format, args)
}

func NewTimeoutError(err error, msg string) error {
	return jujuerr.NewTimeout(err, msg)
}

func IsTimeout(err error) bool {
	return jujuerr.IsTimeout(err)
}

func NewNotImplementedErrorf(format string, args ...interface{}) error {
	return jujuerr.NotImplementedf(format, args)
}

func NewNotImplementedError(err error, msg string) error {
	return jujuerr.NewNotImplemented(err, msg)
}

func IsNotImplemented(err error) bool {
	return jujuerr.IsNotImplemented(err)
}

func NewUserNotFoundErrorf(format string, args ...interface{}) error {
	return jujuerr.UserNotFoundf(format, args)
}

func NewUserNotFoundError(err error, msg string) error {
	return jujuerr.NewUserNotFound(err, msg)
}

func IsUserNotFound(err error) bool {
	return jujuerr.IsUserNotFound(err)
}

type ApplicationError struct {
	jujuerr.Err
}

func NewApplicationError(msg string) error {
	err := &ApplicationError{jujuerr.NewErr("%s application error", msg)}
	err.SetLocation(1)
	return err
}
