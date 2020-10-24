package svc

import (
	"context"
	"time"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	"github.com/thelotter-enterprise/usergo/shared"
)

// NewLoggingMiddleware ... ..
func NewLoggingMiddleware(logger *tlelogger.Manager) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger *tlelogger.Manager
	next   Service
}

func (mw loggingMiddleware) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	defer func(begin time.Time) {
		logger := *mw.logger
		_ = logger.Info(
			ctx,
			"GetUserByID",
			"method", "GetUserByID",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByID(ctx, userID)
}

func (mw loggingMiddleware) ConsumeLoginCommand(ctx context.Context, userID int) error {
	defer func(begin time.Time) {
		logger := *mw.logger
		_ = logger.Info(
			ctx,
			"ConsumeLoginCommand",
			"method", "ConsumeLoginCommand",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.ConsumeLoginCommand(ctx, userID)
}
