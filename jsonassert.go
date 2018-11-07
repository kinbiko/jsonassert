package jsonassert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Printer is what is going to print assertion violation error messages
// In particular, *testing.T adheres to this interface.
type Printer interface {
	// Errorf takes a format string that can be augmented with the fmt.Sprintf
	// arguments given in the vararg.
	Errorf(string, ...interface{})
}

var defaultConfig = Configuration{Verbosity: Default}

// New creates a new Asserter based on the given Printer.
// The Printer will be a *testing.T in 99% of your use cases.
func New(p Printer) Asserter {
	return &asserter{Printer: p, config: defaultConfig}
}

// Asserter exposes methods for asserting that JSON payloads match the given
// string representation of the JSON payload
type Asserter interface {
	Assert(interface{}, string, ...interface{})
}

type verbosity int

const (
	// Default verbosity will print each assertion error that occurs, including
	// the key and expected value if applicable.
	Default verbosity = iota
	// Verbose will print a pretty-printed version of the JSON against which
	// you're asserting in the case of an assertion failure. In the case that
	// the JSON cannot be pretty-printed, e.g. if it cannot be parsed as JSON
	// then the payload will be shown as a string.
	Verbose
)

// Configuration allows you to customise the Asserter further, e.g. to ignore
// the order of array elements or to make the assertion failure output more
// verobse.
type Configuration struct {
	Verbosity verbosity
}

type asserter struct {
	Printer
	config Configuration
}

type jsonType string

const (
	jsonString  jsonType = "string"
	jsonNumber  jsonType = "number"
	jsonObject  jsonType = "object"
	jsonBoolean jsonType = "boolean"
	jsonNull    jsonType = "null"
	jsonArray   jsonType = "array"
)

const presenceKeyword = `"<PRESENCE>"`

func (a *asserter) Assert(jsonPayload interface{}, assertionJSON string, args ...interface{}) {
	if reflect.ValueOf(jsonPayload).Kind() == reflect.Struct {
		b, err := json.Marshal(jsonPayload)
		if err != nil {
			a.fail(jsonPayload, "Unsupported JSON type: '%T'", jsonPayload)
		}
		a.checkMap(string(b), fmt.Sprintf(assertionJSON, args...), "")
	} else {
		switch jsonPayload.(type) {
		case string:
			a.checkMap(jsonPayload.(string), fmt.Sprintf(assertionJSON, args...), "")
		default:
			a.fail(jsonPayload, "Unsupported JSON type: '%T'", jsonPayload)
		}
	}
}

func (a *asserter) checkMap(payload, format, path string) {
	got, err := readStringAsJSON(payload)
	if err != nil {
		a.fail(payload, err.Error())
		return
	}

	exp, err := readStringAsJSON(format)
	if err != nil {
		a.fail(payload, err.Error())
		return
	}

	checkedKeys := make(map[string]bool)
	// Check that everything in the actual payload exists in the expected payload
	for k, actualV := range got {
		checkedKeys[k] = true
		newPath := path + "." + k
		if path == "" {
			newPath = k
		}
		a.checkMapField(actualV, exp[k], newPath)
	}

	// Check that everything in the expected payload exists in the actual payload
	for k, v := range exp {
		newPath := path + "." + k
		if path == "" {
			newPath = k
		}
		if !checkedKeys[k] {
			a.checkMapField(got[k], v, newPath)
		}
	}
}

func (a *asserter) checkMapField(got *json.RawMessage, exp *json.RawMessage, path string) {
	//If got is empty xor exp is empty (both should be impossible) then print a message saying so and return
	if got == nil {
		a.fail(got, `Expected key "%s" to have value %s but was not present in the payload`, path, *exp)
		return
	}
	gotBytes, _ := got.MarshalJSON()

	if exp == nil {
		a.fail(got, `Unexpected key "%s" present in the payload`, path)
		return
	}
	expBytes, _ := exp.MarshalJSON()
	// Then identify the type of both got and exp.
	gotType, expType := findType(gotBytes), findType(expBytes)
	// If the exp type is String and has value <PRESENCE>, then return without doing any further checking
	if expType == jsonString && string(expBytes) == presenceKeyword {
		return
	}

	// If they are the same, split into calling different methods for different types.
	// No need to check for null as we know both got an exp are the same type, and there's only one form of null
	switch gotType {
	case jsonString:
		a.checkString(string(gotBytes), string(expBytes), path)
	case jsonObject:
		a.checkMap(string(gotBytes), string(expBytes), path)
	case jsonArray:
		var g, e []interface{}
		json.Unmarshal(gotBytes, &g)
		json.Unmarshal(expBytes, &e)
		a.checkArray(g, e, path)
	case jsonBoolean:
		g, _ := strconv.ParseBool(string(gotBytes))
		e, _ := strconv.ParseBool(string(expBytes))
		a.checkBool(g, e, path)
	}
}

func (a *asserter) checkString(got, exp, path string) {
	if got != exp {
		if exp != "" {
			a.fail(got, `Expected key "%s" to have value %+v but was %+v`, path, exp, got)
		}
	}
}

func (a *asserter) checkBool(got, exp bool, path string) {
	if got != exp {
		a.fail(got, `Expected key "%s" to have value '%v' but was '%v'`, path, exp, got)
	}
}

func (a *asserter) checkArray(got, exp []interface{}, path string) {
	if len(got) != len(exp) {
		a.fail(got, `Expected key "%s" to have value '%v' but was '%v'`, path, exp, got)
	}
}

func findType(bytes []byte) jsonType {
	if string(bytes) == "true" || string(bytes) == "false" {
		return jsonBoolean
	}
	if bytes[0] == '[' {
		return jsonArray
	}
	if bytes[0] == '{' { //FIXME: Naive, but kidna works
		return jsonObject
	}
	return jsonString
}

func readStringAsJSON(s string) (map[string]*json.RawMessage, error) {
	j := make(map[string]*json.RawMessage)
	if err := json.Unmarshal([]byte(s), &j); err != nil {
		return nil, fmt.Errorf("Invalid JSON given: \"%s\",\nnested error is: %s", s, err.Error())
	}
	return j, nil
}

func (a *asserter) fail(actual interface{}, errorMessage string, args ...interface{}) {
	if a.config.Verbosity == Verbose {
		prettyPrint(actual)
	}
	a.Errorf(errorMessage, args...)
}

func prettyPrint(actual interface{}) {
	switch actual.(type) {
	case string:
		var out bytes.Buffer
		json.Indent(&out, []byte(actual.(string)), "", "    ")
		fmt.Println(string(out.Bytes()))
	default:
		fmt.Println(actual.(string))
	}
}
