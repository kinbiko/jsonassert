package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Asserter) checkString(level, act, exp string) {
	if act != exp {
		a.Printer.Errorf("expected string at '%s' to be '%s' but was '%s'", level, exp, act)
	}
}

func extractString(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return "", fmt.Errorf("cannot parse nothing as string")
	}
	if s[0] != '"' {
		return "", fmt.Errorf("cannot parse '%s' as string", s)
	}
	var str string
	err := json.Unmarshal([]byte(s), &str)
	return str, err
}
