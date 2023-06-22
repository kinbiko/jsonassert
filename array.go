package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Asserter) checkArray(path string, act, exp []interface{}) {
	a.tt.Helper()
	if len(exp) > 0 && exp[0] == "<<UNORDERED>>" {
		a.checkArrayUnordered(path, act, exp[1:])
	} else {
		a.checkArrayOrdered(path, act, exp)
	}
}

//nolint:gocognit,gocyclo,cyclop // function is actually still readable
func (a *Asserter) checkArrayUnordered(path string, act, exp []interface{}) {
	a.tt.Helper()
	if len(act) != len(exp) {
		a.tt.Errorf("length of arrays at '%s' were different. Expected array to be of length %d, but contained %d element(s)", path, len(exp), len(act))
		serializedAct, serializedExp := serialize(act), serialize(exp)
		if len(serializedAct+serializedExp) < maxMsgCharCount {
			a.tt.Errorf("actual JSON at '%s' was: %+v, but expected JSON was: %+v, potentially in a different order", path, serializedAct, serializedExp)
		} else {
			a.tt.Errorf("actual JSON at '%s' was:\n%+v\nbut expected JSON was:\n%+v,\npotentially in a different order", path, serializedAct, serializedExp)
		}
		return
	}

	for i, actEl := range act {
		found := false
		for _, expEl := range exp {
			if a.deepEqual(actEl, expEl) {
				found = true
			}
		}
		if !found {
			serializedEl := serialize(actEl)
			if len(serializedEl) < maxMsgCharCount {
				a.tt.Errorf("actual JSON at '%s[%d]' contained an unexpected element: %s", path, i, serializedEl)
			} else {
				a.tt.Errorf("actual JSON at '%s[%d]' contained an unexpected element:\n%s", path, i, serializedEl)
			}
		}
	}

	for i, expEl := range exp {
		found := false
		for _, actEl := range act {
			found = found || a.deepEqual(actEl, expEl)
		}
		if !found {
			serializedEl := serialize(expEl)
			if len(serializedEl) < maxMsgCharCount {
				a.tt.Errorf("expected JSON at '%s[%d]': %s was missing from actual payload", path, i, serializedEl)
			} else {
				a.tt.Errorf("expected JSON at '%s[%d]':\n%s\nwas missing from actual payload", path, i, serializedEl)
			}
		}
	}
}

func (a *Asserter) checkArrayOrdered(path string, act, exp []interface{}) {
	a.tt.Helper()
	if len(act) != len(exp) {
		a.tt.Errorf("length of arrays at '%s' were different. Expected array to be of length %d, but contained %d element(s)", path, len(exp), len(act))
		serializedAct, serializedExp := serialize(act), serialize(exp)
		if len(serializedAct+serializedExp) < maxMsgCharCount {
			a.tt.Errorf("actual JSON at '%s' was: %+v, but expected JSON was: %+v", path, serializedAct, serializedExp)
		} else {
			a.tt.Errorf("actual JSON at '%s' was:\n%+v\nbut expected JSON was:\n%+v", path, serializedAct, serializedExp)
		}
		return
	}
	for i := range act {
		a.pathassertf(path+fmt.Sprintf("[%d]", i), serialize(act[i]), serialize(exp[i]))
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
		a.tt.Errorf("length of expected array at '%s' was longer (length %d) than the actual array (length %d)", path, len(exp), len(act))
		serializedAct, serializedExp := serialize(act), serialize(exp)
		a.tt.Errorf("actual JSON at '%s' was: %+v, but expected JSON to contain: %+v", path, serializedAct, serializedExp)
		return
	}

	if unordered {
		a.checkContainsUnorderedArray(path, act, exp)
		return
	}
	for i := range exp {
		a.pathContainsf(fmt.Sprintf("%s[%d]", path, i), serialize(act[i]), serialize(exp[i]))
	}
}

func (a *Asserter) checkContainsUnorderedArray(path string, act, exp []interface{}) {
	mismatchedExpPaths := map[string]string{}
	for i := range exp {
		found := false
		serializedExp := serialize(exp[i])
		for j := range act {
			ap := arrayPrinter{}
			serializedAct := serialize(act[j])
			New(&ap).pathContainsf("", serializedAct, serializedExp)
			if len(ap) == 0 {
				found = true
			}
		}
		if !found {
			mismatchedExpPaths[fmt.Sprintf("%s[%d]", path, i+1)] = serializedExp // + 1 because 0th element is "<<UNORDERED>>"
		}
	}
	for path, serializedExp := range mismatchedExpPaths {
		a.tt.Errorf(`element at %s in the expected payload was not found anywhere in the actual JSON array:
%s
not found in
%s`,
			path, serializedExp, serialize(act))
	}
}

type arrayPrinter []string

func (p *arrayPrinter) Errorf(msg string, args ...interface{}) {
	n := append(*p, fmt.Sprintf(msg, args...))
	*p = n
}

func extractArray(s string) ([]interface{}, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, false
	}
	var arr []interface{}
	return arr, json.Unmarshal([]byte(s), &arr) == nil
}
