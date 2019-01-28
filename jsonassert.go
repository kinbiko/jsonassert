package jsonassert

import (
	"fmt"
)

// Printer is any interface that has a testing.T-like Errorf function.
// Most users probably want to pass in a *testing.T instance here.
type Printer interface {
	Errorf(msg string, args ...interface{})
}

// Asserter represents the main type of jsonassert.
// See Asserter.Assert for the main use of this package.
type Asserter struct {
	Printer Printer
}

// New creates a new Asserter for making assertions against JSON.
// Can be reused. I.e. if you are using jsonassert as part of your tests,
// you only need one jsonassert.Asseter per test, which can be re-used in sub-tests.
// In most cases, this will look something like
// ja := jsonassert.New(t) // t is an instance of *testing.T
func New(p Printer) *Asserter {
	return &Asserter{Printer: p}
}

// Assertf takes two strings, the first being the 'actual' JSON that you wish to
// make assertions against. The second string is the 'expected' JSON, which
// can be treated as a template for additional format arguments.
// If any discrepancies are found, these will be given to the Errorf function in the printer.
// E.g. for the JSON {"hello": "world"}, you may use an expected JSON of
// {"hello": "%s"}, along with the "world" format argument.
// For example:
// ja.Assertf(`{"hello": "world"}`, `{"hello":"%s"}`, "world")
func (a *Asserter) Assertf(actualJSON, expectedJSON string, fmtArgs ...interface{}) {
	a.assertf("$", actualJSON, fmt.Sprintf(expectedJSON, fmtArgs...))
}
