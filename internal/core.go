package internal

type asserter struct {
	printer interface {
		Errorf(msg string, args ...interface{})
	}
}
