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

## Docs

You can find the [GoDocs for this package here](https://godoc.org/github.com/kinbiko/jsonassert).

## Contributing

Contributions are welcome. Please discuss feature requests in an issue before opening a PR.
