package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
	switch actType {
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

func (a *Asserter) pathContainsf(path, act, exp string) {
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
	switch expType {
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
		a.checkContainsObject(path, actObject, expObject)
	case jsonArray:
		actArray, _ := extractArray(act)
		expArray, _ := extractArray(exp)
		a.checkContainsArray(path, actArray, expArray)
	case jsonNull:
		// Intentionally don't check as it wasn't expected in the payload
	}
}

func serialize(a interface{}) string {
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
	if _, err := extractString(j); err == nil {
		return jsonString, nil
	}
	if _, err := extractNumber(j); err == nil {
		return jsonNumber, nil
	}
	if j == "null" {
		return jsonNull, nil
	}
	if _, err := extractObject(j); err == nil {
		return jsonObject, nil
	}
	if _, err := extractBoolean(j); err == nil {
		return jsonBoolean, nil
	}
	if _, err := extractArray(j); err == nil {
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
