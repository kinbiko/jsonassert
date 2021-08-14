package jsonassert_test

import (
	"fmt"
	"testing"

	"github.com/kinbiko/jsonassert"
)

func TestAssertf(t *testing.T) {
	tt := []struct {
		name string
		act  string
		exp  string
		msgs []string
	}{
		{name: "different types", act: `"true"`, exp: `true`, msgs: []string{
			`actual JSON (string) and expected JSON (boolean) were of different types at '$'`,
		}},
		{name: "empty", act: "", exp: "", msgs: []string{}},
		{name: "empty v null", act: "", exp: "null", msgs: []string{`'actual' JSON is not valid JSON: unable to identify JSON type of ""`}},
		{name: "null v empty", act: "null", exp: "", msgs: []string{`'expected' JSON is not valid JSON: unable to identify JSON type of ""`}},
		{name: "null", act: "null", exp: "null", msgs: []string{}},
		{name: "true", act: `true`, exp: `true`, msgs: []string{}},
		{name: "false", act: `false`, exp: `false`, msgs: []string{}},
		{name: "true v false", act: `true`, exp: `false`, msgs: []string{`expected boolean at '$' to be false but was true`}},
		{name: "false v true", act: `false`, exp: `true`, msgs: []string{`expected boolean at '$' to be true but was false`}},
		{name: "identical floats", act: `12.45`, exp: `12.45`, msgs: []string{}},
		{name: "identical negative ints", act: `-1245`, exp: `-1245`, msgs: []string{}},
		{name: "different floats", act: `12.45`, exp: `1.245`, msgs: []string{`expected number at '$' to be '1.2450000' but was '12.4500000'`}},
		{name: "different ints", act: `1245`, exp: `-1245`, msgs: []string{`expected number at '$' to be '-1245.0000000' but was '1245.0000000'`}},
		{name: "identical strings", act: `"hello world"`, exp: `"hello world"`, msgs: []string{}},
		{name: "identical empty strings", act: `""`, exp: `""`, msgs: []string{}},
		{name: "different strings", act: `"hello"`, exp: `"world"`, msgs: []string{`expected string at '$' to be 'world' but was 'hello'`}},
		{name: "different strings", act: `"lorem ipsum dolor sit amet lorem ipsum dolor sit amet"`, exp: `"lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?"`, msgs: []string{
			`expected string at '$' to be
'lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?'
but was
'lorem ipsum dolor sit amet lorem ipsum dolor sit amet'`,
		}},
		{name: "empty v non-empty string", act: `""`, exp: `"world"`, msgs: []string{`expected string at '$' to be 'world' but was ''`}},
		{name: "identical objects", act: `{"hello": "world"}`, exp: `{"hello":"world"}`, msgs: []string{}},
		{name: "different keys in objects", act: `{"world": "hello"}`, exp: `{"hello":"world"}`, msgs: []string{
			`unexpected object key(s) ["world"] found at '$'`,
			`expected object key(s) ["hello"] missing at '$'`,
		}},
		{name: "different values in objects", act: `{"foo": "hello"}`, exp: `{"foo": "world" }`, msgs: []string{
			`expected string at '$.foo' to be 'world' but was 'hello'`,
		}},
		{name: "different keys in nested objects", act: `{"foo": {"world": "hello"}}`, exp: `{"foo": {"hello":"world"}   }`, msgs: []string{
			`unexpected object key(s) ["world"] found at '$.foo'`,
			`expected object key(s) ["hello"] missing at '$.foo'`,
		}},
		{name: "different values in nested objects", act: `{"foo": {"hello": "world"}}`, exp: `{"foo": {"hello":"世界"}   }`, msgs: []string{
			`expected string at '$.foo.hello' to be '世界' but was 'world'`,
		}},
		{name: "only one object is nested", act: `{}`, exp: `{"foo": {"hello":"世界"}   }`, msgs: []string{
			`expected 1 keys at '$' but got 0 keys`,
			`expected object key(s) ["foo"] missing at '$'`,
		}},
		{name: "empty array v empty array", act: `[]`, exp: `[ ]`, msgs: []string{}},
		{name: "non-empty array v empty array", act: `[null]`, exp: `[ ]`, msgs: []string{
			`length of arrays at '$' were different. Expected array to be of length 0, but contained 1 element(s)`,
			`actual JSON at '$' was: [null], but expected JSON was: []`,
		}},
		{name: "non-empty array v empty array", act: `[1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]`, exp: `[1,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]`, msgs: []string{
			`length of arrays at '$' were different. Expected array to be of length 22, but contained 30 element(s)`,
			`actual JSON at '$' was:
[1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]
but expected JSON was:
[1,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]`,
		}},
		{name: "identical non-empty arrays", act: `["hello"]`, exp: `["hello"]`, msgs: []string{}},
		{name: "different non-empty arrays", act: `["hello"]`, exp: `["world"]`, msgs: []string{
			`expected string at '$[0]' to be 'world' but was 'hello'`,
		}},
		{name: "identical non-empty unsorted arrays", act: `["hello", "world"]`, exp: `["<<UNORDERED>>", "world", "hello"]`, msgs: []string{}},
		{name: "different non-empty unsorted arrays", act: `["hello", "world"]`, exp: `["<<UNORDERED>>", "世界", "hello"]`, msgs: []string{
			`elements at '$' are different, even when ignoring order within the array:
expected some ordering of
["世界","hello"]
but got
["hello","world"]`,
		}},
		{name: "different length non-empty arrays", act: `["hello", "world"]`, exp: `["world"]`, msgs: []string{
			`length of arrays at '$' were different. Expected array to be of length 1, but contained 2 element(s)`,
			`actual JSON at '$' was: ["hello","world"], but expected JSON was: ["world"]`,
		}},
		{name: "presence against null", act: `{"foo": null}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{
			`expected the presence of any value at '$.foo', but was absent`,
		}},
		{name: "presence against boolean", act: `{"foo": true}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{}},
		{name: "presence against number", act: `{"foo": 1234}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{}},
		{name: "presence against string", act: `{"foo": "hello world"}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{}},
		{name: "presence against object", act: `{"foo": {"bar": "baz"}}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{}},
		{name: "presence against object", act: `{"foo": ["bar", "baz"]}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(st *testing.T) {
			tp, ja := setup()
			ja.Assertf(tc.act, tc.exp)
			if got := len(tp.messages); got != len(tc.msgs) {
				st.Errorf("expected %d assertion message(s) but got %d", len(tc.msgs), got)
				if len(tc.msgs) > 0 {
					st.Errorf("Expected the following messages:")
					for _, msg := range tc.msgs {
						st.Errorf(" - %s", msg)
					}
				}

				if len(tp.messages) > 0 {
					st.Errorf("Got the following messages:")
					for _, msg := range tp.messages {
						st.Errorf(" - %s", msg)
					}
				}
				return
			}
			for i := range tc.msgs {
				if exp, got := tc.msgs[i], tp.messages[i]; got != exp {
					st.Errorf("expected assertion message:\n'%s'\nbut got\n'%s'", exp, got)
				}
			}
		})
	}
}

