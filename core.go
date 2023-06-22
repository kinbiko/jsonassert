package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

// The length at which to consider a message too long to fit on a single line
const maxMsgCharCount = 50

//nolint:gocyclo,cyclop // function is actually still readable
func (a *Asserter) pathassertf(path, act, exp string) {
	a.tt.Helper()
	if act == exp {
		return
	}
	actType, err := findType(act)
	if err != nil {
		a.tt.Errorf("'actual' JSON is not valid JSON: " + err.Error())
		return
	}
	expType, err := findType(exp)
	if err != nil {
		a.tt.Errorf("'expected' JSON is not valid JSON: " + err.Error())
		return
	}

	// If we're only caring about the presence of the key, then don't bother checking any further
	if expPresence, _ := extractString(exp); expPresence == "<<PRESENCE>>" {
		if actType == jsonNull {
			a.tt.Errorf(`expected the presence of any value at '%s', but was absent`, path)
		}
		return
	}

	if actType != expType {
		a.tt.Errorf("actual JSON (%s) and expected JSON (%s) were of different types at '%s'", actType, expType, path)
		return
	}
	switch actType { //nolint:exhaustive // already know it's valid JSON and not null
	case jsonBoolean:
		actBool, _ := extractBoolean(act)
		expBool, _ := extractBoolean(exp)
		a.checkBoolean(path, actBool, expBool)
	case jsonNumber:
		actNumber, _ := extractNumber(act)
		expNumber, _ := extractNumber(exp)
		a.checkNumber(path, actNumber, expNumber)
	case jsonString:
		actString, _ := extractString(act)
		expString, _ := extractString(exp)
		a.checkString(path, actString, expString)
	case jsonObject:
		actObject, _ := extractObject(act)
		expObject, _ := extractObject(exp)
		a.checkObject(path, actObject, expObject)
	case jsonArray:
		actArray, _ := extractArray(act)
		expArray, _ := extractArray(exp)
		a.checkArray(path, actArray, expArray)
	}
}

func serialize(a interface{}) string {
	//nolint:errchkjson // Can be confident this won't return an error: the
	// input will be a nested part of valid JSON, thus valid JSON
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

type jsonType string

const (
	jsonString      jsonType = "string"
	jsonNumber      jsonType = "number"
	jsonBoolean     jsonType = "boolean"
	jsonNull        jsonType = "null"
	jsonObject      jsonType = "object"
	jsonArray       jsonType = "array"
	jsonTypeUnknown jsonType = "unknown"
)

func findType(j string) (jsonType, error) {
	j = strings.TrimSpace(j)
	if _, ok := extractString(j); ok {
		return jsonString, nil
	}
	if _, ok := extractNumber(j); ok {
		return jsonNumber, nil
	}
	if j == "null" {
		return jsonNull, nil
	}
	if _, ok := extractObject(j); ok {
		return jsonObject, nil
	}
	if _, err := extractBoolean(j); err == nil {
		return jsonBoolean, nil
	}
	if _, ok := extractArray(j); ok {
		return jsonArray, nil
	}
	return jsonTypeUnknown, fmt.Errorf(`unable to identify JSON type of "%s"`, j)
}

// *testing.T has a Helper() func that allow testing tools like this package to
// ignore their own frames when calling Errorf on *testing.T instances.
// This interface is here to avoid breaking backwards compatibility in terms of
// the interface we expect in New.
type tt interface {
	Printer
	Helper()
}
