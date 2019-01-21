package internal

import "testing"

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

	if got, expLen := len(tp.messages), 2; got != expLen {
		t.Fatalf("expected %d assertion messages but got %d", expLen, got)
	}
	if got, expMsg := tp.messages[0], "length of arrays at '$' were different. Actual JSON had length 2, whereas expected JSON had length 3"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
	if got, expMsg := tp.messages[1], "actual JSON at '$' was: [hello world], but expected JSON was: [goodbye cruel world]"; got != expMsg {
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
