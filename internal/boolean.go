package internal

import "fmt"

func extractBoolean(b string) (bool, error) {
	if b == "true" {
		return true, nil
	}
	if b == "false" {
		return false, nil
	}
	return false, fmt.Errorf("could not parse '%s' as a boolean", b)
}

func (a *asserter) checkBoolean(level string, act, exp bool) {
	if act != exp {
		a.printer.Errorf("expected value at '%s' to be %v but was %v", level, exp, act)
	}
}
