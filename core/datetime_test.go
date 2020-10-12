package core_test

import (
	"testing"
)

func TestDurationToString(t *testing.T) {
	// duration := time.Second * 10
	// dt := core.NewDateTime()

	// result := dt.DurationToString(duration)

	// if result != "10s" {
	// 	t.Error(result)
	// }
}

// func TestStringToDuration(t *testing.T) {
// 	durationExpected := time.Second * 10
// 	duration := "10s"
// 	dt := core.NewDateTime()

// 	result, err := dt.StringToDuration(duration)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if result != durationExpected {
// 		t.Error(result)
// 	}
// }

// func TestTimeToStringAndStringToTimeUTC(t *testing.T) {
// 	date := time.Date(2020, 10, 16, 9, 30, 59, 0, time.UTC)
// 	dt := core.NewDateTime()

// 	dateAsString := dt.TimeToString(date)
// 	dateAsTime, err := dt.StringToTime(dateAsString)

// 	if dateAsTime != date {
// 		t.Error(err)
// 	}
// }

// func TestTimeToStringAndStringToTimeNotUTC(t *testing.T) {
// 	date := time.Date(2020, 10, 16, 9, 30, 59, 0, time.Local)
// 	dt := core.NewDateTime()

// 	dateAsString := dt.TimeToString(date)
// 	dateAsTimeInUTC, err := dt.StringToTime(dateAsString)
// 	dateAsTimeLocal := dateAsTimeInUTC.Local()

// 	if dateAsTimeLocal != date {
// 		t.Error(err)
// 	}
// }
