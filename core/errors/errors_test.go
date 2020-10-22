package errors_test

import (
	"testing"

	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
)

func TestApplicationError(t *testing.T) {
	format := "error from %s %s"
	err := tleerrors.NewApplicationErrorf(format, "guy", "kolbis")
	isErrString := err.Error()
	wantErrString := "error from guy kolbis application error"

	if tleerrors.IsApplicationError(err) == false {
		t.Fail()
	}

	if wantErrString != isErrString {
		t.Errorf("wanted %s, got %s", wantErrString, isErrString)
	}
}

func TestErrorWithAnnotation(t *testing.T) {
	msg := "original error from guy kolbis"
	err := tleerrors.New(msg)
	err = tleerrors.Annotate(err, "more information about the error")

	isErrString := err.Error()
	wantErrString := "more information about the error: original error from guy kolbis"

	if wantErrString != isErrString {
		t.Errorf("wanted %s, got %s", wantErrString, isErrString)
	}
}

func TestErrorWithWrap(t *testing.T) {
	msg := "original error from guy kolbis"
	err := tleerrors.New(msg)
	errForbidden := tleerrors.NewForbiddenError(err, "forbidden!")
	newerr := tleerrors.Wrap(err, errForbidden)

	if tleerrors.IsForbidden(newerr) == false {
		t.Fail()
	}

	isErrString := newerr.Error()
	wantErrString := "forbidden!: original error from guy kolbis"

	if wantErrString != isErrString {
		t.Errorf("wanted %s, got %s", wantErrString, isErrString)
	}
}
