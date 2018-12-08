# Design

## Nomenclature

- *actualJSON*: An object passed to jsonassert against which to make assertions.
- *expectedJSON*: A string passed in to jsonassert, possibly with format directives. This is the expected value of the JSON, generally written as a template string using backticks, in order for readability in tests.
- *formatParams*: The parameters to pass to parse the format directives in expectedJSON. Optional.
- *payload*: A custom type representing a JSON object. May represent several JSON types. Used internally to represent both the *expectedJSON* and the *actualJSON*.

## Version 1.0

Version 1 will not allow users to configure the library in any way. The JSON must match exactly in terms of content, but not necessarily be in order. Version 1 will only support strings as *actualJSON* types.

## Future (minor) versions

A later minor version bump could include configuration options. Desired configuration options include:

- Allow for certain properties to not be checked. For example:
    - when not present in the *expectedJSON*.
    - when explicitly declared as `"<PRESENCE>"` only the presence of the property should be checked, not the value
- Treat absent values as null
- Enable configuration of the output of the library. E.g:
    - Toggle colourised output
    - Toggle printing of successful assertions
    - Toggle printing of the entire *actualJSON*
    - Toggle printing of the entire *expectedJSON*, formatted
    - Toggle printing of a diff between the *actualJSON* and the *expectedJSON*

Support more than just Strings as the *actualJSON*:

- Strings containing JSON
- Structs with JSON tags
- *http.Request

Allow for truncated error messages. E.g. if the *expectedJSON* declares a highly nested object, whereas the actualJSON declares a short string, then only show the top-level keys of the object.
e.g.
```json
{
    "firstA": {
        "second:": {
            "fish": true
        }
    },
    "firstB": "my string value",
    "firstC": {
        "second:": {
            "fish": false
        }
    }
}
```

becomes

```json
{
    "firstA": "<object (truncated...)>",
    "firstB": "my string value",
    "firstC": "<object (truncated...)"
}
```

Set up a bug reporting shortcut: simple web app that lets you pass in two strings and a bunch of args which gets passed to Bugsnag, which integrates with GitHub issues, in order to automatically create issues.

## Design

### Goals

#### Error messages should be written in JSON and not Go

E.g. when describing an array of strings prefer to use a pretty JSON representation.

```json
[
  "myVal",
  "myOtherValue"
]
```

The use of `[myVal myOtherValue]` is discouraged.

#### Assertion messages should be unambiguous

