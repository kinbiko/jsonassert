# `jsonassert`

[![Build Status](https://travis-ci.com/kinbiko/jsonassert.svg?branch=master)](https://travis-ci.com/kinbiko/jsonassert)

`jsonassert` is a Go test assertion library for asserting JSON payloads.

## Installation

```bash
go get github.com/kinbiko/jsonassert
```

## Usage

Create a new `jsonassert.Asserter` in your test and use this to make assertions against your JSON payloads:

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

Notice that you can pass in `fmt.Sprintf` arguments after the expected JSON structure.

`Asserter.Assertf()` currently supports assertions against strings only.

### Check for presence only

Some properties of a JSON payload may be difficult to know in advance.
E.g. timestamps, UUIDs, or other randomly assigned values.

For these types of values, place the string `"<<PRESENCE>>"` as the expected value, and `jsonassert` will only verify that this key exists (i.e. the actual JSON has the expected key, and its value is not `null`).

For example:

```go
func TestWhatever(t *testing.T) {
    ja := jsonassert.New(t)
    // find some sort of payload
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

The above will fail your test, but:


```go
func TestWhatever(t *testing.T) {
    ja := jsonassert.New(t)
    // find some sort of payload
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

Will fail your tests because the `time` key was not present in the actual JSON, and the `uuid` was `null`.

## Docs

You can find the [GoDocs for this package here](https://godoc.org/github.com/kinbiko/jsonassert).

## Contributing

Contributions are welcome. Please discuss feature requests in an issue before opening a PR.
