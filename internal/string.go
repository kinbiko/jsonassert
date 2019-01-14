package internal

func (a *asserter) checkString(level, act, exp string) {
	if act != exp {
		a.printer.Errorf("expected value at '%s' to be '%s' but was '%s'", level, exp, act)
	}
}
