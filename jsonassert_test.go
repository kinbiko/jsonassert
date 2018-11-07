package jsonassert_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
)

type fakeT struct {
	receivedMessages []string
}

func (ft *fakeT) Errorf(format string, args ...interface{}) {
	ft.receivedMessages = append(ft.receivedMessages, fmt.Sprintf(format, args...))
}

type NestedStruct struct {
	Whatever bool `json:"whatever"`
}

type taggedStruct struct {
	MyNested *NestedStruct `json:"my_nested"`
	MyBool   bool          `json:"my_bool"`
	MyInt    int           `json:"my_int"`
	MyFloat  float64       `json:"my_float"`
	MyString string        `json:"my_string"`
}

var structExample = taggedStruct{
	MyNested: &NestedStruct{
		Whatever: false,
	},
	MyBool:   true,
	MyInt:    4123,
	MyFloat:  3.01,
	MyString: "foobar",
}

func TestAssert(t *testing.T) {
	jsonStr := []byte(`{"ok": "hello", "seven":7, "fish": true, "foo": {"bar": "baz"}}`)
	req, _ := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(jsonStr))
	tt := []struct {
		useCase       string
		payload       interface{}
		assertionJSON string
		args          []interface{}
		expAssertions []string
	}{
		{
			useCase:       "Simple valid check",
			payload:       `{"check": "ok"}`,
			assertionJSON: `{"check": "ok"}`,
			expAssertions: []string{},
		},

		{
			useCase:       "Unparseable payload",
			payload:       `Can't parse this`,
			assertionJSON: `{"check": "ok"}`,
			expAssertions: []string{`Invalid JSON given: "Can't parse this",
nested error is: invalid character 'C' looking for beginning of value`},
		},

		{
			useCase:       "Unparseable assertion JSON",
			payload:       `{"check": "ok"}`,
			assertionJSON: `Can't parse this`,
			expAssertions: []string{`Invalid JSON given: "Can't parse this",
nested error is: invalid character 'C' looking for beginning of value`},
		},

		{
			useCase:       "Mutiple violations, including string formatting",
			payload:       `{"check": "nope", "ok": "nah"}`,
			assertionJSON: `{"check": "%s", "ok": "yup"}`,
			args:          []interface{}{"works"},
			expAssertions: []string{
				`Expected key "check" to have value "works" but was "nope"`,
				`Expected key "ok" to have value "yup" but was "nah"`,
			},
		},

		{
			useCase:       "Payload < Assertion JSON",
			payload:       `{"ok": "yup"}`,
			assertionJSON: `{"check": "%s", "ok": "yup"}`,
			args:          []interface{}{"works"},
			expAssertions: []string{
				`Expected key "check" to have value "works" but was not present in the payload`,
			},
		},

		{
			useCase:       "Payload > Assertion JSON",
			payload:       `{"check": "works", "ok": "yup"}`,
			assertionJSON: `{"ok": "yup"}`,
			expAssertions: []string{
				`Unexpected key "check" present in the payload`,
			},
		},

		{
			useCase:       "Payload > Assertion JSON",
			payload:       `{"numbah": 3, "fish": "here"}`,
			assertionJSON: `{"numbah": 3}`,
			expAssertions: []string{`Unexpected key "fish" present in the payload`},
		},

		{
			useCase:       "Null in payload",
			payload:       `{"key": null}`,
			assertionJSON: `{"key": "hello"}`,
			expAssertions: []string{
				`Expected key "key" to have value "hello" but was not present in the payload`,
			},
		},

		{
			useCase:       "Null in assertion JSON",
			payload:       `{"key": "hello"}`,
			assertionJSON: `{"key": null}`,
			expAssertions: []string{
				`Unexpected key "key" present in the payload`,
			},
		},

		{
			useCase:       "Nested payload",
			payload:       `{"nested": {"check": "ok"}}`,
			assertionJSON: `{"nested": {"check": "%s"}}`,
			args:          []interface{}{"not ok"},
			expAssertions: []string{
				`Expected key "nested.check" to have value "not ok" but was "ok"`,
			},
		},

		{
			useCase: "<PRESENCE> keyword",
			payload: `{
				"uuid": "cb5230fc-f98f-4c63-abb7-d0588295983b",
				"timestamp": "2018-10-26T23:43:50+00:00"
			}`,
			assertionJSON: `{"uuid": "<PRESENCE>", "timestamp": "<PRESENCE>"}`,
		},

		{
			useCase:       "Differing types of value",
			payload:       `{"key": 539}`,
			assertionJSON: `{"key": "539"}`,
			expAssertions: []string{
				`Expected key "key" to have value "539" but was 539`,
			},
		},

		{
			useCase:       "Unsupported json payload type",
			payload:       []string{"wat"},
			assertionJSON: `{"key": "kagi"}`,
			expAssertions: []string{
				`Unsupported JSON type: '[]string'`,
			},
		},

		{
			useCase:       "Booleans",
			payload:       `{"key": true}`,
			assertionJSON: `{"key": false}`,
			expAssertions: []string{
				`Expected key "key" to have value 'false' but was 'true'`,
			},
		},

		{
			useCase:       "Arrays",
			payload:       `{"key": ["first", "second"]}`,
			assertionJSON: `{"key": []}`,
			expAssertions: []string{
				// TODO: Ideally this would be more JSON-like, but for now this'll do.
				`Expected key "key" to have value '[]' but was '[first second]'`,
			},
		},

		{
			useCase:       "Tagged struct",
			payload:       structExample,
			assertionJSON: `{"my_float": 3.01, "my_nested": {"whatever": false}, "my_string": "foobar", "my_int": 4123, "my_bool": true}`,
		},

		{
			useCase: "Non-Tagged struct",
			payload: struct {
				MyFloat  float64
				MyString string
				MyInt    int
				MyBool   bool
			}{
				MyFloat:  3.01,
				MyString: "foobar",
				MyInt:    4123,
				MyBool:   true,
			},
			assertionJSON: `{"MyFloat": 3.01, "MyString": "foobar", "MyInt": 4123, "MyBool": true}`,
		},

		{
			useCase:       "*http.Request",
			payload:       req,
			assertionJSON: `{"fish": true, "foo": {"bar": "baz"}, "seven": 7, "ok": "hello"}`,
		},
	}
	for _, tc := range tt {
		ft := new(fakeT)
		ja := jsonassert.New(ft)
		ja.Assert(tc.payload, tc.assertionJSON, tc.args...)

		msgs := ft.receivedMessages
		if exp, got := len(tc.expAssertions), len(msgs); exp != got {
			t.Errorf("'%s': Expected %d error messages to be written, but there were %d", tc.useCase, exp, got)
			if len(tc.expAssertions) > 0 {
				t.Errorf("'%s': Expected the following messages:", tc.useCase)
				for _, msg := range tc.expAssertions {
					t.Errorf(" - %s", msg)
				}
			}

			if len(msgs) > 0 {
				t.Errorf("'%s': Got the following messages:", tc.useCase)
				for _, msg := range msgs {
					t.Errorf(" - %s", msg)
				}
			}
			return //Don't attempt the following assertions

		}

		// The order of the JSON does not matter, so have to do a double subset check
		// Combines the issues in the end in order to make deciphering the test failure easier to parse
		unexpectedAssertions := ""
		for _, got := range msgs {
			found := false
			for _, exp := range tc.expAssertions {
				if got == exp {
					found = true
				}
			}
			if !found {
				if unexpectedAssertions == "" {
					unexpectedAssertions = "Got unexpected assertion failure:"
				}
				unexpectedAssertions += "\n - " + got
			}
		}

		missingAssertions := ""
		for _, got := range tc.expAssertions {
			found := false
			for _, exp := range msgs {
				if got == exp {
					found = true
				}
			}
			if !found {
				if missingAssertions == "" {
					missingAssertions = "\nExpected assertion failure but was not found:"
				}
				missingAssertions += "\n - " + got
			}
		}

		if totalError := unexpectedAssertions + missingAssertions; totalError != "" {
			t.Errorf("'%s': Inconsistent assertions:\n%s", tc.useCase, totalError)
		}
	}
}