func TestContainsf(t *testing.T) {
	tt := []struct {
		name string
		act  string
		exp  string
		msgs []string
	}{
		{name: "actual not valid json", act: `foo`, exp: `"foo"`, msgs: []string{
			`'actual' JSON is not valid JSON: unable to identify JSON type of "foo"`,
		}},
		{name: "expected not valid json", act: `"foo"`, exp: `foo`, msgs: []string{
			`'expected' JSON is not valid JSON: unable to identify JSON type of "foo"`,
		}},
		{name: "number contains a number", act: `5`, exp: `5`, msgs: nil},
		{name: "number does not contain a different number", act: `5`, exp: `-2`, msgs: []string{
			"expected number at '$' to be '-2.0000000' but was '5.0000000'",
		}},
		{name: "string contains a string", act: `"foo"`, exp: `"foo"`, msgs: nil},
		{name: "string does not contain a different string", act: `"foo"`, exp: `"bar"`, msgs: []string{
			"expected string at '$' to be 'bar' but was 'foo'",
		}},
		{name: "boolean contains a boolean", act: `true`, exp: `true`, msgs: nil},
		{name: "boolean does not contain a different boolean", act: `true`, exp: `false`, msgs: []string{
			"expected boolean at '$' to be false but was true",
		}},
		{name: "empty array contains empty array", act: `[]`, exp: `[]`, msgs: nil},
		{name: "single-element array contains empty array", act: `["fish"]`, exp: `[]`, msgs: nil},
		{name: "unordered empty array contains empty array", act: `[]`, exp: `["<<UNORDERED>>"]`, msgs: nil},
		{name: "unordered single-element array contains empty array", act: `["fish"]`, exp: `["<<UNORDERED>>"]`, msgs: nil},
		{name: "empty array contains single-element array", act: `[]`, exp: `["fish"]`, msgs: []string{
			"length of expected array at '$' was longer (length 1) than the actual array (length 0)",
			`actual JSON at '$' was: [], but expected JSON to contain: ["fish"]`,
		}},
		{name: "unordered multi-element array contains subset", act: `["alpha", "beta", "gamma"]`, exp: `["<<UNORDERED>>", "beta", "alpha"]`, msgs: nil},
		{name: "unordered multi-element array does not contain single element", act: `["alpha", "beta", "gamma"]`, exp: `["<<UNORDERED>>", "delta", "alpha"]`, msgs: []string{
			`element at $[1] in the expected payload was not found anywhere in the actual JSON array:
"delta"
not found in
["alpha","beta","gamma"]`,
		}},
		{name: "unordered multi-element array contains none of multi-element array", act: `["alpha", "beta", "gamma"]`, exp: `["<<UNORDERED>>", "delta", "pi", "omega"]`, msgs: []string{
			`element at $[1] in the expected payload was not found anywhere in the actual JSON array:
"delta"
not found in
["alpha","beta","gamma"]`,
			`element at $[2] in the expected payload was not found anywhere in the actual JSON array:
"pi"
not found in
["alpha","beta","gamma"]`,
			`element at $[3] in the expected payload was not found anywhere in the actual JSON array:
"omega"
not found in
["alpha","beta","gamma"]`,
		}},
		{name: "multi-element array contains itself", act: `["alpha", "beta"]`, exp: `["alpha", "beta"]`, msgs: nil},
		{name: "multi-element array does not contain itself permuted", act: `["alpha", "beta"]`, exp: `["beta" ,"alpha"]`, msgs: []string{
			"expected string at '$[0]' to be 'beta' but was 'alpha'",
			"expected string at '$[1]' to be 'alpha' but was 'beta'",
		}},
		// Allow users to test against a subset of the payload without erroring out.
		// This is to avoid the frustraion and unintuitive solution of adding "<<UNORDERED>>" in order to "enable" subsetting,
		// which is really implied with the `contains` part of the API name.
		{name: "multi-element array does contain its subset", act: `["alpha", "beta"]`, exp: `["alpha"]`, msgs: []string{}},
		{name: "multi-element array does not contain its superset", act: `["alpha", "beta"]`, exp: `["alpha", "beta", "gamma"]`, msgs: []string{
			"length of expected array at '$' was longer (length 3) than the actual array (length 2)",
			`actual JSON at '$' was: ["alpha","beta"], but expected JSON to contain: ["alpha","beta","gamma"]`,
		}},
		{name: "expected and actual have different types", act: `{"foo": "bar"}`, exp: `null`, msgs: []string{
			"actual JSON (object) and expected JSON (null) were of different types at '$'",
		}},
		{name: "expected any value, but got null", act: `{"foo": null}`, exp: `{"foo": "<<PRESENCE>>"}`, msgs: []string{
			"expected the presence of any value at '$.foo', but was absent",
		}},
		{name: "unordered multi-element array of different types contains subset", act: `["alpha", 5, false, ["foo"], {"bar": "baz"}]`, exp: `["<<UNORDERED>>", 5, "alpha", {"bar": "baz"}]`, msgs: nil},

		{name: "object contains its subset", act: `{"foo": "bar", "alpha": "omega"}`, exp: `{"alpha": "omega"}`, msgs: nil},
		/*
			{
				name: "big fat test",
				act: `{
						"arr": [
							"alpha",
							5,
							false,
							["foo"],
							{
								"bar": "baz",
								"fork": {
									"start": "stop"
								},
								"nested": ["really", "fast"]
							}
						],
						"fish": "mooney"
					}`,
				exp: `{
						"arr": [
							"<<UNORDERED>>",
							5,
							{
								"fork": {
									"start": "stop"
								},
								"nested": ["fast"]
							}
						]
					}`, msgs: nil},
		*/
	}
	for _, tc := range tt {
		t.Run(tc.name, func(st *testing.T) {
			tp, ja := setup()
			ja.Containsf(tc.act, tc.exp)
			if got := len(tp.messages); got != len(tc.msgs) {
				st.Errorf("expected %d assertion message(s) but got %d", len(tc.msgs), got)
				if len(tc.msgs) > 0 {
					st.Errorf("Expected the following messages:")
					for _, msg := range tc.msgs {
						st.Errorf(" - %s", msg)
					}
				}

				if len(tp.messages) > 0 {
					st.Errorf("Got the following messages:")
					for _, msg := range tp.messages {
						st.Errorf(" - %s", msg)
					}
				}
				return
			}

			if len(tc.msgs) == 1 {
				if exp, got := tc.msgs[0], tp.messages[0]; got != exp {
					st.Errorf("expected assertion message:\n'%s'\nbut got\n'%s'", exp, got)
				}
				return
			}

			for _, exp := range tc.msgs {
				found := false
				for _, got := range tp.messages {
					if got == exp {
						found = true
					}
				}
				if !found {
					st.Errorf("couldn't find expected assertion message:\n'%s'", exp)
				}
			}
		})
	}

}

func setup() (*testPrinter, *jsonassert.Asserter) {
	tp := &testPrinter{}
	return tp, jsonassert.New(tp)
}

type testPrinter struct {
	messages []string
}

func (tp *testPrinter) Errorf(msg string, args ...interface{}) {
	tp.messages = append(tp.messages, fmt.Sprintf(msg, args...))
}

func (tp *testPrinter) Helper() {
	// Do nothing in tests
}
