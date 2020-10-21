package errors_test

import (
	"context"
	"testing"

	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
)

func TestApplicationError(t *testing.T) {
	ctx := context.Background()
	msg := "some message"
	err := tleerrors.New("some error :(")
	format := ""
	helper := tleerrors.ErrorHelper{}

	appErr := helper.ApplicationError(ctx, err, msg, format)

	// if tleerrors.IsApplicationError(appErr) == false {
	// 	t.Fail()
	// }
}
