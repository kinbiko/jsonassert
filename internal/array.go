package internal

func (a *asserter) checkArray(level string, act, exp []interface{}) {
	if len(act) != len(exp) {
		a.printer.Errorf("length of arrays at '%s' were different. Actual JSON had length %d, whereas expected JSON had length %d", level, len(act), len(exp))
	}
	if unique2Act := difference(act, exp); len(unique2Act) != 0 {
		a.printer.Errorf("element present in actual JSON but not in expected JSON: %v", unique2Act)
	}
	if unique2Exp := difference(exp, act); len(unique2Exp) != 0 {
		a.printer.Errorf("element present in expected JSON but not in actual JSON: %v", unique2Exp)
	}
	if len(act) == len(exp) {
		for i, actEl := range act {
			if expEl := exp[i]; actEl != expEl {
				a.printer.Errorf("expected element in position %s[%d] to be '%v' but was '%v'", level, i, expEl, actEl)
			}
		}
	}
}

func difference(a, b []interface{}) []interface{} {
	unique := []interface{}{}
	for _, e := range a {
		if !contains(b, e) {
			unique = append(unique, e)
		}
	}
	return unique
}

func contains(container []interface{}, containee interface{}) bool {
	for _, e := range container {
		if e == containee {
			return true
		}
	}
	return false
}
