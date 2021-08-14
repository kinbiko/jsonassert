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
	ja.Assertf(
		`{ "name": "Jayne Cobb", "age": 36, "averageTestScore": "88%" }`,
		`{ "name": "Jayne Cobb", "age": 36, "averageTestScore": "%s" }`, "28%",
	)
	//output:
	//expected string at '$.averageTestScore' to be '28%' but was '88%'
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
