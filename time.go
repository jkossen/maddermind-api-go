package main

import (
	"fmt"
	"time"
)

// StartOfDayEpoch determines the Unix timestamp in seconds for 00:00am this day
// It returns the resulting int64
func StartOfDayEpoch(t time.Time) int64 {
	year, month, day := t.Year(), t.Month(), t.Day()
	tt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%d-%02d-%02dT00:00:00+00:00", year, month, day))
	return tt.Unix()
}
