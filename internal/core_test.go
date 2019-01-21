package internal

import "testing"

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
