package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Printer is what is going to print assertion violation error messages
type Printer interface {
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
		a.p.Errorf("jsonassert: cannot parse *payload* JSON: %s", err.Error())
		return // Really a t.Fatalf, but want a minimal Printer interface
	}
	expectedJSON, err := simplejson.NewJson([]byte(fmt.Sprintf(format, args...)))
	if err != nil {
		a.p.Errorf("jsonassert: cannot parse *expected JSON*: %s", err.Error())
		return // Really a t.Fatalf, but want a minimal Printer interface
	}
	assertion := &assertion{asserter: a, exp: expectedJSON}
	assertion.checkMap(jsonMap, "")
}

func (a *assertion) checkMap(jsonMap map[string]*json.RawMessage, path string) {
	for k, v := range jsonMap {
		newPath := path + "." + k
		if path == "" {
			newPath = k
		}
		a.checkField(v, newPath)
	}
}

func (a *assertion) checkField(j *json.RawMessage, path string) {
	// Should be safe to ignore error here: we have already parsed the payload
	// as JSON at thiis point, and checkField is only called when we check
	// values for keys
	bytes, _ := j.MarshalJSON()
	valueAtPath := strings.Split(path, ".")
	// TODO: Do not actually know the type at this stage, this assertion will crash in a horrible fire pretty soon
	exp := a.exp.GetPath(valueAtPath...).MustString()
	got := string(bytes)[1 : len(bytes)-1] //TODO: this isn't very nice. I want to escape the quotes surrounding the JSON string here.
	if got != exp {
		a.asserter.p.Errorf(`Expected key: "%s" to have value "%+v" but was "%+v"`, path, exp, got)
	}
}
