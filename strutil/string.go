package strutil

import (
	"math/rand"
	"time"
)

// SepEveryNth inserts strutil c after every nth character meant as a separator
// It returns the resulting strutil
func SepEveryNth(s string, n int, c string) string {
	if n < 1 {
		return s
	}

	for i := n; i < len(s); i += n + 1 {
		s = s[:i] + c + s[i:]
	}

	return s
}

// Rand generates a random strutil of n characters in length
// It returns the resulting strutil
func Rand(n int) string {
	rand.Seed(time.Now().UnixNano())

	srcChars := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	srcCnt := len(srcChars)

	dstChars := make([]rune, n)

	for i := range dstChars {
		dstChars[i] = srcChars[rand.Intn(srcCnt)]
	}

	return string(dstChars)
}
