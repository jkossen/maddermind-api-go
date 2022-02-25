package main

import "errors"

func ChkAttempt(g []int, c []int) ([]int, error) {
	size := len(c)

	// array for collecting return values
	ret := make([]int, size)

	// expect nr of items in guess is equal to nr of items in code
	if len(g) != size {
		return ret, errors.New("Length of guess and code is not equal")
	}

	// count for both correct numbers and correct position
	correct := 0

	// count for correct numbers in wrong position
	wrongPos := 0

	// collect which guesses were in the wrong position
	var rightPosses = make(map[int]bool)

	// collect which guesses were in the wrong position
	var wrongPosses = make(map[int]bool)

	for i, gE := range g {
		for j, cE := range c {
			if gE == cE && i == j {
				// number and position correct
				correct++

				// record that we've seen gE at the right position
				_, isRightPos := rightPosses[gE]
				if !isRightPos {
					rightPosses[gE] = true
				}

				// if gE was in some wrong position earlier but in correct position here don't count it in wrongPosses
				_, isWrongPos := wrongPosses[gE]
				if isWrongPos {
					wrongPos--
					delete(wrongPosses, gE)
				}

				// gE was found, break out of loop to prevent counting it multiple times if the value in some
				// other position as well
				break
			} else if gE == cE {
				//  number correct but position wrong
				_, isWrongPos := wrongPosses[gE]
				_, isRightPos := rightPosses[gE]

				// count wrongPosses only once even if it's tried multiple times
				if !isWrongPos && !isRightPos {
					wrongPosses[gE] = true
					wrongPos++
				}
			}
		}
	}

	var i int

	// 2: correct; 1: wrongPos
	for i = 0; i < (correct + wrongPos); i++ {
		if i < correct {
			ret[i] = 2
		} else {
			ret[i] = 1
		}
	}

	return ret, nil
}
