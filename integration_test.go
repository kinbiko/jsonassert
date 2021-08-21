package jsonassert_test

import (
	"fmt"
	"testing"

	"github.com/kinbiko/jsonassert"
)

func TestAssertf(t *testing.T) {
	t.Run("primitives", func(t *testing.T) {
		t.Run("equality", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"0 bytes":         {``, ``, nil},
				"null":            {`null`, `null`, nil},
				"empty objects":   {`{}`, `{ }`, nil},
				"empty arrays":    {`[]`, `[ ]`, nil},
				"empty strings":   {`""`, `""`, nil},
				"zero":            {`0`, `0`, nil},
				"booleans":        {`false`, `false`, nil},
				"positive ints":   {`125`, `125`, nil},
				"negative ints":   {`-1245`, `-1245`, nil},
				"positive floats": {`12.45`, `12.45`, nil},
				"negative floats": {`-12.345`, `-12.345`, nil},
				"strings":         {`"hello world"`, `"hello world"`, nil},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})

		t.Run("difference", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"types":                    {`"true"`, `true`, []string{`actual JSON (string) and expected JSON (boolean) were of different types at '$'`}},
				"0 bytes v null":           {``, `null`, []string{`'actual' JSON is not valid JSON: unable to identify JSON type of ""`}},
				"booleans":                 {`false`, `true`, []string{`expected boolean at '$' to be true but was false`}},
				"floats":                   {`12.45`, `1.245`, []string{`expected number at '$' to be '1.2450000' but was '12.4500000'`}},
				"ints":                     {`1245`, `-1245`, []string{`expected number at '$' to be '-1245.0000000' but was '1245.0000000'`}},
				"strings":                  {`"hello"`, `"world"`, []string{`expected string at '$' to be 'world' but was 'hello'`}},
				"empty v non-empty string": {`""`, `"world"`, []string{`expected string at '$' to be 'world' but was ''`}},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})
	})

	t.Run("objects", func(t *testing.T) {
		t.Run("flat", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"identical objects": {
					`{"hello": "world"}`,
					`{"hello":"world"}`,
					nil,
				},
				"empty v non-empty object": {
					`{}`,
					`{"a": "b"}`,
					[]string{
						`expected 1 keys at '$' but got 0 keys`,
						`expected object key(s) ["a"] missing at '$'`,
					},
				},
				"different values in objects": {
					`{"foo": "hello"}`,
					`{"foo": "world" }`,
					[]string{`expected string at '$.foo' to be 'world' but was 'hello'`},
				},
				"different keys in objects": {
					`{"world": "hello"}`,
					`{"hello":"world"}`,
					[]string{
						`unexpected object key(s) ["world"] found at '$'`,
						`expected object key(s) ["hello"] missing at '$'`,
					}},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})

		t.Run("nested", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"different keys in nested objects": {
					`{"foo": {"world": "hello"}}`,
					`{"foo": {"hello": "world"}}`,
					[]string{
						`unexpected object key(s) ["world"] found at '$.foo'`,
						`expected object key(s) ["hello"] missing at '$.foo'`,
					},
				},
				"different values in nested objects": {
					`{"foo": {"hello": "world"}}`,
					`{"foo": {"hello":"世界"}}`,
					[]string{`expected string at '$.foo.hello' to be '世界' but was 'world'`},
				},
				"only one object is nested": {
					`{}`,
					`{ "foo": { "hello": "世界" } }`,
					[]string{
						`expected 1 keys at '$' but got 0 keys`,
						`expected object key(s) ["foo"] missing at '$'`,
					},
				},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})

		t.Run("with PRESENCE directives", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"presence against null": {
					`{"foo": null}`,
					`{"foo": "<<PRESENCE>>"}`,
					[]string{`expected the presence of any value at '$.foo', but was absent`},
				},
				"presence against boolean": {
					`{"foo": true}`,
					`{"foo": "<<PRESENCE>>"}`,
					nil,
				},
				"presence against number": {
					`{"foo": 1234}`,
					`{"foo": "<<PRESENCE>>"}`,
					nil,
				},
				"presence against string": {
					`{"foo": "hello world"}`,
					`{"foo": "<<PRESENCE>>"}`,
					nil,
				},
				"presence against object": {
					`{"foo": {"bar": "baz"}}`,
					`{"foo": "<<PRESENCE>>"}`,
					nil,
				},
				"presence against array": {
					`{"foo": ["bar", "baz"]}`,
					`{"foo": "<<PRESENCE>>"}`,
					nil,
				},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})
	})

	t.Run("arrays", func(t *testing.T) {
		t.Run("flat", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"empty array v empty array": {
					`[]`,
					`[ ]`,
					nil,
				},
				"non-empty array v empty array": {
					`[null]`,
					`[ ]`,
					[]string{
						`length of arrays at '$' were different. Expected array to be of length 0, but contained 1 element(s)`,
						`actual JSON at '$' was: [null], but expected JSON was: []`,
					},
				},
				"non-empty array v different non-empty array": {
					`[1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]`,
					`[1,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]`,
					[]string{
						`length of arrays at '$' were different. Expected array to be of length 22, but contained 30 element(s)`,
						`actual JSON at '$' was:
[1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]
but expected JSON was:
[1,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0]`,
					},
				},
				"identical non-empty arrays": {
					`["hello"]`,
					`["hello"]`,
					nil,
				},
				"different non-empty arrays": {
					`["hello"]`,
					`["world"]`,
					[]string{`expected string at '$[0]' to be 'world' but was 'hello'`},
				},
				"different length non-empty arrays": {
					`["hello", "world"]`,
					`["world"]`,
					[]string{
						`length of arrays at '$' were different. Expected array to be of length 1, but contained 2 element(s)`,
						`actual JSON at '$' was: ["hello","world"], but expected JSON was: ["world"]`,
					},
				},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})

		t.Run("composite elements", func(t *testing.T) {
			for name, tc := range map[string]*testCase{
				"single object with different values": {
					`[{"hello": "world"}]`,
					`[{"hello": "世界"}]`,
					[]string{`expected string at '$[0].hello' to be '世界' but was 'world'`},
				},
				"multiple nested object with different values": {
					`[
						{"hello": "world"},
						{"foo": {"bar": "baz"}}
					]`,
					`[
						{"hello": "世界"},
						{"foo": {"bat": "baz"}}
					]`,
					[]string{
						`expected string at '$[0].hello' to be '世界' but was 'world'`,
						`unexpected object key(s) ["bar"] found at '$[1].foo'`,
						`expected object key(s) ["bat"] missing at '$[1].foo'`,
					},
				},
				"array as array element": {
					`[["hello", "world"]]`,
					`[["hello", "世界"]]`,
					[]string{`expected string at '$[0][1]' to be '世界' but was 'world'`},
				},
				"multiple array elements": {
					`[["hello", "world"], [["foo"], "barz"]]`,
					`[["hello", "世界"], [["food"], "barz"]]`,
					[]string{
						`expected string at '$[0][1]' to be '世界' but was 'world'`,
						`expected string at '$[1][0][0]' to be 'food' but was 'foo'`,
					},
				},
			} {
				t.Run(name, func(t *testing.T) { tc.check(t) })
			}
		})
	})

	t.Run("extra long strings should be formatted on a new line", func(t *testing.T) {
		tc := &testCase{
			`"lorem ipsum dolor sit amet lorem ipsum dolor sit amet"`,
			`"lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?"`,
			[]string{`expected string at '$' to be
'lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?'
but was
'lorem ipsum dolor sit amet lorem ipsum dolor sit amet'`}}
		tc.check(t)
	})
}

type testCase struct {
	act, exp string
	msgs     []string
}

func (tc *testCase) check(t *testing.T) {
	tp := &testPrinter{}
	jsonassert.New(tp).Assertf(tc.act, tc.exp)
	if got := len(tp.messages); got != len(tc.msgs) {
		t.Errorf("expected %d assertion message(s) but got %d", len(tc.msgs), got)
		if len(tc.msgs) > 0 {
			t.Errorf("Expected the following messages:")
			for _, msg := range tc.msgs {
				t.Errorf(" - %s", msg)
			}
		}

		if len(tp.messages) > 0 {
			t.Errorf("Got the following messages:")
			for _, msg := range tp.messages {
				t.Errorf(" - %s", msg)
			}
		}
		return
	}
	for i := range tc.msgs {
		if exp, got := tc.msgs[i], tp.messages[i]; got != exp {
			t.Errorf("expected assertion message:\n'%s'\nbut got\n'%s'", exp, got)
		}
	}
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
