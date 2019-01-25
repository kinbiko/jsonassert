package internal

import (
	"testing"
)

func TestAssert(t *testing.T) {
	tt := []struct {
		name string
		act  string
		exp  string
		msgs []string
	}{
		{name: "different types", act: `"true"`, exp: `true`, msgs: []string{
			`actual JSON (string) and expected JSON (boolean) were of different types.`,
		}},
		{name: "empty", act: "", exp: "", msgs: []string{}},
		{name: "empty v null", act: "", exp: "null", msgs: []string{`could not find type for actual JSON: unable to identify JSON type of ""`}},
		{name: "null v empty", act: "null", exp: "", msgs: []string{`could not find type for expected JSON: unable to identify JSON type of ""`}},
		{name: "null", act: "null", exp: "null", msgs: []string{}},
		{name: "true", act: `true`, exp: `true`, msgs: []string{}},
		{name: "false", act: `false`, exp: `false`, msgs: []string{}},
		{name: "true v false", act: `true`, exp: `false`, msgs: []string{`expected value at '$' to be false but was true`}},
		{name: "false v true", act: `false`, exp: `true`, msgs: []string{`expected value at '$' to be true but was false`}},
		{name: "identical floats", act: `12.45`, exp: `12.45`, msgs: []string{}},
		{name: "identical negative ints", act: `-1245`, exp: `-1245`, msgs: []string{}},
		{name: "different floats", act: `12.45`, exp: `1.245`, msgs: []string{`expected value at '$' to be '1.2450000' but was '12.4500000'`}},
		{name: "different ints", act: `1245`, exp: `-1245`, msgs: []string{`expected value at '$' to be '-1245.0000000' but was '1245.0000000'`}},
		{name: "identical strings", act: `"hello world"`, exp: `"hello world"`, msgs: []string{}},
		{name: "identical empty strings", act: `""`, exp: `""`, msgs: []string{}},
		{name: "different strings", act: `"hello"`, exp: `"world"`, msgs: []string{`expected value at '$' to be 'world' but was 'hello'`}},
		{name: "empty v non-empty string", act: `""`, exp: `"world"`, msgs: []string{`expected value at '$' to be 'world' but was ''`}},
		{name: "identical objects", act: `{"hello": "world"}`, exp: `{"hello":"world"}`, msgs: []string{}},
		{name: "different keys in objects", act: `{"world": "hello"}`, exp: `{"hello":"world"}`, msgs: []string{
			`at level '$', key(s) [world] present in actual JSON but not in expected JSON`,
			`at level '$', key(s) [hello] present in expected JSON but not in actual JSON`,
		}},
		{name: "different values in objects", act: `{"foo": "hello"}`, exp: `{"foo": "world" }`, msgs: []string{
			`expected value at '$.foo' to be 'world' but was 'hello'`,
		}},
		{name: "different keys in nested objects", act: `{"foo": {"world": "hello"}}`, exp: `{"foo": {"hello":"world"}   }`, msgs: []string{
			`at level '$.foo', key(s) [world] present in actual JSON but not in expected JSON`,
			`at level '$.foo', key(s) [hello] present in expected JSON but not in actual JSON`,
		}},
		{name: "different values in nested objects", act: `{"foo": {"hello": "world"}}`, exp: `{"foo": {"hello":"世界"}   }`, msgs: []string{
			`expected value at '$.foo.hello' to be '世界' but was 'world'`,
		}},
		{name: "only one object is nested", act: `{}`, exp: `{"foo": {"hello":"世界"}   }`, msgs: []string{
			`different number of keys at level '$' in actual JSON (0) and expected JSON (1)`,
			`at level '$', key(s) [foo] present in expected JSON but not in actual JSON`,
		}},
		{name: "empty array v empty array", act: `[]`, exp: `[ ]`, msgs: []string{}},
		{name: "non-empty array v empty array", act: `[null]`, exp: `[ ]`, msgs: []string{
			`length of arrays at '$' were different. Actual JSON had length 1, whereas expected JSON had length 0`,
			`actual JSON at '$' was: [<nil>], but expected JSON was: []`,
		}},
		{name: "identical non-empty arrays", act: `["hello"]`, exp: `["hello"]`, msgs: []string{}},
		{name: "different non-empty arrays", act: `["hello"]`, exp: `["world"]`, msgs: []string{
			`expected value at '$[0]' to be 'world' but was 'hello'`,
		}},
		{name: "different length non-empty arrays", act: `["hello", "world"]`, exp: `["world"]`, msgs: []string{
			`length of arrays at '$' were different. Actual JSON had length 2, whereas expected JSON had length 1`,
			`actual JSON at '$' was: [hello world], but expected JSON was: [world]`,
		}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(st *testing.T) {
			tp, a := setup()
			a.Assert("$", tc.act, tc.exp)
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

func TestStringRepOf(t *testing.T) {
	tt := []struct {
		input interface{}
		exp   string
	}{
		{input: nil, exp: "null"},
		{input: true, exp: "true"},
		{input: false, exp: "false"},
		{input: 12.23, exp: "12.23"},
		{input: "", exp: `""`},
		{input: "abc", exp: `"abc"`},
		{input: map[string]interface{}{"hello": "world"}, exp: `{"hello":"world"}`},
		{input: map[string]interface{}{"hello": map[string]interface{}{"world": "世界"}}, exp: `{"hello":{"world":"世界"}}`},
		{
			input: []interface{}{"hello", 123, nil, map[string]interface{}{"hello": "world"}, []interface{}{"ok"}},
			exp:   `["hello",123,null,{"hello":"world"},["ok"]]`,
		},
	}
	for _, tc := range tt {
		if got := serialize(tc.input); got != tc.exp {
			t.Errorf("failed to get string rep of '%+v', expected\n'%s'\nbut got\n'%s'", tc.input, tc.exp, got)
		}
	}
}

func TestFindType(t *testing.T) {
	tt := []struct {
		input   string
		expType jsonType
	}{
		{input: `""`, expType: jsonString},
		{input: `123`, expType: jsonNumber},
		{input: `true`, expType: jsonBoolean},
		{input: `null`, expType: jsonNull},
		{input: `{}`, expType: jsonObject},
		{input: `[]`, expType: jsonArray},
	}

	for _, tc := range tt {
		t.Run(string(tc.expType), func(st *testing.T) {
			if got, err := findType(tc.input); got != tc.expType {
				if err != nil {
					st.Errorf("got error message when attempting to find type for '%s': '%s'", tc.input, err.Error())
				} else {
					st.Errorf("Expected input of '%s' to yield type '%s', but was '%s'", tc.input, tc.expType, got)
				}
			}
		})
	}
}
