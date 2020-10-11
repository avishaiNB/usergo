package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestNewFrom(t *testing.T) {
	ctx := context.Background()
	c := core.NewCtx()
	corrid := "12345"
	duration := time.Second * 10
	deadline := time.Now().UTC().Add(duration)

	c = c.NewFrom(ctx, corrid, duration, deadline)

	if c.Cancel == nil {
		t.Fail()
	}

	if c.CorrelationID != corrid {
		t.Error("correlation does not match")
	}

	if c.Duration != duration {
		t.Error("duration does not match")
	}

	if c.Deadline != deadline {
		t.Error("deadline does not match")
	}

	actualDeadline, _ := c.Context.Deadline()
	if actualDeadline != deadline {
		t.Error("context deadline does not match")
	}
}

func TestSetContext(t *testing.T) {
	ctx := context.Background()
	c := core.NewCtx()
	corrid := "12345"
	duration := time.Second * 10
	deadline := time.Now().UTC().Add(duration)

	ctx = c.SetCorrealtionToContext(ctx, corrid)
	corridResult := c.GetCorrelationFromContext(ctx)

	if corrid != corridResult {
		t.Error("correlation does not match")
	}

	ctx = c.SetTimeoutToContext(ctx, duration, deadline)
	durationResult, deadlineResult := c.GetTimeoutFromContext(ctx)

	if durationResult != duration {
		t.Error("duration does not match")
	}

	if deadlineResult != deadline {
		t.Error("deadline does not match")
	}
}

func TestGetOrCreateCorrelationID_Create(t *testing.T) {
	ctx := context.Background()
	c := core.NewCtx()
	var corrid string
	corrid, ctx = c.GetOrCreateCorrelationFromContext(ctx, true)
	actualCorrelationID := c.GetCorrelationFromContext(ctx)

	if actualCorrelationID == "" {
		t.Error("correlation was not set to context")
	}

	if corrid != actualCorrelationID {
		t.Error("correlation does not match")
	}
}

func TestGetOrCreateTimeout_Create(t *testing.T) {
	ctx := context.Background()
	c := core.NewCtx()

	duration, deadline, ctx := c.GetOrCreateTimeoutFromContext(ctx, true)
	actualDuration, actualDeadlne := c.GetTimeoutFromContext(ctx)

	if actualDuration != duration {
		t.Error("duration does not match")
	}

	if deadline != actualDeadlne {
		t.Error("deadline does not match")
	}
}
