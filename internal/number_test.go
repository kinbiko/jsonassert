package internal

import "testing"

func TestNumberComparison(t *testing.T) {
	t.Run("degenerate case", func(st *testing.T) {
		tp := &testPrinter{}
		a := &asserter{printer: tp}
		a.checkNumber("$", 0, 0)
		if got := len(tp.messages); got != 0 {
			st.Errorf("expect no printed messages but there were %d", got)
		}
	})

	t.Run("unequal case integer", func(st *testing.T) {
		tp := &testPrinter{}
		a := &asserter{printer: tp}
		a.checkNumber("$", 42, 1337)
		if len(tp.messages) != 1 {
			st.Errorf("expect exactly one printed message but there were %d", len(tp.messages))
		} else {
			exp, got := "expected value at '$' to be '1337.0000000' but was '42.0000000'", tp.messages[0]
			if exp != got {
				st.Errorf("Expected error message '%s' but got '%s'", exp, got)
			}
		}
	})

	t.Run("deeper level decimal", func(st *testing.T) {
		tp := &testPrinter{}
		a := &asserter{printer: tp}
		a.checkNumber("$.here.wat", 12.34, 43.21)
		if len(tp.messages) != 1 {
			st.Errorf("expect exactly one printed message but there were %d", len(tp.messages))
		} else {
			exp, got := "expected value at '$.here.wat' to be '43.2100000' but was '12.3400000'", tp.messages[0]
			if exp != got {
				st.Errorf("Expected error message '%s' but got '%s'", exp, got)
			}
		}
	})

	t.Run("unequal but within accepted difference", func(st *testing.T) {
		tp := &testPrinter{}
		a := &asserter{printer: tp}
		a.checkNumber("$.here.wat", 1.0000000, 1.0000003)
		if got := len(tp.messages); got != 0 {
			st.Errorf("expect no printed messages but there were %d", got)
			st.Errorf(tp.messages[0])
		}
	})
}

func TestExtractNumber(t *testing.T) {
	t.Run("case integer", func(st *testing.T) {
		got, err := extractNumber("24")
		if err != nil {
			t.Fatalf("Error was: '%s'", err)
		}
		exp := 24.0
		if got != exp {
			t.Errorf("expected %f but got %f", exp, got)
		}
	})

	t.Run("decimal", func(st *testing.T) {
		got, err := extractNumber("24.1340632")
		if err != nil {
			t.Fatalf("Error was: '%s'", err)
		}
		exp := 24.1340632
		if got != exp {
			t.Errorf("expected %f but got %f", exp, got)
		}
	})
}
