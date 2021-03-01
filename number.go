package jsonassert

import (
	"math"
	"strconv"
)

// This is *probably* good enough. Can change this to be even smaller if necessary
const minDiff = 0.000001

func (a *Asserter) checkNumber(path string, act, exp float64) {
	a.tt.Helper()
	if diff := math.Abs(act - exp); diff > minDiff {
		a.tt.Errorf("expected number at '%s' to be '%.7f' but was '%.7f'", path, exp, act)
	}
}

func extractNumber(n string) (float64, error) {
	return strconv.ParseFloat(n, 64)
}
