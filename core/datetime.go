package core

import (
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

// DurationToString ...
// TODO: work with unix time???
func (dt DateTime) DurationToString(duration time.Duration) string {
	return ""
}

// StringToDuration ...
// TODO: work with unix time???
func (dt DateTime) StringToDuration(duration string) (time.Duration, error) {
	return time.ParseDuration(duration)
}

// TimeToString ...
// TODO: work with unix time???
func (dt DateTime) TimeToString(t time.Time) string {

	return ""
}

// StringToTime ...
// TODO: work with unix time???
func (dt DateTime) StringToTime(duration string) time.Time {
	return dt.Now()
}
