![logo](./logo.png)

`jsonassert` is a Go test assertion library for verifying that two representations of JSON are semantically equal.

## NEW ADDED

fork from https://github.com/kinbiko/jsonassert with new added functionality

### Check for key exists or not (Does not care its value is `null` or not)

Some properties of a JSON payload may be difficult to know in advance.
E.g. timestamps, UUIDs, or other randomly assigned values.

For these types of values, place the string `"<<KEY_EXISTS>>"` as the expected value, and `jsonassert` will only verify that this key exists (i.e. the actual JSON has the expected key, Does not care its value is `null` or not).

For example:

```go
func TestWhatever(t *testing.T) {
    ja := jsonassert.New(t)
    ja.Assertf(`
    {
        "time": "2019-01-28T21:19:42",
        "uuid": null
    }`, `
    {
        "time": "<<KEY_EXISTS>>",
        "uuid": "<<KEY_EXISTS>>"

    }`)
}
```


## Usage

Create a new `*jsonassert.Asserter` in your test and use this to make assertions against your JSON payloads:

```go
func TestWhatever(t *testing.T) {
    ja := jsonassert.New(t)
    // find some sort of payload
    ja.Assertf(payload, `
    {
        "name": "%s",
        "age": %d,
        "skills": [
            { "name": "martial arts", "level": 99 },
            { "name": "intelligence", "level": 100 },
            { "name": "mental fortitude", "level": 4 }
        ]
    }`, "River Tam", 16)
}
```

You may pass in `fmt.Sprintf` arguments after the expected JSON structure.

`ja.Assertf()` currently supports assertions against **strings only**.

### Check for presence only

Some properties of a JSON payload may be difficult to know in advance.
E.g. timestamps, UUIDs, or other randomly assigned values.

For these types of values, place the string `"<<PRESENCE>>"` as the expected value, and `jsonassert` will only verify that this key exists (i.e. the actual JSON has the expected key, and its value is not `null`), but this does not check its value.

For example:

```go
func TestWhatever(t *testing.T) {
    ja := jsonassert.New(t)
    ja.Assertf(`
    {
        "time": "2019-01-28T21:19:42",
        "uuid": "94ae1a31-63b2-4a55-a478-47764b60c56b"
    }`, `
    {
        "time": "<<PRESENCE>>",
        "uuid": "<<PRESENCE>>"

    }`)
}
```

The above will pass your test, but:

```go
func TestWhatever(t *testing.T) {
    ja := jsonassert.New(t)
    ja.Assertf(`
    {
        "date": "2019-01-28T21:19:42",
        "uuid": null
    }`, `
    {
        "time": "<<PRESENCE>>",
        "uuid": "<<PRESENCE>>"
    }`)
}
```

The above will fail your tests because the `time` key was not present in the actual JSON, and the `uuid` was `null`.

## Docs

You can find the [GoDocs for this package here](https://godoc.org/github.com/okisetiawan/jsonassert).

## Contributing

Contributions are welcome. Please discuss feature requests in an issue before opening a PR.
