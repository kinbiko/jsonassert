package jsonassert

import (
	"encoding/json"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

type stringAssertion struct {
	asserter *asserter
	exp      *simplejson.Json
}

func (a *stringAssertion) checkMap(jsonMap map[string]*json.RawMessage, path string) {
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

func (a *stringAssertion) checkField(actualVal *json.RawMessage, expVal interface{}, path string) {
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
	// TODO: Do not actually know the type at this stage, this stringAssertion will crash in a horrible fire pretty soon
	exp := a.exp.GetPath(pathSegments...).MustString()
	got := string(bytes)[1 : len(bytes)-1] //TODO: this isn't very nice. I want to escape the quotes surrounding the JSON string here.
	if got != exp {
		if exp == "" {
			a.asserter.p.Errorf(`Unexpected key "%s" present in the payload`, path)
		} else {
			a.asserter.p.Errorf(`Expected key "%s" to have value "%+v" but was "%+v"`, path, exp, got)
		}
	}
}
