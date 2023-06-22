package jsonassert_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kinbiko/jsonassert"
)

func TestAssertf(t *testing.T) {
	t.Parallel()
	t.Run("primitives", func(t *testing.T) {
		t.Parallel()
		t.Run("equality", func(t *testing.T) {
			t.Parallel()
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
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})

		t.Run("difference", func(t *testing.T) {
			t.Parallel()
			for name, tc := range map[string]*testCase{
				"types":                    {`"true"`, `true`, []string{`actual JSON (string) and expected JSON (boolean) were of different types at '$'`}},
				"0 bytes v null":           {``, `null`, []string{`'actual' JSON is not valid JSON: unable to identify JSON type of ""`}},
				"booleans":                 {`false`, `true`, []string{`expected boolean at '$' to be true but was false`}},
				"floats":                   {`12.45`, `1.245`, []string{`expected number at '$' to be '1.2450000' but was '12.4500000'`}},
				"ints":                     {`1245`, `-1245`, []string{`expected number at '$' to be '-1245.0000000' but was '1245.0000000'`}},
				"strings":                  {`"hello"`, `"world"`, []string{`expected string at '$' to be 'world' but was 'hello'`}},
				"empty v non-empty string": {`""`, `"world"`, []string{`expected string at '$' to be 'world' but was ''`}},
			} {
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})
	})

	t.Run("objects", func(t *testing.T) {
		t.Parallel()
		t.Run("flat", func(t *testing.T) {
			t.Parallel()
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
					},
				},
			} {
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})

		t.Run("nested", func(t *testing.T) {
			t.Parallel()
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
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})

		t.Run("with PRESENCE directives", func(t *testing.T) {
			t.Parallel()
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
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})
	})

	t.Run("arrays", func(t *testing.T) {
		t.Parallel()
		t.Run("flat", func(t *testing.T) {
			t.Parallel()
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
				"identical non-empty unsorted arrays": {
					`["hello", "world"]`,
					`["<<UNORDERED>>", "world", "hello"]`,
					[]string{},
				},
			} {
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})

		t.Run("composite elements", func(t *testing.T) {
			t.Parallel()
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
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})

		t.Run("with UNORDERED directive", func(t *testing.T) {
			t.Parallel()
			for name, tc := range map[string]*testCase{
				"no elements":            {`[]`, `["<<UNORDERED>>"]`, nil},
				"only one equal element": {`["foo"]`, `["<<UNORDERED>>", "foo"]`, nil},
				"two elements ordered": {
					`["foo", "bar"]`,
					`["<<UNORDERED>>", "foo", "bar"]`,
					nil,
				},
				"two elements unordered": {
					`["bar", "foo"]`,
					`["<<UNORDERED>>", "foo", "bar"]`,
					nil,
				},
				"different number of elements": {
					`["foo"]`,
					`["<<UNORDERED>>", "foo", "bar"]`,
					[]string{
						`length of arrays at '$' were different. Expected array to be of length 2, but contained 1 element(s)`,
						`actual JSON at '$' was: ["foo"], but expected JSON was: ["foo","bar"], potentially in a different order`,
					},
				},
				"two different elements": {
					`["far", "boo"]`,
					`["<<UNORDERED>>", "foo", "bar"]`,
					[]string{
						`actual JSON at '$[0]' contained an unexpected element: "far"`,
						`actual JSON at '$[1]' contained an unexpected element: "boo"`,
						`expected JSON at '$[0]': "foo" was missing from actual payload`,
						`expected JSON at '$[1]': "bar" was missing from actual payload`,
					},
				},
				"valid array of different primitive types": {
					`["far", 1, null, true, [], {}]`,
					`["<<UNORDERED>>", true, 1, null, [], "far", {} ]`,
					nil,
				},
				"duplicates should still error out": {
					`["foo", "boo", "foo"]`,
					`["<<UNORDERED>>", "foo", "boo"]`,
					[]string{
						`length of arrays at '$' were different. Expected array to be of length 2, but contained 3 element(s)`,
						`actual JSON at '$' was: ["foo","boo","foo"], but expected JSON was: ["foo","boo"], potentially in a different order`,
					},
				},
				"nested unordered arrays": {
					// really long object means that serializing it the same is
					// highly unlikely should the determinism of JSON
					// serialization go away.
					`[{"20": 20}, {"19": 19}, {"18": 18 }, {"17": 17 }, {"16": 16 }, {"15": 15 }, {"14": 14 }, {"13": 13 }, {"12": 12 }, {"11": 11 }, {"10": 10 }, {"9": 9 }, {"8": 8 }, {"7": 7 }, {"6": 6 }, {"5": 5 }, {"4": 4 }, {"3": 3 }, {"2": 2 }, {"1": 1}]`,
					`["<<UNORDERED>>", {"1": 1}, {"2": 2}, {"3": 3}, {"4": 4}, {"5": 5}, {"6": 6}, {"7": 7}, {"8": 8}, {"9": 9}, {"10": 10}, {"11": 11}, {"12": 12}, {"13": 13}, {"14": 14}, {"15": 15}, {"16": 16}, {"17": 17}, {"18": 18}, {"19": 19}, {"20": 20}]`,
					nil,
				},
				"unordered array with objects with PRESENCE directives": {
					// Regression test for https://github.com/kinbiko/jsonassert/issues/39
					`{ "data": [ { "foo": 1, "bar": 2 }, { "foo": 11, "bar": 22 } ] }`,
					`{ "data": [ "<<UNORDERED>>", { "foo": 11, "bar": "<<PRESENCE>>" }, { "foo": 1, "bar": "<<PRESENCE>>" } ] }`,
					nil,
				},
			} {
				tc := tc
				t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
			}
		})
	})

	t.Run("extra long strings should be formatted on a new line", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]*testCase{
			"simple test string": {
				`"lorem ipsum dolor sit amet lorem ipsum dolor sit amet"`,
				`"lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?"`,
				[]string{
					`expected string at '$' to be
'lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?'
but was
'lorem ipsum dolor sit amet lorem ipsum dolor sit amet'`,
				},
			},
			"nested unordered arrays": {
				`["lorem ipsum dolor sit amet lorem ipsum dolor sit amet", "lorem ipsum dolor sit amet lorem ipsum dolor sit amet"]`,
				`["<<UNORDERED>>", "lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?"]`,
				[]string{
					`length of arrays at '$' were different. Expected array to be of length 1, but contained 2 element(s)`,
					`actual JSON at '$' was:
["lorem ipsum dolor sit amet lorem ipsum dolor sit amet","lorem ipsum dolor sit amet lorem ipsum dolor sit amet"]
but expected JSON was:
["lorem ipsum dolor sit amet lorem ipsum dolor sit amet why do I have to be the test string?"],
potentially in a different order`,
				},
			},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) { tc.checkAssertf(t) })
		}
	})

	t.Run("big fat test", func(t *testing.T) {
		t.Parallel()
		var (
			bigFatPayloadActual, _   = os.ReadFile("testdata/big-fat-payload-actual.json")
			bigFatPayloadExpected, _ = os.ReadFile("testdata/big-fat-payload-expected.json")
		)

		tc := testCase{
			act: fmt.Sprintf(`{
				"null": null,
				"emptyObject": {},
				"emptyArray": [],
				"emptyString": "",
				"zero": 0,
				"boolean": false,
				"positiveInt": 125,
				"negativeInt": -1245,
				"positiveFloats": 12.45,
				"negativeFloats": -12.345,
				"strings": "hello 世界",
				"flatArray": ["foo", "bar", "baz"],
				"flatObject": {"boo": "far", "biz": "qwerboipqwerb"},
				"nestedArray": ["boop", ["poob", {"bat": "boi", "asdf": 14, "oi": ["boy"]}], {"n": null}],
				"nestedObject": %s
			}`, string(bigFatPayloadActual)),
			exp: fmt.Sprintf(`{
				"nil": null,
				"emptyObject": [],
				"emptyArray": [null],
				"emptyString": " ",
				"zero": 0.00001,
				"boolean": true,
				"positiveInt": 124,
				"negativeInt": -1246,
				"positiveFloats": 11.45,
				"negativeFloats": -13.345,
				"strings": "hello world",
				"flatArray": ["fo", "ar", "baz"],
				"flatObject": {"bo": "far", "biz": "qwerboipqwer"},
				"nestedArray": ["oop", ["pob", {"bat": "oi", "asdf": 13, "oi": ["by"]}], {"m": null}],
				"nestedObject": %s
			}`, string(bigFatPayloadExpected)),
			msgs: []string{
				`unexpected object key(s) ["null"] found at '$'`,
				`expected object key(s) ["nil"] missing at '$'`,

				`actual JSON (object) and expected JSON (array) were of different types at '$.emptyObject'`,

				`length of arrays at '$.emptyArray' were different. Expected array to be of length 1, but contained 0 element(s)`,
				`actual JSON at '$.emptyArray' was: [], but expected JSON was: [null]`,

				`expected string at '$.emptyString' to be ' ' but was ''`,

				`expected number at '$.zero' to be '0.0000100' but was '0.0000000'`,

				`expected boolean at '$.boolean' to be true but was false`,

				`expected number at '$.positiveInt' to be '124.0000000' but was '125.0000000'`,

				`expected number at '$.negativeInt' to be '-1246.0000000' but was '-1245.0000000'`,

				`expected number at '$.positiveFloats' to be '11.4500000' but was '12.4500000'`,

				`expected number at '$.negativeFloats' to be '-13.3450000' but was '-12.3450000'`,

				`expected string at '$.strings' to be 'hello world' but was 'hello 世界'`,

				`expected string at '$.flatArray[0]' to be 'fo' but was 'foo'`,
				`expected string at '$.flatArray[1]' to be 'ar' but was 'bar'`,

				`unexpected object key(s) ["boo"] found at '$.flatObject'`,
				`expected object key(s) ["bo"] missing at '$.flatObject'`,
				`expected string at '$.flatObject.biz' to be 'qwerboipqwer' but was 'qwerboipqwerb'`,

				`expected string at '$.nestedArray[0]' to be 'oop' but was 'boop'`,
				`expected string at '$.nestedArray[1][0]' to be 'pob' but was 'poob'`,
				`expected number at '$.nestedArray[1][1].asdf' to be '13.0000000' but was '14.0000000'`,
				`expected string at '$.nestedArray[1][1].bat' to be 'oi' but was 'boi'`,
				`expected string at '$.nestedArray[1][1].oi[0]' to be 'by' but was 'boy'`,
				`unexpected object key(s) ["n"] found at '$.nestedArray[2]'`,
				`expected object key(s) ["m"] missing at '$.nestedArray[2]'`,

				`expected boolean at '$.nestedObject.is_full_report' to be false but was true`,
				`expected string at '$.nestedObject.id' to be 's869n10s9000060596qs3007' but was 's869n10s9000060s96qs3007'`,
				`actual JSON (object) and expected JSON (null) were of different types at '$.nestedObject.request.headers'`,
				`expected string at '$.nestedObject.metaData.device.userAgent' to be
'Mozilla/4.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36'
but was
'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36'`,
				`expected 7 keys at '$.nestedObject.source_map_failure' but got 8 keys`,
				`unexpected object key(s) ["source_map_url"] found at '$.nestedObject.source_map_failure'`,
				`expected boolean at '$.nestedObject.source_map_failure.has_uploaded_source_maps_for_version' to be true but was false`,
				`actual JSON at '$.nestedObject.breadcrumbs[1]' contained an unexpected element:
"Something that is most definitely missing from the expected one, right??"`,
				`expected JSON at '$.nestedObject.breadcrumbs[1]':
"Something that is most definitely missing from the actual one, right??"
was missing from actual payload`,
			},
		}
		tc.checkAssertf(t)
	})
}

