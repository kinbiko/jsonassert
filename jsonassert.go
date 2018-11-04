package jsonassert

// Printer is what is going to print assertion violation error messages
// In particular, *testing.T adheres to this interface.
type Printer interface {
	// Errorf takes a format string that can be augmented with the fmt.Sprintf
	// arguments given in the vararg.
	Errorf(string, ...interface{})
}

// New creates a new Asserter based on the given Printer.
// The Printer will be a *testing.T in 99% of your use cases.
func New(p Printer) Asserter {
	return &asserter{p}
}

// Asserter exposes methods for asserting that JSON payloads match the given
// string representation of the JSON payload
type Asserter interface {
	Assert(interface{}, string, ...interface{})
}
