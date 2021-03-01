package jsonassert_test

import (
	"fmt"

	"github.com/kinbiko/jsonassert"
)

type printer struct{}

func (p *printer) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// using the varible name 't' to mimic a *testing.T variable
var t *printer

func ExampleNew() {
	ja := jsonassert.New(t)
	ja.Assertf(`{"hello":"world"}`, `
		{
			"hello": "world"
		}`)
}

func ExampleAsserter_Assertf_formatArguments() {
	ja := jsonassert.New(t)
	ja.Assertf(`{"hello":"世界"}`, `
		{
			"hello": "%s"
		}`, "world")
	//output:
	//expected string at '$.hello' to be 'world' but was '世界'
}

func ExampleAsserter_Assertf_presenceOnly() {
	ja := jsonassert.New(t)
	ja.Assertf(`{"hi":"not the right key name"}`, `
		{
			"hello": "<<PRESENCE>>"
		}`)
	//output:
	//unexpected object key(s) ["hi"] found at '$'
	//expected object key(s) ["hello"] missing at '$'
}
