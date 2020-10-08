package core

import (
	"context"
	"time"
)

const (
	// MaxTimeout is 15 seconds
	MaxTimeout time.Duration = time.Second * 15
)

// Ctx ...
type Ctx struct {
	Context context.Context
}

// NewCtx will create a new context
func NewCtx(ctx context.Context) Ctx {
	return Ctx{
		Context: ctx,
	}
}

// CalcTimeout will return the timeout (deadline) for waiting an external response to come back
// TODO: now I return max, need to change it
func (ctx Ctx) CalcTimeout() time.Duration {
	return MaxTimeout
}
