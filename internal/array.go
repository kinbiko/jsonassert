package internal

import "fmt"

func (a *asserter) checkArray(level string, act, exp []interface{}) {
	if len(act) != len(exp) {
		a.printer.Errorf("length of arrays at '%s' were different. Actual JSON had length %d, whereas expected JSON had length %d", level, len(act), len(exp))
		a.printer.Errorf("actual JSON at '%s' was: %+v, but expected JSON was: %+v", level, act, exp)
		return
	}
	for i := range act {
		a.Assert(level+fmt.Sprintf("[%d]", i), serialize(act[i]), serialize(exp[i]))
	}
}
