package jsonassert

import (
	"encoding/json"
	"strings"
)

func (a *Asserter) checkString(path, act, exp string) {
	a.Helper()
	if act != exp {
		if len(exp+act) < maxMsgCharCount {
			a.Errorf("expected string at '%s' to be '%s' but was '%s'", path, exp, act)
		} else {
			a.Errorf("expected string at '%s' to be\n'%s'\nbut was\n'%s'", path, exp, act)
		}
	}
}

func extractString(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	if s[0] != '"' {
		return "", false
	}
	var str string
	return str, json.Unmarshal([]byte(s), &str) == nil
}
