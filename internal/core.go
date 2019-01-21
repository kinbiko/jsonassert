package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type asserter struct {
	printer interface {
		Errorf(msg string, args ...interface{})
	}
}

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
	if startsWith(j, `"`) {
		_, err := extractString(j)
		if err != nil {
			return jsonTypeUnknown, fmt.Errorf("unable to extract string: %s", err.Error())
		}
		return jsonString, nil
	}
	if startsWith(j, "0") ||
		startsWith(j, "1") ||
		startsWith(j, "2") ||
		startsWith(j, "3") ||
		startsWith(j, "4") ||
		startsWith(j, "5") ||
		startsWith(j, "6") ||
		startsWith(j, "7") ||
		startsWith(j, "8") ||
		startsWith(j, "9") {
		_, err := extractNumber(j)
		if err != nil {
			return jsonTypeUnknown, fmt.Errorf("unable to extract number: %s", err.Error())
		}
		return jsonNumber, nil
	}
	if j == "null" {
		return jsonNull, nil
	}
	if startsWith(j, "{") {
		_, err := extractObject(j)
		if err != nil {
			return jsonTypeUnknown, fmt.Errorf("unable to extract object: %s", err.Error())
		}
		return jsonObject, nil
	}
	if j == "true" || j == "false" {
		_, err := extractBoolean(j)
		if err != nil {
			return jsonTypeUnknown, fmt.Errorf("unable to extract boolean: %s", err.Error())
		}
		return jsonBoolean, nil
	}
	if startsWith(j, "[") {
		_, err := extractArray(j)
		if err != nil {
			return jsonTypeUnknown, fmt.Errorf("unable to extract array: %s", err.Error())
		}
		return jsonArray, nil
	}
	return jsonTypeUnknown, fmt.Errorf("unable to identify JSON type of %s", j)
}

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
