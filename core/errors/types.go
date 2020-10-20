package errors

import (
	jujuerr "github.com/juju/errors"
)

// NewUnauthorized returns an error which wraps err and satisfies IsUnauthorized().
func NewUnauthorized(err error, msg string) error {
	return jujuerr.NewUnauthorized(err, msg)
}

// NewUnauthorizedf returns an error which satisfies IsUnauthorized().
func NewUnauthorizedf(format string, args ...interface{}) error {
	return jujuerr.Unauthorizedf(format, args)
}

func IsUnauthorized(err error) bool {
	return jujuerr.IsUnauthorized(err)
}

// NewNotValidf returns an error which satisfies IsNotValid().
func NewNotValidf(format string, args ...interface{}) error {
	return jujuerr.NotValidf(format, args)
}

// NewNotValid returns an error which wraps err and satisfies IsNotValid().
func NewNotValid(err error, msg string) error {
	return jujuerr.NewNotValid(err, msg)
}

func IsNotValid(err error) bool {
	return jujuerr.IsNotValid(err)
}

func NewNotSupportedf(format string, args ...interface{}) error {
	return jujuerr.NotSupportedf(format, args)
}

func NewNotSupported(err error, msg string) error {
	return jujuerr.NewNotSupported(err, msg)
}

func IsNotSupported(err error) bool {
	return jujuerr.IsNotSupported(err)
}

func NewBadRequestf(format string, args ...interface{}) error {
	return jujuerr.BadRequestf(format, args)
}

func NewBadRequest(err error, msg string) error {
	return jujuerr.NewBadRequest(err, msg)
}

func IsBadRequest(err error) bool {
	return jujuerr.IsBadRequest(err)
}

func NewForbiddenf(format string, args ...interface{}) error {
	return jujuerr.Forbiddenf(format, args)
}

func NewForbidden(err error, msg string) error {
	return jujuerr.NewForbidden(err, msg)
}

func IsForbidden(err error) bool {
	return jujuerr.IsForbidden(err)
}

func NewMethodNotAllowedf(format string, args ...interface{}) error {
	return jujuerr.MethodNotAllowedf(format, args)
}

func NewMethodNotAllowed(err error, msg string) error {
	return jujuerr.NewMethodNotAllowed(err, msg)
}

func IsMethodNotAllowed(err error) bool {
	return jujuerr.IsMethodNotAllowed(err)
}

func NewNotFoundf(format string, args ...interface{}) error {
	return jujuerr.NotFoundf(format, args)
}

func NewNotFound(err error, msg string) error {
	return jujuerr.NewNotFound(err, msg)
}

func IsNotFound(err error) bool {
	return jujuerr.IsNotFound(err)
}

func NewTimeoutf(format string, args ...interface{}) error {
	return jujuerr.Timeoutf(format, args)
}

func NewTimeout(err error, msg string) error {
	return jujuerr.NewTimeout(err, msg)
}

func IsTimeout(err error) bool {
	return jujuerr.IsTimeout(err)
}

func NewNotImplementedf(format string, args ...interface{}) error {
	return jujuerr.NotImplementedf(format, args)
}

func NewNotImplemented(err error, msg string) error {
	return jujuerr.NewNotImplemented(err, msg)
}

func IsNotImplemented(err error) bool {
	return jujuerr.IsNotImplemented(err)
}

func NewUserNotFoundf(format string, args ...interface{}) error {
	return jujuerr.UserNotFoundf(format, args)
}

func NewUserNotFound(err error, msg string) error {
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
