package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Asserter) checkObject(path string, act, exp map[string]interface{}) {
	a.tt.Helper()
	if len(act) != len(exp) {
		a.tt.Errorf("expected %d keys at '%s' but got %d keys", len(exp), path, len(act))
	}
	if unique := difference(act, exp); len(unique) != 0 {
		a.tt.Errorf("unexpected object key(s) %+v found at '%s'", serialize(unique), path)
	}
	if unique := difference(exp, act); len(unique) != 0 {
		a.tt.Errorf("expected object key(s) %+v missing at '%s'", serialize(unique), path)
	}
	for key := range act {
		if contains(exp, key) {
			a.pathassertf(path+"."+key, serialize(act[key]), serialize(exp[key]))
		}
	}
}

func (a *Asserter) checkContainsObject(path string, act, exp map[string]interface{}) {
	a.tt.Helper()

	if missingExpected := difference(exp, act); len(missingExpected) != 0 {
		a.tt.Errorf("expected object key(s) %+v missing at '%s'", serialize(missingExpected), path)
	}
	for key := range exp {
		if contains(act, key) {
			a.pathContainsf(path+"."+key, serialize(act[key]), serialize(exp[key]))
		}
	}
}

// difference returns a slice of the keys that were found in a but not in b.
func difference(a, b map[string]interface{}) []string {
	unique := []string{}
	for key := range a {
		if !contains(b, key) {
			unique = append(unique, key)
		}
	}
	return unique
}

func contains(container map[string]interface{}, candidate string) bool {
	for key := range container {
		if key == candidate {
			return true
		}
	}
	return false
}

func extractObject(s string) (map[string]interface{}, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, fmt.Errorf("cannot parse empty string as object")
	}
	if s[0] != '{' {
		return nil, fmt.Errorf("cannot parse '%s' as object", s)
	}
	var arr map[string]interface{}
	err := json.Unmarshal([]byte(s), &arr)
	return arr, err
}
