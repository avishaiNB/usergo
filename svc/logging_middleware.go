package svc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/shared"
)

// NewLoggingMiddleware ... ..
func NewLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetUserByID",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByID(ctx, userID)
}

func (mw loggingMiddleware) ConsumeLoginCommand(ctx context.Context, userID int) error {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ConsumeLoginCommand",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.ConsumeLoginCommand(ctx, userID)
}
