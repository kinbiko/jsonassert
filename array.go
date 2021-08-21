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

func (a *Asserter) checkArrayUnordered(path string, act, exp []interface{}) {
	a.tt.Helper()
	if len(act) != len(exp) {
		a.tt.Errorf("length of arrays at '%s' were different. Expected array to be of length %d, but contained %d element(s)", path, len(exp), len(act))
		serializedAct, serializedExp := serialize(act), serialize(exp)
		if len(serializedAct+serializedExp) < 50 {
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
			if len(serializedEl) < 50 {
				a.tt.Errorf("actual JSON at '%s[%d]' contained an unexpected element: %s", path, i, serializedEl)
			} else {
				a.tt.Errorf("actual JSON at '%s[%d]' contained an unexpected element:\n%s", path, i, serializedEl)
			}
		}
	}

	for i, expEl := range exp {
		found := false
		for _, actEl := range act {
			found = found || a.deepEqual(expEl, actEl)
		}
		if !found {
			serializedEl := serialize(expEl)
			if len(serializedEl) < 50 {
				a.tt.Errorf("expected JSON at '%s[%d]': %s was missing from actual payload", path, i, serializedEl)
			} else {
				a.tt.Errorf("expected JSON at '%s[%d]':\n%s\nwas missing from actual payload", path, i, serializedEl)
			}
		}
	}
}

func (a *Asserter) deepEqual(act, exp interface{}) bool {
	// There's a non-zero chance that JSON serialization will *not* be
	// deterministic in the future like it is in v1.16.
	// However, until this is the case, I can't seem to find a test case that
	// makes this evaluation return a false positive.
	// The benefit is a lot of simplicity and considerable performance benefits
	// for large nested structures.
	return serialize(act) == serialize(exp)
}

func (a *Asserter) checkArrayOrdered(path string, act, exp []interface{}) {
	a.tt.Helper()
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
	for i := range act {
		a.pathassertf(path+fmt.Sprintf("[%d]", i), serialize(act[i]), serialize(exp[i]))
	}
}

func extractArray(s string) ([]interface{}, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, fmt.Errorf("cannot parse empty string as array")
	}
	var arr []interface{}
	err := json.Unmarshal([]byte(s), &arr)
	return arr, err
}
