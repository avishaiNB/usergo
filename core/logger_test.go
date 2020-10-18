package core_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/core"
)

func TestLog(t *testing.T) {
	message := "test error message"
	err := errors.New("test error")
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stdout)
	log := core.NewLog(logger, core.LogLevelCritical)

	wasLogged := log.Error(ctx, message, err, logger)

	if wasLogged == true {
		t.Error("Should not have logged due to low log level")
	}
}
