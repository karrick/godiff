package godiff

import (
	"strings"
	"testing"
)

func ensureError(tb testing.TB, err error, contains ...string) {
	tb.Helper()
	if len(contains) == 0 || (len(contains) == 1 && contains[0] == "") {
		if err != nil {
			tb.Fatalf("GOT: %v; WANT: %v", err, contains)
		}
	} else if err == nil {
		tb.Errorf("GOT: %v; WANT: %v", err, contains)
	} else {
		for _, stub := range contains {
			if stub != "" && !strings.Contains(err.Error(), stub) {
				tb.Errorf("GOT: %v; WANT: %q", err, stub)
			}
		}
	}
}

func ensureStringSlicesMatch(tb testing.TB, actual, expected []string) {
	tb.Helper()

	al := len(actual)
	el := len(expected)
	var ai, ei int

	debug := false
	if debug {
		tb.Logf("SLICE: GOT:  %v", actual)
		tb.Logf("SLICE: WANT: %v", expected)
	}

	for ai < al || ei < el {
		if ai == al {
			if ei == el {
				break // only way out of loop when both ends reached
			}
			tb.Errorf("%10v | %10v", "", expected[ei])
			ei++
		} else if ei == el {
			tb.Errorf("%10v | %10v", actual[ai], "")
			ai++
		} else {
			if actual[ai] != expected[ei] {
				tb.Errorf("%10v | %10v", actual[ai], expected[ei])
			} else if debug {
				tb.Logf("%10v | %10v", actual[ai], expected[ei])
			}
			ai++
			ei++
		}
	}
}
