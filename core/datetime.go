package core

import (
	"strconv"
	"time"
)

// DateTime ...
type DateTime struct {
}

// NewDateTime creates a new instance of the date time
func NewDateTime() DateTime {
	return DateTime{}
}

// Now returns the current date and time in UTC
func (dt DateTime) Now() time.Time {
	return time.Now().UTC()
}

// AddDuration will add duration to now date in utc
func (dt DateTime) AddDuration(duration time.Duration) time.Time {
	date := dt.Now().Add(duration)
	return date
}

// DurationToString will return the duration as string
func (dt DateTime) DurationToString(duration time.Duration) string {
	return duration.String()
}

// StringToDuration will return a time.Duration from the given duration string
func (dt DateTime) StringToDuration(duration string) (time.Duration, error) {
	return time.ParseDuration(duration)
}

// TimeToString will convert the time into unix int64, then return it as a string
func (dt DateTime) TimeToString(t time.Time) string {
	time := t.UTC().Unix()
	return strconv.FormatInt(time, 10)
}

// StringToTime will convert the string time to int64 which represents the unix time
// then it will create a new time and will return it
func (dt DateTime) StringToTime(t string) (time.Time, error) {
	timeInUnix, err := strconv.ParseInt(t, 10, 64)
	var date time.Time
	if err == nil {
		date = time.Unix(timeInUnix, 0).UTC()
	}
	return date, err
}
