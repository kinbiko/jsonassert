package jsonassert

import (
	"encoding/json"
	"fmt"
)

type asserter struct{ p Printer }

type jsonType string

const (
	jsonString  jsonType = "string"
	jsonArray   jsonType = "array"
	jsonNumber  jsonType = "number"
	jsonNull    jsonType = "null"
	jsonObject  jsonType = "object"
	jsonBoolean jsonType = "boolean"
)

const presenceKeyword = "<PRESENCE>"

func (a *asserter) Assert(jsonPayload interface{}, assertionJSON string, args ...interface{}) {
	switch jsonPayload.(type) {
	case string:
		a.checkMap(jsonPayload.(string), fmt.Sprintf(assertionJSON, args...), "")
	default:
		a.p.Errorf("Unsupported jsonPayload type: '%T'", jsonPayload)
	}
}

func (a *asserter) checkMap(payload, format, path string) {
	got, err := readStringAsJSON(payload)
	if err != nil {
		a.p.Errorf(err.Error())
		return
	}

	exp, err := readStringAsJSON(format)
	if err != nil {
		a.p.Errorf(err.Error())
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
		a.p.Errorf(`Expected key "%s" to have value %s but was not present in the payload`, path, *exp)
		return
	}
	gotBytes, err := got.MarshalJSON()
	if err != nil {
		a.p.Errorf("Unexpected error when marshalling payload: %s", err)
		return
	}

	if exp == nil {
		a.p.Errorf(`Unexpected key "%s" present in the payload`, path)
		return
	}
	expBytes, err := exp.MarshalJSON()
	if err != nil {
		a.p.Errorf("Unexpected error when marshalling expected JSON: %s", err)
	}

	// Then identify the type of both got and exp.
	gotType, gotVal := a.findVal(gotBytes)
	expType, expVal := a.findVal(expBytes)
	// If the exp type is String and has value <PRESENCE>, then return without doing any further checking
	if expType == jsonString && expVal == presenceKeyword {
		return
	}

	// If their types are different, then write an error naming their types and their values.
	if expType != gotType {
		a.p.Errorf(`Types of key "%s" different in payload (%s) and expected payload (%s)`, gotType, expType)
		a.p.Errorf(`Got: '%+v'\nExp: '%+v'`, gotVal, expVal)
		a.p.Errorf(`Types of key "%s" different in payload (%s) and expected payload (%s)`, gotType, expType)
	}

	// If they are the same, split into calling different methods for different types.
	// No need to check for null as we know both got an exp are the same type, and there's only one form of null
	switch gotType {
	case jsonString:
		a.checkString(gotVal.(string), expVal.(string), path)
	case jsonObject:
		a.checkMap(string(gotBytes), string(expBytes), path)
	}
}

func (a *asserter) checkString(got, exp, path string) {
	if got != exp {
		if exp != "" {
			a.p.Errorf(`Expected key "%s" to have value %+v but was %+v`, path, exp, got)
		}
	}
}

func (a *asserter) findVal(bytes []byte) (jsonType, interface{}) {
	return jsonString, string(bytes)
}

func readStringAsJSON(s string) (map[string]*json.RawMessage, error) {
	j := make(map[string]*json.RawMessage)
	if err := json.Unmarshal([]byte(s), &j); err != nil {
		return nil, fmt.Errorf("Invalid JSON given: \"%s\",\nnested error is: %s", s, err.Error())
	}
	return j, nil
}
