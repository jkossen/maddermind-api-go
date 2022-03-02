package main

import (
	"testing"
	"time"
)

func TestStartOfDayEpoch(t *testing.T) {
	testTime := time.Date(2022, 03, 02, 20, 34, 58, 651387237, time.UTC)
	res := StartOfDayEpoch(testTime)
	var want int64 = 1646179200

	if res != want {
		t.Fatalf(`StartOfDayEpoch(%v): %v != %v`, testTime, res, want)
	}

	testTime = time.Date(2022, 03, 01, 1, 02, 28, 651387237, time.UTC)
	res = StartOfDayEpoch(testTime)
	want = 1646092800

	if res != want {
		t.Fatalf(`StartOfDayEpoch(%v): %v != %v`, testTime, res, want)
	}

	testTime = time.Date(2015, 12, 29, 14, 42, 59, 0, time.UTC)
	res = StartOfDayEpoch(testTime)
	want = 1451347200

	if res != want {
		t.Fatalf(`StartOfDayEpoch(%v): %v != %v`, testTime, res, want)
	}

	testTime = time.Date(2043, 8, 9, 23, 28, 0, 0, time.UTC)
	res = StartOfDayEpoch(testTime)
	want = 2322691200

	if res != want {
		t.Fatalf(`StartOfDayEpoch(%v): %v != %v`, testTime, res, want)
	}

}
