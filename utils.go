package main

import (
	"fmt"
	"math/rand"
	"time"
)

// SepEveryNth inserts string c after every nth character meant as a separator
// It returns the resulting string
func SepEveryNth(s string, n int, c string) string {
	for i := n; i < len(s); i += n + 1 {
		s = s[:i] + c + s[i:]
	}

	return s
}

// RandString generates a random string of n characters in length
// It returns the resulting string
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())

	srcChars := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	srcCnt := len(srcChars)

	dstChars := make([]rune, n)

	for i := range dstChars {
		dstChars[i] = srcChars[rand.Intn(srcCnt)]
	}

	return string(dstChars)
}

// StartOfDayEpoch determines the Unix timestamp in seconds for 00:00am this day
// It returns the resulting int64
func StartOfDayEpoch(t time.Time) int64 {
	year, month, day := t.Year(), t.Month(), t.Day()
	tt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%d-%02d-%02dT00:00:00+00:00", year, month, day))
	return tt.Unix()
}
