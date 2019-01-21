package internal

import (
	"testing"
)

func TestEmptyObjectsAreEqual(t *testing.T) {
	tp, a := setup()
	act := map[string]interface{}{}
	exp := map[string]interface{}{}

	a.checkObject("$", act, exp)

	if got := len(tp.messages); got != 0 {
		t.Fatalf("expected 0 assertion messages written but got %d", got)
	}
}

func TestMoreKeysInActThanExp(t *testing.T) {
	tp, a := setup()
	act := map[string]interface{}{"hello": "world"}
	exp := map[string]interface{}{}

	a.checkObject("$", act, exp)

	if got := len(tp.messages); got != 2 {
		t.Fatalf("Expected 2 assertion messages but got %d messages", got)
	}
	if got, expMsg := tp.messages[0], "different number of keys at level '$' in actual JSON (1) and expected JSON (0)"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
	if got, expMsg := tp.messages[1], "at level '$', key(s) [hello] present in actual JSON but not in expected JSON"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
}

func TestMoreKeysExpInThanAct(t *testing.T) {
	tp, a := setup()
	act := map[string]interface{}{}
	exp := map[string]interface{}{"hello": "world"}

	a.checkObject("$", act, exp)

	if got := len(tp.messages); got != 2 {
		t.Fatalf("Expected 2 assertion messages but got %d messages", got)
	}
	if got, expMsg := tp.messages[0], "different number of keys at level '$' in actual JSON (0) and expected JSON (1)"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
	if got, expMsg := tp.messages[1], "at level '$', key(s) [hello] present in expected JSON but not in actual JSON"; got != expMsg {
		verifyAssertions(t, expMsg, got)
	}
}

func TestExtractObject(t *testing.T) {
	t.Run("empty object", func(st *testing.T) {
		s := "{}"
		got, err := extractObject(s)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if len(got) != 0 {
			t.Errorf("expected empty object but got '%+v'", got)
		}
	})
}
