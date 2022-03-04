package main

import (
	"reflect"
	"testing"
)

func TestGenCode(t *testing.T) {
	var chal Challenge = challenge{}

	// test a couple of times that we don't get the same chal twice
	for n := 0; n <= 10; n++ {
		c1 := chal.Gen(6)
		c2 := chal.Gen(6)

		if reflect.DeepEqual(c1.Code(), c2.Code()) {
			t.Fatalf("GenCode(6): generated same chal twice")
		}
	}

	// test that we get a chal of the specified length
	for n := 0; n <= 10; n++ {
		c := chal.Gen(n)
		if len(c.Code()) != n {
			t.Fatalf(`GenCode(%v): outputted length != %v`, n, n)
		}
	}

}

func TestChkAttempt(t *testing.T) {
	var chal Challenge = challenge{}

	var knowngoods = [][][]int{
		{[]int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{2, 2, 2, 2}},
		{[]int{0, 0, 0, 0}, []int{0, 0, 0, 1}, []int{2, 2, 2, 0}},
		{[]int{0, 0, 0, 0}, []int{0, 0, 1, 1}, []int{2, 2, 0, 0}},
		{[]int{0, 0, 0, 0}, []int{0, 1, 1, 1}, []int{2, 0, 0, 0}},
		{[]int{0, 0, 0, 0}, []int{1, 1, 1, 1}, []int{0, 0, 0, 0}},
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 1}, []int{1, 1, 1, 1}},
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 0}, []int{1, 1, 1, 0}},
		{[]int{1, 2, 3, 4}, []int{4, 3, 0, 0}, []int{1, 1, 0, 0}},
		{[]int{1, 2, 3, 4}, []int{4, 0, 0, 0}, []int{1, 0, 0, 0}},
		{[]int{1, 2, 3, 4}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}},
		{[]int{1, 2, 3, 4}, []int{1, 2, 4, 3}, []int{2, 2, 1, 1}},
		{[]int{1, 2, 3, 4}, []int{1, 4, 2, 3}, []int{2, 1, 1, 1}},
		{[]int{5, 5, 1, 0}, []int{5, 5, 5, 5}, []int{2, 2, 0, 0}},
		{[]int{5, 5, 1, 0}, []int{1, 1, 5, 5}, []int{1, 1, 1, 0}},
		{[]int{5, 5, 1, 1}, []int{1, 1, 5, 5}, []int{1, 1, 1, 1}},
		{[]int{5, 9, 4, 0}, []int{5, 0, 4, 9}, []int{2, 2, 1, 1}},
		{[]int{7, 6, 7, 8}, []int{7, 0, 0, 0}, []int{2, 0, 0, 0}},
		{[]int{7, 6, 7, 8}, []int{7, 0, 0, 6}, []int{2, 1, 0, 0}},
		{[]int{7, 6, 7, 8}, []int{7, 6, 0, 0}, []int{2, 2, 0, 0}},
		{[]int{7, 6, 7, 8}, []int{7, 6, 0, 6}, []int{2, 2, 0, 0}},
		{[]int{7, 6, 7, 8}, []int{7, 6, 6, 6}, []int{2, 2, 0, 0}},
		{[]int{7, 6, 7, 8}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}},
		{[]int{7, 6, 7, 8}, []int{7, 6, 7, 8}, []int{2, 2, 2, 2}},
		{[]int{7, 6, 7, 8}, []int{6, 6, 8, 7}, []int{2, 1, 1, 0}},
		{[]int{7, 6, 7, 8}, []int{6, 7, 7, 8}, []int{2, 2, 1, 1}},
	}

	for _, row := range knowngoods {
		var code = row[0]
		var guess = row[1]
		var want = row[2]

		chal = chal.WithCode(code)
		res, _ := chal.Check(guess)

		if !reflect.DeepEqual(res, want) {
			t.Fatalf(`ChkAttempt(%v, %v): %v != %v`, guess, code, res, want)
		}
	}
}
