/*
Package jsonassert is a Go test assertion library for verifying that two
representations of JSON are semantically equal. Create a new

	*jsonassert.Asserter

in your test and use this to make assertions against your JSON payloads:

	ja := jsonassert.New(t)

E.g. for the JSON

	{"hello": "world"}

you may use an expected JSON of

	{"hello": "%s"}

along with the "world" format argument. For example:

	ja.Assertf(`{"hello": "world"}`, `{"hello":"%s"}`, "world")

You may wish to make assertions against the *presence* of a value, but not
against its value. For example:

	ja.Assertf(`{"uuid": "94ae1a31-63b2-4a55-a478-47764b60c56b"}`, `{"uuid":"<<PRESENCE>>"}`)

will verify that the UUID field is present, but does not check its actual value.
You may use "<<PRESENCE>>" against any type of value. The only exception is null, which
will result in an assertion failure.

If you don't know / care about the order of the elements in an array in your
payload, you can ignore the ordering:

	payload := `["bar", "foo", "baz"]`
	ja.Assertf(payload, `["<<UNORDERED>>", "foo", "bar", "baz"]`)

The above will verify that "foo", "bar", and "baz" are exactly the elements in
the payload, but will ignore the order in which they appear.
*/
package jsonassert

import (
	"fmt"
)

// Printer is any type that has a testing.T-like Errorf function.
// You probably want to pass in a *testing.T instance here if you are using
// this in your tests.
type Printer interface {
	Errorf(msg string, args ...interface{})
}

// Asserter represents the main type within the jsonassert package.
// See Asserter.Assertf for the main use of this package.
type Asserter struct {
	tt
}

/*
New creates a new *jsonassert.Asserter for making assertions against JSON payloads.
This type can be reused. I.e. if you are using jsonassert as part of your tests,
you only need one *jsonassert.Asseter per (sub)test.
In most cases, this will look something like

	ja := jsonassert.New(t)

*/
func New(p Printer) *Asserter {
	// Initially this package was written without the assumption that the
	// provided Printer will implement testing.tt, which includes the Helper()
	// function to get better stacktraces in your testing utility functions.
	// This assumption was later added in order to get more accurate stackframe
	// information in test failures. In most cases users will pass in a
	// *testing.T to this function, which does adhere to that interface.
	// However, in order to be backwards compatible we also permit the use of
	// printers that do not implement Helper(). This is done by wrapping the
	// provided Printer into another struct that implements a NOOP Helper
	// method.
	if t, ok := p.(tt); ok {
		return &Asserter{tt: t}
	}
	return &Asserter{tt: &noopHelperTT{Printer: p}}
}

/*
Assertf takes two strings, the first being the 'actual' JSON that you wish to
make assertions against. The second string is the 'expected' JSON, which
can be treated as a template for additional format arguments.
If any discrepancies are found, these will be given to the Errorf function in the Printer.
E.g. for the JSON

	{"hello": "world"}

you may use an expected JSON of

	{"hello": "%s"}

along with the "world" format argument. For example:

	ja.Assertf(`{"hello": "world"}`, `{"hello":"%s"}`, "world")

You may also use format arguments in the case when your expected JSON contains
a percent character, which would otherwise be interpreted as a
format-directive.

	ja.Assertf(`{"averageTestScore": "99%"}`, `{"averageTestScore":"%s"}`, "99%")

You may wish to make assertions against the *presence* of a value, but not
against its value. For example:

	ja.Assertf(`{"uuid": "94ae1a31-63b2-4a55-a478-47764b60c56b"}`, `{"uuid":"<<PRESENCE>>"}`)

will verify that the UUID field is present, but does not check its actual value.
You may use "<<PRESENCE>>" against any type of value. The only exception is null, which
will result in an assertion failure.

If you don't know / care about the order of the elements in an array in your
payload, you can ignore the ordering:

	payload := `["bar", "foo", "baz"]`
	ja.Assertf(payload, `["<<UNORDERED>>", "foo", "bar", "baz"]`)

The above will verify that "foo", "bar", and "baz" are exactly the elements in
the payload, but will ignore the order in which they appear.
*/
func (a *Asserter) Assertf(actualJSON, expectedJSON string, fmtArgs ...interface{}) {
	a.tt.Helper()
	a.pathassertf("$", actualJSON, fmt.Sprintf(expectedJSON, fmtArgs...))
}
