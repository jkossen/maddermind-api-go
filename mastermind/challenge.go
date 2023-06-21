// Package mastermind provides the algorithms and functions for the game Mastermind.
package mastermind

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// challenge is the interface that contains the methods for a Mastermind challenge.
type challenge interface {
	Code() []int
	RetrieveOrGen(cs ChallengeStorage, timestamp int64, codeLen int) (challenge, error)
	Gen(n int) challenge
	Check(guess []int) ([]int, error)
	WithCode(code []int) challenge
	FromJson(cStr string) (challenge, error)
	ToString() string
}

// ChallengeStorage is the interface for a storage layer such as a database.
type ChallengeStorage interface {
	ErrNoChal() error
	DSN(dsn string)
	Open() error
	Close() error
	Challenge(time int64, len int) (string, error)
	Create(time int64, len int, code string) error
}

// Challenge is the struct wrapping data needed for a mastermind challenge
// The main part is the code that the player has to guess
type Challenge struct {
	code []int
}

// Code returns the code that the player has to guess
func (c Challenge) Code() []int {
	return c.code
}

// GetOrCreate tries to retrieve a challenge for the given timestamp and codelength from the ChallengeStorage.
// If it fails it will generate a new challenge
// It returns the resulting challenge
func (c Challenge) GetOrCreate(cs ChallengeStorage, timestamp int64, codeLen int) (Challenge, error) {
	cStr, err := cs.Challenge(timestamp, codeLen)

	if cStr == "" || err == cs.ErrNoChal() {
		c = c.Gen(codeLen)
		err = cs.Create(timestamp, codeLen, c.String())

		return c, err
	}

	return c.FromJson(cStr)
}

// Gen generates a new challenge
// It returns the generated challenge
func (c Challenge) Gen(n int) Challenge {
	randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))

	srcNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	srcCnt := len(srcNumbers)

	dstNrs := make([]int, n)
	for i := range dstNrs {
		dstNrs[i] = srcNumbers[randSrc.Intn(srcCnt)]
	}

	c.code = dstNrs

	return c
}

// Check will check whether the given guess contains any correct numbers
// It will return a slice of ints where:
// - 2 means 'correct number, correct position'
// - 1 means 'correct number, wrong position'
// - 0 means 'incorrect number'
func (c Challenge) Check(guess []int) ([]int, error) {
	size := len(c.code)

	// array for collecting return values
	ret := make([]int, size)

	// expect nr of items in guess is equal to nr of items in challenge
	if len(guess) != size {
		return ret, errors.New("length of guess and challenge is not equal")
	}

	// collect which guesses were in the right position
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

// WithCode will instantiate a new Challenge based on the given code
// It will return the resulting Challenge
func (c Challenge) WithCode(code []int) Challenge {
	c.code = code

	return c
}

// FromJson will instantiate a new Challenge from the given JSON encoded array of ints
// It will return the resulting Challenge and an error if it could not parse the JSON code
func (c Challenge) FromJson(cStr string) (Challenge, error) {
	err := json.Unmarshal([]byte(cStr), (&c.code))

	return c, err
}

// ToString will convert the code for this Challenge to a string
// It will return the resulting string
func (c Challenge) String() string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(c.code)), ","), "")
}
