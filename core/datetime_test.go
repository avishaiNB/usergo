package core_test

import (
	"testing"
	"time"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestDurationToString(t *testing.T) {
	duration := time.Second * 10
	dt := core.NewDateTime()

	result := dt.DurationToString(duration)

	if result != "10s" {
		t.Error(result)
	}
}

func TestStringToDuration(t *testing.T) {
	durationExpected := time.Second * 10
	duration := "10s"
	dt := core.NewDateTime()

	result, err := dt.StringToDuration(duration)

	if err != nil {
		t.Error(err)
	}

	if result != durationExpected {
		t.Error(result)
	}
}
