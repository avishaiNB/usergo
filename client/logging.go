package client

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/shared"
)

func makeLoggingMiddleware(logger log.Logger) UserServiceClientMiddleware {
	return func(next UserServiceClient) UserServiceClient {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	UserServiceClient
}

func (mw loggingMiddleware) GetUserByID(id int) (response shared.HTTPResponse) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetUserByID",
			"input", id,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.UserServiceClient.GetUserByID(id)
}

func (mw loggingMiddleware) GetUserByEmail(email string) (response shared.HTTPResponse) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetUserByEmail",
			"input", email,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.UserServiceClient.GetUserByEmail(email)
}
