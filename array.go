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

func extractArray(s string) ([]interface{}, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, false
	}
	var arr []interface{}
	return arr, json.Unmarshal([]byte(s), &arr) == nil
}
