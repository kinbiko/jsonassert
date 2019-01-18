# Algorithm design:

Level is `"$"` initially

## `internal/unknown.go`

### `Assert(level, act, exp string)`

1. Check the type of `exp` using `findType`
1. Check the type of `act` using `findType`

### findType(j string) (jsonType, error)

```go
if startsWith(j, "\"") {
    err := parseString(j)
    if err != nil { return nil, "friendlyErrorMessage" }
}
if startsWith(j, "%d+") {
    err := parseNumber(j)
    if err != nil { return nil, "friendlyErrorMessage" }
    return "number", nil
}
if j == "null" {
    err := parseNull(j)
    if err != nil { return nil, "friendlyErrorMessage" }
    return "null", nil
}
if startsWith(j,"{") {
    err := parseObject(j)
    if err != nil { return nil, "friendlyErrorMessage" }
    return "object", nil
}
if j == "true" || j = "false" {
    err := parseBoolean(j)
    if err != nil { return nil, "friendlyErrorMessage" }
    return "boolean", nil
}
if (startsWith(j, "[")) {
    err := parseArray(j)
    if err != nil { return nil, "friendlyErrorMessage" }
    return "array", nil
}
return nil, fmt.Errorf("unable to identify type of %s", j)
```

1. If an error is returned, then print error and return from this level.
1. If the types of `act` and `exp` are different, then print error and return from this level.

1. If the type is "null" for both, then return from this level (everything OK).

1. Switch on the type:
  1. String  -> checkString(level, act, exp)
  1. Number  -> checkNumber(level, extractNumber(act), extractNumber(exp))
  1. Boolean -> checkBoolean(level, extractBoolean(act), extractBoolean(exp))
  1. Object  -> checkObject(level, extractObject(act), extractObject(exp))
  1. Array   -> checkArray(level, extractArray(act), extractArray(exp))

### `serialize(a interface{}) string`

Essentially just a wrapper around `json.Marshal`. Should really never have an error at this level so let's panic.

```golang
bytes, err := json.Marshal(a)
```

## `internal/string.go`

### `checkString(level, act, exp string)`

1. Simple string comparison. Print error message if unequal.

## `internal/number.go`

### `checkNumber(level string, act, exp float64)`

1. Simple `float64` comparison. Print error message if unequal.

### `extractNumber(n string) float64`

Convert a `string` to a float64. This should be the same type that numbers are in JS.

- TODO: Precision? Raise a ticket to investigate if this becomes an issue during development.
- If this conversion errors: Panic, with an error message saying to raise an issue on the repo with the JSON that was passed in the API.
  - This should have been caught already in the type checking

## `internal/boolean.go`

### `extractBoolean(b string) bool`

Convert a `string` to a bool. Make sure to only accept `true` and `false`, and perhaps avoid a built-in converter that might not be as strict.

### `checkBoolean(level string, act, exp bool)`

1. Simple `bool` comparison. Print error message if unequal.

## `internal/object.go`

1. get number of keys in `exp`
1. get number of keys in `act`
1. If the numbers are different:
  1. print an error saying the number of keys is different
  1. gather all keys that exist in `exp` and not in `act`
      1. for all of these keys: report an error saying that this key was missing from the actual JSON
  1. gather all keys that exist in `act` and not in `exp`
      1. for all of these keys: report an error saying that this additional key was present
1. For each key that exists in both the expected JSON and the actual JSON:
  1. call `Assert("<level>.<key>", serialize(act, <key>), serialize())`

## `internal/array.go`

1. If the lengths of `act` and `get` are different:
    1. report an error saying they're of different length
    1. Use the default go representation of the array to print an assertion message of roughly how they're different
1. Else:
    1. call `Assert("<level>.<arrayIndex>", serialize(act, <arrayIndex>), serialize(exp, <arrayIndex>))`for each element index:
