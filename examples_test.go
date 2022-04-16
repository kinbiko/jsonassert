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
// nolint:gochecknoglobals // this is global to make the examples look like valid test code
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
	expTestScore := "28%"
	ja.Assertf(
		`{ "name": "Jayne Cobb", "age": 36, "averageTestScore": "88%" }`,
		`{ "name": "Jayne Cobb", "age": 36, "averageTestScore": "%s" }`, expTestScore,
	)
	// output:
	// expected string at '$.averageTestScore' to be '28%' but was '88%'
}

func ExampleAsserter_Assertf_presenceOnly() {
	ja := jsonassert.New(t)
	ja.Assertf(`{"hi":"not the right key name"}`, `
		{
			"hello": "<<PRESENCE>>"
		}`)
	// output:
	// unexpected object key(s) ["hi"] found at '$'
	// expected object key(s) ["hello"] missing at '$'
}

func ExampleAsserter_Assertf_unorderedArray() {
	ja := jsonassert.New(t)
	ja.Assertf(
		`["zero", "one", "two"]`,
		`["<<UNORDERED>>", "one", "two", "three"]`,
	)
	// output:
	// actual JSON at '$[0]' contained an unexpected element: "zero"
	// expected JSON at '$[2]': "three" was missing from actual payload
}
