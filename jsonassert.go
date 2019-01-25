package jsonassert

import (
	"fmt"

	"github.com/kinbiko/jsonassert/internal"
)

type Printer interface {
	Errorf(msg string, args ...interface{})
}

type Asserter struct {
	asserter *internal.Asserter
}

func New(p Printer) *Asserter {
	return &Asserter{asserter: &internal.Asserter{Printer: p}}
}

func (a *Asserter) Assert(actualJSON, expectedJSON string, fmtArgs ...interface{}) {
	a.asserter.Assert("$", actualJSON, fmt.Sprintf(expectedJSON, fmtArgs...))
}
