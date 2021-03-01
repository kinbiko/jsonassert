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
