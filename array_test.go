package jsonassert

import (
	"testing"
)

func TestCheckArray(t *testing.T) {
	t.Run("empty arrays are equal", func(st *testing.T) {
		tp, a := setup()
		act := []interface{}{}
		exp := []interface{}{}

		a.checkArray("$", act, exp)

		if got := len(tp.messages); got != 0 {
			st.Fatalf("expected 0 assertion messages written but got %d", got)
		}
	})

	t.Run("different length arrays", func(st *testing.T) {
		tp, a := setup()
		act := []interface{}{"hello", "world"}
		exp := []interface{}{"goodbye", "cruel", "world"}

		a.checkArray("$", act, exp)

		if got, expLen := len(tp.messages), 2; got != expLen {
			st.Fatalf("expected %d assertion messages but got %d", expLen, got)
		}
		if got, expMsg := tp.messages[0], "length of arrays at '$' were different. Actual JSON had length 2, whereas expected JSON had length 3"; got != expMsg {
			verifyAssertions(st, expMsg, got)
		}
		if got, expMsg := tp.messages[1], "actual JSON at '$' was: [hello world], but expected JSON was: [goodbye cruel world]"; got != expMsg {
			verifyAssertions(st, expMsg, got)
		}
	})

	t.Run("expected subset of actual", func(st *testing.T) {
		tp, a := setup()
		act := []interface{}{"The", "first", "word"}
		exp := []interface{}{"The", "first"}

		a.checkArray("$", act, exp)

		if got, expLen := len(tp.messages), 2; got != expLen {
			t.Fatalf("expected %d assertion messages written but got %d", expLen, got)
		}
		if got, expMsg := tp.messages[0], "length of arrays at '$' were different. Actual JSON had length 3, whereas expected JSON had length 2"; got != expMsg {
			verifyAssertions(t, expMsg, got)
		}
		if got, expMsg := tp.messages[1], "actual JSON at '$' was: [The first word], but expected JSON was: [The first]"; got != expMsg {
			verifyAssertions(t, expMsg, got)
		}
	})

	t.Run("actual subset of expected", func(st *testing.T) {
		tp, a := setup()
		act := []interface{}{"The", "first"}
		exp := []interface{}{"The", "first", "word"}

		a.checkArray("$", act, exp)

		if got, expLen := len(tp.messages), 2; got != expLen {
			t.Fatalf("expected %d assertion messages written but got %d", expLen, got)
		}
		if got, expMsg := tp.messages[0], "length of arrays at '$' were different. Actual JSON had length 2, whereas expected JSON had length 3"; got != expMsg {
			verifyAssertions(t, expMsg, got)
		}
		if got, expMsg := tp.messages[1], "actual JSON at '$' was: [The first], but expected JSON was: [The first word]"; got != expMsg {
			verifyAssertions(t, expMsg, got)
		}
	})
}

func TestExtractArray(t *testing.T) {
	t.Run("empty array", func(st *testing.T) {
		s := "[]"
		got, err := extractArray(s)
		if err != nil {
			st.Fatalf(err.Error())
		}
		if len(got) != 0 {
			st.Errorf(`Expected "%s" to have length 1 but had length %d`, s, len(got))
		}
	})
	t.Run("multiple elements in array", func(st *testing.T) {
		s := `[null, 1, "hello", true, ["world"], {"foo": "bar"}]`
		got, err := extractArray(s)
		if err != nil {
			st.Fatalf(err.Error())
		}
		if len(got) != 6 {
			st.Errorf(`Expected "%s" to have length 6 but had length %d`, s, len(got))
		}
		ind := 0
		if el, exp := got[ind], interface{}(nil); el != exp {
			st.Errorf(`Expected to find '%+v' in position %d but found %+v of type %T`, exp, ind, el, el)
		}
		ind = 1
		if el, exp := got[ind], float64(1); el != exp {
			st.Errorf(`Expected to find '%+v' in position %d but found %+v of type %T`, exp, ind, el, el)
		}
		ind = 2
		if el, exp := got[ind], "hello"; el != exp {
			st.Errorf(`Expected to find '%+v' in position %d but found %+v of type %T`, exp, ind, el, el)
		}
		ind = 3
		if el, exp := got[ind], true; el != exp {
			st.Errorf(`Expected to find '%+v' in position %d but found %+v of type %T`, exp, ind, el, el)
		}
		ind = 4
		if el, exp := got[ind].([]interface{}), "world"; el[0] != exp {
			st.Errorf(`Expected to find '%+v' in position %d but found %+v of type %T`, exp, ind, el, el)
		}
		ind = 5
		if el, exp := got[ind].(map[string]interface{}), "bar"; el["foo"] != exp {
			st.Errorf(`Expected to find '%+v' in position %d but found %+v of type %T`, exp, ind, el, el)
		}
	})
}
