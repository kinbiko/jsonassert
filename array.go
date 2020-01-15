package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Asserter) checkArray(path string, act, exp []interface{}) {
	a.tt.Helper()

	var unordered bool
	if len(exp) > 0 && exp[0] == "<<UNORDERED>>" {
		unordered = true
		exp = exp[1:]
	}

	if len(act) != len(exp) {
		a.tt.Errorf("length of arrays at '%s' were different. Expected array to be of length %d, but contained %d element(s)", path, len(exp), len(act))
		serializedAct, serializedExp := serialize(act), serialize(exp)
		if len(serializedAct+serializedExp) < 50 {
			a.tt.Errorf("actual JSON at '%s' was: %+v, but expected JSON was: %+v", path, serializedAct, serializedExp)
		} else {
			a.tt.Errorf("actual JSON at '%s' was:\n%+v\nbut expected JSON was:\n%+v", path, serializedAct, serializedExp)
		}
		return
	}

	if unordered {
		for i := range act {
			hasMatch := false
			for j := range act {
				ap := arrayPrinter{}
				New(&ap).pathassertf("", serialize(act[i]), serialize(exp[j]))
				hasMatch = hasMatch || len(ap) == 0
			}
			if !hasMatch {
				serializedAct, serializedExp := serialize(act), serialize(exp)
				a.tt.Errorf("elements at '%s' are different, even when ignoring order within the array:\nexpected some ordering of\n%s\nbut got\n%s", path, serializedExp, serializedAct)
			}
		}
	} else {
		for i := range act {
			a.pathassertf(path+fmt.Sprintf("[%d]", i), serialize(act[i]), serialize(exp[i]))
		}
	}
}

func (a *Asserter) checkContainsArray(path string, act, exp []interface{}) {
	a.tt.Helper()

	var unordered bool
	if len(exp) > 0 && exp[0] == "<<UNORDERED>>" {
		unordered = true
		exp = exp[1:]
	}

	if len(act) < len(exp) {
		a.tt.Errorf("length of expected array at '%s' was longer (length %d) than the actual array (length %d).", path, len(exp), len(act))
		serializedAct, serializedExp := serialize(act), serialize(exp)
		if len(serializedAct+serializedExp) < 50 {
			a.tt.Errorf("actual JSON at '%s' was: %+v, but expected JSON was: %+v", path, serializedAct, serializedExp)
		} else {
			a.tt.Errorf("actual JSON at '%s' was:\n%+v\nbut expected JSON was:\n%+v", path, serializedAct, serializedExp)
		}
		return
	}

	if unordered {
		mismatchedExpPaths := map[string]string{}
		for i := range exp {
			for j := range act {
				ap := arrayPrinter{}
				serializedAct, serializedExp := serialize(act[i]), serialize(exp[j])
				New(&ap).pathContainsf("", serializedAct, serializedExp)
				if len(ap) > 0 {
					mismatchedExpPaths[fmt.Sprintf("%s[%d]", path, i+1)] = serializedExp // + 1 because 0th element is "<<UNORDERED>>"
				}
			}
		}
		for path, serializedExp := range mismatchedExpPaths {
			a.tt.Errorf("element at %s in the expected payload was not found anywhere in the actual JSON array:\n%s\nnot found in\n%s", path, serializedExp, serialize(act))
		}
	} else {
		for i := range exp {
			a.pathContainsf(fmt.Sprintf("%s[%d]", path, i), serialize(act[i]), serialize(exp[i]))
		}
	}
}

type arrayPrinter []string

func (p *arrayPrinter) Errorf(msg string, args ...interface{}) {
	n := append(*p, fmt.Sprintf(msg, args...))
	*p = n
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
