package internal

import (
	"fmt"
	"testing"
)

type testPrinter struct {
	messages []string
}

func (tp *testPrinter) Errorf(msg string, args ...interface{}) {
	tp.messages = append(tp.messages, fmt.Sprintf(msg, args...))
}

func setup() (*testPrinter, *asserter) {
	tp := &testPrinter{}
	return tp, &asserter{printer: tp}
}

func verifyAssertions(t *testing.T, exp, got string) {
	t.Errorf("expected assertion message\n'%s' but got\n'%s'", exp, got)
}
