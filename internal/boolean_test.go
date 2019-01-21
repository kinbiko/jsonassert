package internal

import "testing"

func TestExtractBoolean(t *testing.T) {
	t.Run("Not a boolean", func(st *testing.T) {
		_, err := extractBoolean("fish")
		if err == nil {
			st.Fatalf("expected non-nil error message but was nil")
		}
		exp := "could not parse 'fish' as a boolean"
		if got := err.Error(); got != exp {
			st.Errorf("expected error message '%s' but was '%s'", exp, got)
		}
	})

	t.Run("true", func(st *testing.T) {
		got, err := extractBoolean("true")
		if err != nil {
			st.Fatalf("expected nil error message but was '%s'", err.Error())
		}
		if !got {
			st.Errorf("expected true returned false")
		}
	})

	t.Run("false", func(st *testing.T) {
		got, err := extractBoolean("false")
		if err != nil {
			st.Fatalf("expected nil error message but was '%s'", err.Error())
		}
		if got {
			st.Errorf("expected false returned true")
		}
	})
}

func TestCheckBoolean(t *testing.T) {
	t.Run("equal", func(st *testing.T) {
		tp, a := setup()
		a.checkBoolean("$", true, true)
		if got := len(tp.messages); got != 0 {
			st.Errorf("expect no printed messages but there were %d", got)
		}
	})

	t.Run("unequal", func(st *testing.T) {
		tp, a := setup()
		a.checkBoolean("$", true, false)
		if len(tp.messages) != 1 {
			st.Errorf("expect exactly one printed message but there were %d", len(tp.messages))
		} else {
			exp, got := "expected value at '$' to be false but was true", tp.messages[0]
			if exp != got {
				st.Errorf("Expected error message '%s' but got '%s'", exp, got)
			}
		}
	})
}
