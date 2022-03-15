package strutil

import (
	"testing"
)

type testSet struct {
	in, want, sep string
	n             int
}

func TestSepEveryNth(t *testing.T) {
	sets := []testSet{
		testSet{"aaaaaaaaaaa", "aaaa-aaaa-aaa", "-", 4},
		testSet{"aaaaaaaaaaaa", "aaaa-aaaa-aaaa", "-", 4},
		testSet{"aaaaaaaa", "aaaa-aaaa", "-", 4},
		testSet{"aaaa", "aaaa", "-", 4},
		testSet{"aaa", "aaa", "-", 3},
		testSet{"aaaaaaaaaaa", "a-a-a-a-a-a-a-a-a-a-a", "-", 1},
		testSet{"aaaaaaaaaaa", "aaa-aaa-aaa-aa", "-", 3},
		testSet{"-----", "------", "-", 4},
		testSet{"", "", "-", 1},
		testSet{"", "", "-", 0},
		testSet{"aaa", "aaa", "-", 0},
	}

	for _, row := range sets {
		var res = SepEveryNth(row.in, row.n, row.sep)

		if res != row.want {
			t.Fatalf(`SepEveryNth(%v, %v, %v): %v != %v`, row.in, row.n, row.sep, res, row.want)
		}
	}
}
