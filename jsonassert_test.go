package jsonassert_test

import (
	"testing"

	"github.com/kinbiko/jsonassert"
)

func TestExternally(t *testing.T) {
	ja := jsonassert.New(t)
	ja.Assertf(`{"hello": "world"}`, `{"hello": "%s"}`, "world")
}
