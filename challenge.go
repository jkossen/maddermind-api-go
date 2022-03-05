package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Challenge interface {
	Code() []int
	RetrieveOrGen(cs ChallengeStorage, timestamp int64, codeLen int) (challenge, error)
	Gen(n int) challenge
	Check(guess []int) ([]int, error)
	WithCode(code []int) challenge
	FromJson(cStr string) (challenge, error)
	ToString() string
}

type ChallengeStorage interface {
	ErrNoChal() error
	DSN(dsn string)
	Open() error
	Close() error
	Challenge(time int64, len int) (string, error)
	Create(time int64, len int, code string) error
}

type challenge struct {
	code []int
}

func (c challenge) Code() []int {
	return c.code
}

func (c challenge) RetrieveOrGen(cs ChallengeStorage, timestamp int64, codeLen int) (challenge, error) {
	cStr, err := cs.Challenge(timestamp, codeLen)

	if cStr == "" || err == cs.ErrNoChal() {
		c = c.Gen(codeLen)
		err = cs.Create(timestamp, codeLen, c.ToString())

		return c, err
	}

	return c.FromJson(cStr)
}

func (c challenge) Gen(n int) challenge {
	rand.Seed(time.Now().UnixNano())

	srcNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	srcCnt := len(srcNumbers)

	dstNrs := make([]int, n)

	for i := range dstNrs {
		dstNrs[i] = srcNumbers[rand.Intn(srcCnt)]
	}

	c.code = dstNrs

	return c
}

func (c challenge) Check(guess []int) ([]int, error) {
	size := len(c.code)

	// array for collecting return values
	ret := make([]int, size)

	// expect nr of items in guess is equal to nr of items in challenge
	if len(guess) != size {
		return ret, errors.New("Length of guess and challenge is not equal")
	}

	// collect which guesses were in the wrong position
	var rightPosses = make(map[int]int)

	// collect which guesses were in the wrong position
	var wrongPosses = make(map[int]bool)

	nrRightPosses := 0
	nrWrongPosses := 0

	// gVal: value of guess item
	// cVal: value of challenge item
	for i, gVal := range guess {
		for j, cVal := range c.code {
			if gVal == cVal && i == j {
				rightPosses[gVal]++
				nrRightPosses++
				// It's the right color in the right pos. It will not get better than this. Break the loop.
				break
			} else if gVal == cVal && i != j {
				// do not count as wrongPos if it's also in the right pos in guess or challenge
				if guess[j] != gVal && c.code[i] != gVal && !wrongPosses[j] {
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

func (c challenge) WithCode(code []int) challenge {
	c.code = code

	return c
}

func (c challenge) FromJson(cStr string) (challenge, error) {
	err := json.Unmarshal([]byte(cStr), (&c.code))

	return c, err
}

func (c challenge) ToString() string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(c.code)), ","), "")
}
