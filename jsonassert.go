package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
)

// Printer is what is going to print assertion violation error messages
// In particular, *testing.T adheres to this interface.
type Printer interface {
	// Errorf takes a format string that can be augmented with the fmt.Sprintf
	// arguments given in the vararg.
	Errorf(string, ...interface{})
}

// Asserter exposes methods for asserting that JSON payloads match the given
// string representation of the JSON payload
type Asserter interface {
	AssertString(string, string, ...interface{})
}

// New creates a new Asserter based on the given Printer.
// The Printer will be a *testing.T in 99% of your use cases.
func New(p Printer) Asserter {
	return &asserter{p: p}

}

type asserter struct{ p Printer }

type assertion struct {
	asserter *asserter
	exp      *simplejson.Json
}

func (a *asserter) AssertString(payload, format string, args ...interface{}) {
	var jsonMap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(payload), &jsonMap)
	if err != nil {
		a.p.Errorf("The given payload is not JSON: \"%s\",\nnested error is: %s", payload, err.Error())
		return // Really a t.Fatalf, but want a minimal Printer interface
	}
	formatted := fmt.Sprintf(format, args...)
	expectedJSON, err := simplejson.NewJson([]byte(formatted))
	if err != nil {
		a.p.Errorf("The expected payload is not JSON: \"%s\",\nnested error is: %s", formatted, err.Error())
		return // Really a t.Fatalf, but want a minimal Printer interface
	}
	assertion := &assertion{asserter: a, exp: expectedJSON}
	assertion.checkMap(jsonMap, "")
}

func (a *assertion) checkMap(jsonMap map[string]*json.RawMessage, path string) {
	m, err := a.exp.Map()
	if err != nil {
		a.asserter.p.Errorf("Attempted to turn expected payload into a map and failed because the actual payload has a nested json object at path '%s'", path)
	}
	checkedKeys := make(map[string]bool)
	// Check that everything in the actual payload exists in the expected payload
	for k, actualV := range jsonMap {
		checkedKeys[k] = true
		newPath := path + "." + k
		if path == "" {
			newPath = k
		}
		a.checkField(actualV, m[k], newPath)
	}

	// Check that everything in the expected payload exists in the actual payload
	for k, v := range m {
		newPath := path + "." + k
		if path == "" {
			newPath = k
		}
		if !checkedKeys[k] {
			a.checkField(jsonMap[k], v, newPath)
		}
	}
}

func (a *assertion) checkField(actualVal *json.RawMessage, expVal interface{}, path string) {
	if actualVal == nil {
		a.asserter.p.Errorf("Expected key \"%s\" to have value \"%+v\" but was not present in the payload", path, expVal)
		return
	}
	bytes, err := actualVal.MarshalJSON()
	if err != nil {
		a.asserter.p.Errorf("Unexpected error when marshalling JSON: %s", err)
		return
	}
	pathSegments := strings.Split(path, ".")
	// TODO: Do not actually know the type at this stage, this assertion will crash in a horrible fire pretty soon
	exp := a.exp.GetPath(pathSegments...).MustString()
	got := string(bytes)[1 : len(bytes)-1] //TODO: this isn't very nice. I want to escape the quotes surrounding the JSON string here.
	if got != exp {
		a.asserter.p.Errorf(`Expected key: "%s" to have value "%+v" but was "%+v"`, path, exp, got)
	}
}
