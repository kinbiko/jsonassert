package jsonassert

import (
	"encoding/json"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
)

type asserter struct{ p Printer }

func (a *asserter) Assert(jsonPayload interface{}, assertionJSON string, args ...interface{}) {
	switch jsonPayload.(type) {
	case string:
		a.assertString(jsonPayload.(string), assertionJSON, args...)
	default:
		a.p.Errorf("Unsupported jsonPayload type: '%T'", jsonPayload)
	}
}

func (a *asserter) assertString(payload, format string, args ...interface{}) {
	jsonMap := make(map[string]*json.RawMessage)
	err := json.Unmarshal([]byte(payload), &jsonMap)
	if err != nil {
		a.p.Errorf("The given payload is not JSON: \"%s\",\nnested error is: %s", payload, err.Error())
		return
	}

	formatted := fmt.Sprintf(format, args...)
	expectedJSON, err := simplejson.NewJson([]byte(formatted))
	if err != nil {
		a.p.Errorf("The expected payload is not JSON: \"%s\",\nnested error is: %s", formatted, err.Error())
		return
	}
	assertion := &stringAssertion{asserter: a, exp: expectedJSON}
	assertion.checkMap(jsonMap, "")
}
