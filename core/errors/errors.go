package errors

import (
	jujuerr "github.com/juju/errors"
)

func NewError(msg string) error {
	return jujuerr.New(msg)
}

func Annotate(err error, msg string) error {
	return jujuerr.Annotate(err, msg)
}

func Annotatef(err error, format string, args ...interface{}) error {
	return jujuerr.Annotatef(err, format, args)
}

func Cause(err error) error {
	return jujuerr.Cause(err)
}

func Details(err error) string {
	return jujuerr.Details(err)
}

func ErrorStack(err error) string {
	return jujuerr.ErrorStack(err)
}

func Errorf(format string, args ...interface{}) error {
	return jujuerr.Errorf(format, args)
}

func Mask(err error) error {
	return jujuerr.Mask(err)
}

func Maskf(err error, format string, args ...interface{}) error {
	return jujuerr.Maskf(err, format, args)
}

func Wrap(err error, newDescriptive error) error {
	return jujuerr.Wrap(err, newDescriptive)
}

func Wrapf(other error, newDescriptive error, format string, args ...interface{}) error {
	return jujuerr.Wrapf(other, newDescriptive, format, args)
}

// Trace adds the location of the Trace call to the stack.
// The Cause of the resulting error is the same as the error parameter.
// If the other error is nil, the result will be nil.
func Trace(err error) error {
	return jujuerr.Trace(err)
}