func TestContainsf(t *testing.T) {
	t.Parallel()
	tt := map[string]*testCase{
		"actual not valid json": {
			`foo`,
			`"foo"`,
			[]string{`'actual' JSON is not valid JSON: unable to identify JSON type of "foo"`},
		},
		"expected not valid json":                             {`"foo"`, `foo`, []string{`'expected' JSON is not valid JSON: unable to identify JSON type of "foo"`}},
		"number contains a number":                            {`5`, `5`, nil},
		"number does not contain a different number":          {`5`, `-2`, []string{"expected number at '$' to be '-2.0000000' but was '5.0000000'"}},
		"string contains a string":                            {`"foo"`, `"foo"`, nil},
		"string does not contain a different string":          {`"foo"`, `"bar"`, []string{"expected string at '$' to be 'bar' but was 'foo'"}},
		"boolean contains a boolean":                          {`true`, `true`, nil},
		"boolean does not contain a different boolean":        {`true`, `false`, []string{"expected boolean at '$' to be false but was true"}},
		"empty array contains empty array":                    {`[]`, `[]`, nil},
		"single-element array contains empty array":           {`["fish"]`, `[]`, nil},
		"unordered empty array contains empty array":          {`[]`, `["<<UNORDERED>>"]`, nil},
		"unordered single-element array contains empty array": {`["fish"]`, `["<<UNORDERED>>"]`, nil},
		"empty array contains single-element array":           {`[]`, `["fish"]`, []string{"length of expected array at '$' was longer (length 1) than the actual array (length 0)", `actual JSON at '$' was: [], but expected JSON to contain: ["fish"]`}},
		"unordered multi-element array contains subset":       {`["alpha", "beta", "gamma"]`, `["<<UNORDERED>>", "beta", "alpha"]`, nil},
		"unordered multi-element array does not contain single element": {`["alpha", "beta", "gamma"]`, `["<<UNORDERED>>", "delta", "alpha"]`, []string{
			`element at $[1] in the expected payload was not found anywhere in the actual JSON array:
"delta"
not found in
["alpha","beta","gamma"]`,
		}},
		"unordered multi-element array contains none of multi-element array": {`["alpha", "beta", "gamma"]`, `["<<UNORDERED>>", "delta", "pi", "omega"]`, []string{
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
		"multi-element array contains itself": {`["alpha", "beta"]`, `["alpha", "beta"]`, nil},
		// NOTE: There's an important design decision to be made here.
		// Currently, in the case of "Containsf" there's an implicit "<<UNORDERED>>" (if it's explicitly written it will be ignored)
		// This is so that nested arrays don't have to repeatedly say "<<UNORDERED">> assuming the user just wants to check for the existence of some element of an array.
		// However, this makes jsonassert useless for cases where you want to partially assert that an ordered array exists.
		// Ideally this package should be able to support both nicely.
		"multi-element array does contain itself permuted": {`["alpha", "beta"]`, `["beta" ,"alpha"]`, []string{}},
		// Allow users to test against a subset of the payload without erroring out.
		// This is to avoid the frustraion and unintuitive solution of adding "<<UNORDERED>>" in order to "enable" subsetting,
		// which is really implied with the `contains` part of the API name.
		"multi-element array does contain its subset":                      {`["alpha", "beta"]`, `["beta"]`, []string{}},
		"multi-element array does not contain its superset":                {`["alpha", "beta"]`, `["alpha", "beta", "gamma"]`, []string{"length of expected array at '$' was longer (length 3) than the actual array (length 2)", `actual JSON at '$' was: ["alpha","beta"], but expected JSON to contain: ["alpha","beta","gamma"]`}},
		"expected and actual have different types":                         {`{"foo": "bar"}`, `null`, []string{"actual JSON (object) and expected JSON (null) were of different types at '$'"}},
		"expected any value but got null":                                  {`{"foo": null}`, `{"foo": "<<PRESENCE>>"}`, []string{"expected the presence of any value at '$.foo', but was absent"}},
		"unordered multi-element array of different types contains subset": {`["alpha", 5, false, ["foo"], {"bar": "baz"}]`, `["<<UNORDERED>>", 5, "alpha", {"bar": "baz"}]`, nil},
		"object contains its subset":                                       {`{"foo": "bar", "alpha": "omega"}`, `{"alpha": "omega"}`, nil},
		/*
			"array inside object": {
				`{ "arr": [                  { "fork": { "start": "stop" }, "nested": ["really",        "fast"] } ] }`,
				`{ "arr": [ "<<UNORDERED>>", { "fork": { "start": "stop" }, "nested": ["<<UNORDERED>>", "fast"] } ] }`,
				nil,
			},
		*/
	}
	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc.checkContainsf(t)
		})
	}
}

type testCase struct {
	act, exp string
	msgs     []string
}

func (tc *testCase) checkContainsf(t *testing.T) {
	t.Helper()
	tp := &testPrinter{messages: nil}
	jsonassert.New(tp).Containsf(tc.act, tc.exp)
	tc.check(t, tp)
}

func (tc *testCase) checkAssertf(t *testing.T) {
	t.Helper()
	tp := &testPrinter{messages: nil}
	jsonassert.New(tp).Assertf(tc.act, tc.exp)
	tc.check(t, tp)
}

func (tc *testCase) check(t *testing.T, tp *testPrinter) {
	t.Helper()
	if got := len(tp.messages); got != len(tc.msgs) {
		t.Errorf("expected %d assertion message(s) but got %d", len(tc.msgs), got)
	}

	for _, expMsg := range tc.msgs {
		found := false
		for _, printedMsg := range tp.messages {
			found = found || expMsg == printedMsg
		}
		if !found {
			t.Errorf("missing assertion message:\n%s", expMsg)
		}
	}

	for _, printedMsg := range tp.messages {
		found := false
		for _, expMsg := range tc.msgs {
			found = found || printedMsg == expMsg
		}
		if !found {
			t.Errorf("unexpected assertion message:\n%s", printedMsg)
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
