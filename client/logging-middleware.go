package client

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/core"
)

func makeLoggingMiddleware(logger log.Logger) UserServiceMiddleware {
	return func(next UserService) UserService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	UserService
}

func (mw loggingMiddleware) GetUserByID(id int) (response core.Response) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetUserByID",
			"input", id,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.UserService.GetUserByID(id)
}

func (mw loggingMiddleware) GetUserByEmail(email string) (response core.Response) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetUserByEmail",
			"input", email,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.UserService.GetUserByEmail(email)
}
