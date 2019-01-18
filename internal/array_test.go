package internal

import (
	"testing"
)

func TestEmptyArraysAreEqual(t *testing.T) {
	tp, a := setup()
	act := []interface{}{}
	exp := []interface{}{}

	a.checkArray("$", act, exp)

	if got := len(tp.messages); got != 0 {
		t.Fatalf("expected 0 assertion messages written but got %d", got)
	}
}

func TestDifferentLengthArrays(t *testing.T) {
	tp, a := setup()
	act := []interface{}{"hello", "world"}
	exp := []interface{}{"goodbye", "cruel", "world"}

	a.checkArray("$", act, exp)

	if got, expLen := len(tp.messages), 3; got != expLen {
		t.Fatalf("expected %d assertion messages but got %d", expLen, got)
	}
	if got, expMsg := tp.messages[0], "length of arrays at '$' were different. Actual JSON had length 2, whereas expected JSON had length 3"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
	// TODO: instead expect unique elements more like how they look in JSON
	if got, expMsg := tp.messages[1], "element present in actual JSON but not in expected JSON: [hello]"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
	if got, expMsg := tp.messages[2], "element present in expected JSON but not in actual JSON: [goodbye cruel]"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
}

func TestEmptyArraysDifferentButSameLength(t *testing.T) {
	tp, a := setup()
	act := []interface{}{"The", "first", "letters"}
	exp := []interface{}{"The", "second", "word"}

	a.checkArray("$", act, exp)

	if got, expLen := len(tp.messages), 4; got != expLen {
		t.Fatalf("expected %d assertion messages written but got %d", expLen, got)
	}
	if got, expMsg := tp.messages[2], "expected element in position $[1] to be 'second' but was 'first'"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
	if got, expMsg := tp.messages[3], "expected element in position $[2] to be 'word' but was 'letters'"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
}

func TestSuperset(t *testing.T) {
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
		if got, expMsg := tp.messages[1], "element present in actual JSON but not in expected JSON: [word]"; got != expMsg {
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
		if got, expMsg := tp.messages[1], "element present in expected JSON but not in actual JSON: [word]"; got != expMsg {
			verifyAssertions(t, expMsg, got)
		}
	})
}

func verifyAssertions(t *testing.T, exp, got string) {
	t.Errorf("expected assertion message \n'%s' but got \n'%s'", exp, got)
}

func setup() (*testPrinter, *asserter) {
	tp := &testPrinter{}
	return tp, &asserter{printer: tp}
}
