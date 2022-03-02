package main

import (
	"errors"
	"math/rand"
	"time"
)

func GenCode(n int) []int {
	rand.Seed(time.Now().UnixNano())

	srcNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	srcCnt := len(srcNumbers)

	dstNrs := make([]int, n)

	for i := range dstNrs {
		dstNrs[i] = srcNumbers[rand.Intn(srcCnt)]
	}

	return dstNrs
}

func ChkAttempt(g []int, c []int) ([]int, error) {
	size := len(c)

	// array for collecting return values
	ret := make([]int, size)

	// expect nr of items in guess is equal to nr of items in code
	if len(g) != size {
		return ret, errors.New("Length of guess and code is not equal")
	}

	// collect which guesses were in the wrong position
	var rightPosses = make(map[int]int)

	// collect which guesses were in the wrong position
	var wrongPosses = make(map[int]bool)

	nrRightPosses := 0
	nrWrongPosses := 0

	// gVal: value of guess item
	// cVal: value of code item
	for i, gVal := range g {
		for j, cVal := range c {
			if gVal == cVal && i == j {
				rightPosses[gVal]++
				nrRightPosses++
				// It's the right color in the right pos. It will not get better than this. Break the loop.
				break
			} else if gVal == cVal && i != j {
				// do not count as wrongPos if it's also in the right pos in guess or code
				if g[j] != gVal && c[i] != gVal && !wrongPosses[j] {
					wrongPosses[j] = true
					nrWrongPosses++
					// only count this wrongpos once, so break the loop
					break
				}
			}
		}
	}

	var i int

	// 2: correct; 1: wrongPos
	for i = 0; i < (nrRightPosses + nrWrongPosses); i++ {
		if i < nrRightPosses {
			ret[i] = 2
		} else {
			ret[i] = 1
		}
	}

	return ret, nil
}
