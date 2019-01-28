package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Asserter struct {
	Printer interface {
		Errorf(msg string, args ...interface{})
	}
}

// Assertf checks that the given actual and expected strings are identical representations of JSON.
// If any discrepancies are found, these will be given to the Errorf function in the printer.
func (a *Asserter) Assertf(level string, act, exp string) {
	if act == exp {
		return
	}
	actType, err := findType(act)
	if err != nil {
		a.Printer.Errorf("could not find type for actual JSON: " + err.Error())
		return
	}
	expType, err := findType(exp)
	if err != nil {
		a.Printer.Errorf("could not find type for expected JSON: " + err.Error())
		return
	}
	if actType != expType {
		a.Printer.Errorf("actual JSON (%s) and expected JSON (%s) were of different types.", actType, expType)
		return
	}
	switch actType {
	case jsonBoolean:
		actBool, _ := extractBoolean(act)
		expBool, _ := extractBoolean(exp)
		a.checkBoolean(level, actBool, expBool)
	case jsonNumber:
		actNumber, _ := extractNumber(act)
		expNumber, _ := extractNumber(exp)
		a.checkNumber(level, actNumber, expNumber)
	case jsonString:
		actString, _ := extractString(act)
		expString, _ := extractString(exp)
		a.checkString(level, actString, expString)
	case jsonObject:
		actObject, _ := extractObject(act)
		expObject, _ := extractObject(exp)
		a.checkObject(level, actObject, expObject)
	case jsonArray:
		actArray, _ := extractArray(act)
		expArray, _ := extractArray(exp)
		a.checkArray(level, actArray, expArray)
	}
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
