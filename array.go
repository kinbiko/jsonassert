package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Asserter) checkArray(level string, act, exp []interface{}) {
	if len(act) != len(exp) {
		a.Printer.Errorf("length of arrays at '%s' were different. Expected array to be of length %d, but contained %d element(s)", level, len(exp), len(act))
		a.Printer.Errorf("actual JSON at '%s' was: %+v, but expected JSON was: %+v", level, serialize(act), serialize(exp))
		return
	}
	for i := range act {
		a.pathassertf(level+fmt.Sprintf("[%d]", i), serialize(act[i]), serialize(exp[i]))
	}
}

func extractArray(s string) ([]interface{}, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, fmt.Errorf("cannot parse empty string as array")
	}
	if s[0] != '[' {
		return nil, fmt.Errorf("cannot parse '%s' as array", s)
	}
	var arr []interface{}
	err := json.Unmarshal([]byte(s), &arr)
	return arr, err
}
