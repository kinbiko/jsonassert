package jsonassert

import (
	"fmt"
)

type testPrinter struct {
	messages []string
}

func (tp *testPrinter) Errorf(msg string, args ...interface{}) {
	tp.messages = append(tp.messages, fmt.Sprintf(msg, args...))
}

func setup() (*testPrinter, *Asserter) {
	tp := &testPrinter{}
	return tp, &Asserter{Printer: tp}
}
