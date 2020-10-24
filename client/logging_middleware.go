package client

import (
	"context"
	"time"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
)

// NewLoggingMiddleware ...
func NewLoggingMiddleware(logger *tlelogger.Manager) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger *tlelogger.Manager
	next   Service
}

func (mw loggingMiddleware) GetUserByID(id int) (response tlehttp.Response) {
	defer func(begin time.Time) {
		logger := *mw.logger
		_ = logger.Info(
			context.Background(),
			"GetUseByID middleware",
			"method", "GetUserByID",
			"input", id,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByID(id)
}

func (mw loggingMiddleware) GetUserByEmail(email string) (response tlehttp.Response) {
	defer func(begin time.Time) {
		logger := *mw.logger
		_ = logger.Info(
			context.Background(),
			"GetUseByEmail middleware",
			"method", "GetUserByEmail",
			"input", email,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByEmail(email)
}
