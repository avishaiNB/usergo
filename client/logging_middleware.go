package client

import (
	"time"

	"github.com/go-kit/kit/log"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
)

// NewLoggingMiddleware ...
func NewLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) GetUserByID(id int) (response tlehttp.Response) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
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
		_ = mw.logger.Log(
			"method", "GetUserByEmail",
			"input", email,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByEmail(email)
}