All assertion messages should be associated with a JSONPath. This JSON path follows the [JSONPath specification](https://goessner.net/articles/JsonPath/).

#### Exported API should be minimal

The exported API, and nothing else, should all reside in one file: `jsonassert.go`, and be minimal in terms of its implementation, delegating wherever possible.

```go
// Asserter exposes the Assert function for checking that the passed in actualJSON
// matches the expectedJSON and the given formatParams.
type Asserter interface {
    Assert(actualJSON interface{}, expectedJSON string, formatParams... interface{})
}

// Failer is generally implemented by *testing.T
type Failer interface {
    Errorf(format string, params... interface{})
}

// New creates a new Asserter
func New(t Failer) Asserter
```

#### Pretty print JSON when showing error messages

[Use this](https://stackoverflow.com/questions/19038598/how-can-i-pretty-print-json-using-go) to make sure JSON is easily readable.

### Logic flow

1. Basic validation
    1. Assert that the *actualJSON* is in fact JSON
      1. error if not and return right away
    1. Assert that the *expectedJSON* is in fact JSON
      1. error if not and return right away

1. start level `$` [`parseObj`]
    1. Get type of *expectedJSON*
    1. Get type of *actualJSON*
    1. If they are different types
        1. Report error saying that the types are different
        1. return
    1. Else if they are the same:
        1. Switch on the type
            1. if strings are different report an error
            1. if numbers are different report an error
            1. if booleans are different report an error
            1. if object:
                1. start level `$.<objectKey>`
                1. get number of keys in expectedObj
                1. get number of keys in actualObj
                1. If the numbers are different:
                    1. report an error saying the number of keys is different
                    1. gather all keys that exist in expectedObj and not in actualObj
                        1. for all of these keys: report an error saying that this key was missing from the actual JSON, along with the pretty-printed JSON of the expectedObj
                    1. gather all keys that exist in actualObj and not in expectedObj
                        1. for all of these keys: report an error saying that this additional key was present, along with the pretty-printed JSON of actualObj
                1. If the keys are identical then we dive deeper into the rabbit hole:
                    1. for each key:
                        1. start level `<level>.<key>`
                        1. call [`parseObj`] with `expectedObj[<key>]` and `actualObj[<key>]`.
            1. if array:
                1. start level `$.<arrayKey>`
                1. get length of expectedArray
                1. get length of actualArray
                1. If the lengths are different:
                    1. report an error saying they're of different length
                    1. gather all elements that exist in expectedArray and not in actualArray
                        1. for all of these elements: report an error saying that this element was missing from the actual JSON, along with the pretty-printed JSON
                    1. gather all elements that exist in actualArray and not in expectedArray
                        1. for all of these elements: report an error saying that this additional element was present, along with the pretty-printed JSON
                1. If the lengths are identical then we also validate the order of the elements:
                    1. for each element index:
                        1. start level `<level>.<arrayKey>[<index>]`
                        1. call [`parseObj`] with `expectedArray[<index>]` and `actualArray[<index>]`.

#### Pseudocode:

```go
func Assert(actual, expected interface{}) {
    if !validJSON(actual) {
        p.Errorf("actual JSON could not be parsed as JSON")
        return
    }
    if !validJSON(expected) {
        p.Errorf("expected JSON could not be parsed as JSON")
        return
    }

    actualType := getType(actual);
    expectedType := getType(expected);

    if actualType != expectedType {
      p.Errorf("actual JSON (%s) and expected JSON (%s) are of different types and cannot be compared", actualType, expectedType)
      return
    }

    findTypeAsserter(actualType).assertEqual("$", actual, expected)
}

type typeAsserter interface {
    assertEqual(key string, actual, expected interface{})
}

func findTypeAsserter(jsonType string) {
    switch (jsonType) {
    case "string": return stringAsserter
    case "number": return numberAsserter
    case "boolean": return booleanAsserter
    case "null": return nullAsserter
    case "object": return objectAsserter
    case "array": return arrayAsserter
    default: panic("jsonassert fucked up")
    }
}

func stringAsserter(key string, a, e interface{}) {
    var actual, expected string
    if actual, err := a.(string); err == nil {
        panic(fmt.Sprintf("bug in jsonassert, incorrect asserter type called. Expected 'a' to be a string but was %T", a))
    }
    if expected, err := e.(string); err == nil {
        panic(fmt.Sprintf("bug in jsonassert, incorrect asserter type called. Expected 'e' to be a string but was %T", e))
    }
    if actual != expected {
        t.Errorf("Expected string to be 'a' but was 'b'")
    }
}

func arrayAsserter(key string, actual, expected interface{}) {
    aLen, eLen := len(actual), len(expected); aLen != eLen {
    }
                1. get length of expectedArray
                1. get length of actualArray
                1. If the lengths are different:
                    1. report an error saying they're of different length
                    1. gather all elements that exist in expectedArray and not in actualArray
                        1. for all of these elements: report an error saying that this element was missing from the actual JSON, along with the pretty-printed JSON
                    1. gather all elements that exist in actualArray and not in expectedArray
                        1. for all of these elements: report an error saying that this additional element was present, along with the pretty-printed JSON
                1. If the lengths are identical then we also validate the order of the elements:
                    1. for each element index:
                        1. start level `<level>.<arrayKey>[<index>]`
                        1. call [`parseObj`] with `expectedArray[<index>]` and `actualArray[<index>]`.
}
```

Note: Always base type checks on the expectedJSON, not the actualJSON.

#### Key components

- Validator: Validates that the user input is in fact valid JSON.
- FooBar: Struct that holds: (maybe it's a stack?)
    - Path
    - The string representation of the current actual JSON
    - The string representation of the current expected JSON
    - The type of the actual JSON, if known
    - The type of the expected JSON, if known

#### Key features

- level: our representation of a JSONPath. Likely a wrapper around [this library](https://github.com/yalp/jsonpath). These types should be created in a context-like fashion, in that you can create a new from an old, but there's no global state.
- jsonValue: our representation of any form of legal JSON value.
    - type identification: given a jsonValue, what's the jsonType
    - retrieval of a jsonValue's real value in methods per type.
- error logging
