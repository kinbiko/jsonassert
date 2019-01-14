package internal

import "math"

// This is *probably* good enough. Can change this to be even smaller if necessary
const minDiff = 0.000001

func (a *asserter) checkNumber(level string, act, exp float64) {
	diff := math.Abs(act - exp)
	if diff > minDiff {
		a.printer.Errorf("expected value at '%s' to be '%.7f' but was '%.7f'", level, exp, act)
	}
}
