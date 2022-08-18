package jsonassert

// noopHelperTT is used to wrap the Printer in the case that users pass in an
// Printer which does not implement a Helper() method. *testing.T does
// implement this method so it is believed that this utility will be largely
// unused.
type noopHelperTT struct {
	Printer
}

// Helper does nothing, intentionally. See New(Printer).
func (*noopHelperTT) Helper() {
	// Intentional NOOP
}

// deepEqualityPrinter is a utility Printer that lets the jsonassert package
// verify equality internally without affecting the external facing output,
// e.g. to verify where one potential candidate is matching or not.
type deepEqualityPrinter struct{ count int }

func (p *deepEqualityPrinter) Errorf(msg string, args ...interface{}) { p.count++ }
func (p *deepEqualityPrinter) Helper()                                { /* Intentional NOOP */ }

func (a *Asserter) deepEqual(act, exp interface{}) bool {
	p := &deepEqualityPrinter{count: 0}
	deepEqualityAsserter := &Asserter{tt: p}
	deepEqualityAsserter.pathassertf("", serialize(act), serialize(exp))
	return p.count == 0
}
