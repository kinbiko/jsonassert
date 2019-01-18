package internal

import (
	"encoding/json"
	"errors"
)

// Assert checks that the given actual and expected strings are identical representations of JSON.
// If any discrepancies are found, these will be given to the Errorf function in the printer.
func (a *asserter) Assert(level string, act, exp string) {
}

func serialize(a interface{}) string {
	bytes, err := json.Marshal(a)
	if err != nil {
		// Really don't want to panic here, but I can't see a reasonable solution.
		// If this line *does* get executed then we should really investigate what kind of input was given
		panic(errors.New("unexpected failure to re-serialize nested JSON. Please raise an issue including this error message and both the expected and actual JSON strings you used to trigger this panic" + err.Error()))
	}
	return string(bytes)
}
