package manager_test

import (
	"context"
	"testing"
	"time"

	tlectx "github.com/thelotter-enterprise/usergo/core/context/manager"
)

func TestNewFrom(t *testing.T) {
	ctx := context.Background()
	corrid := "12345"
	duration := time.Second * 10
	deadline := time.Now().UTC().Add(duration)

	ctx = tlectx.SetCorrealtion(ctx, correlationID)
	ctx = tlectx.SetTimeout(ctx, duration, deadline)
	ctx, cancel = context.WithDeadline(ctx, deadline)

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
	corrid := "12345"
	duration := time.Second * 10
	deadline := time.Now().UTC().Add(duration)

	ctx = tlectx.SetCorrealtion(ctx, corrid)
	corridResult := tlectx.GetCorrelationID(ctx)

	if corrid != corridResult {
		t.Error("correlation does not match")
	}

	ctx = tlectx.SetTimeout(ctx, duration, deadline)
	durationResult, deadlineResult := tlectx.GetTimeout(ctx)

	if durationResult != duration {
		t.Error("duration does not match")
	}

	if deadlineResult != deadline {
		t.Error("deadline does not match")
	}
}

func TestGetOrCreateCorrelationID_Create(t *testing.T) {
	ctx := context.Background()
	var corrid string
	corrid, ctx = tlectx.GetOrCreateCorrelationFromContext(ctx, true)
	actualCorrelationID := tlectx.GetCorrelationID(ctx)

	if actualCorrelationID == "" {
		t.Error("correlation was not set to context")
	}

	if corrid != actualCorrelationID {
		t.Error("correlation does not match")
	}
}

func TestGetOrCreateTimeout_Create(t *testing.T) {
	ctx := context.Background()

	duration, deadline, ctx := tlectx.GetOrCreateTimeoutFromContext(ctx, true)
	actualDuration, actualDeadlne := tlectx.GetTimeout(ctx)

	if actualDuration != duration {
		t.Error("duration does not match")
	}

	if deadline != actualDeadlne {
		t.Error("deadline does not match")
	}
}
