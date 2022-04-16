# Contributing

## PRs

If you wish to contribute a change to this project, please create an issue first to discuss.
If you do not raise an issue before submitting a (significant) PR then your PR may be dismissed without much consideration.

PRs only improving the documentation are welcome without raising an issue first.

Ensure that:

1. The tests pass.
1. The linter has 0 issues.
1. You don't introduce any new dependencies (ask in an issue if you feel strongly that it's necessary).
1. You follow the existing commit message convention.

## Goals of `jsonassert`

- Accurately solve the problem of "Are these two JSONs semantically the same?"
- Be easy to comprehend
- Be easy to maintain
- Have no dependencies outside of the standard library
- Be well tested

### Non-goals

- Performance. There's a lot of unnecessary back and forth, mainly for the purpose of maintainability. If performance is crucial to you, feel free to fork the repo. Changes which don't negatively impact the readability and maintainability of the project are likely to get pulled if so requested, otherwise I'm happy to provide links to more performant forks in the README of this project.

## Structure

The `exports.go` file contains all the exported (publicly available) code. This file is the entry point to this entire package, and should be kept well documented and sparsely coded. The primary method is the `Assertf` function, which calls `pathassert`.

`pathassertf` is triggers the main algorithm:

- Keep track of where in the JSON we're at with the first string arg. This will get further appended onto for each nested object or array traversal (breadth-first).
- Checks that both representations are in fact JSON
- Checks if the string reps of the two JSON are literally equal.
- Checks if the types of the string reps are unequal, if they are, print and abort.
- Calls a specific `check*` function for each of the possible JSON types against the extracted value (from JSON string to Go). These extraction and check functions live in the file named after each type.
- If the type is object or array, each sub-property (key or element) will be compared by turning the value into a JSON string again, and calling `pathassertf` with an appropriately appended location inside each JSON.
