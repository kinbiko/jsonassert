# Design

json is the fundamental type used throughout.
There should be methods for:
    - Identifying the type of a json
    - Getting a number out of a json
    - Getting a boolean out of a json
    - Getting an array out of a json
    - Getting a nested object out of a json
    - Getting an element of an array json


type json interface{
    getType()
    String()
}

type jsonArray interface {
    json
    get(int) (json, error)
}

type jsonObject interface {
    json
    get(string) (json, error)
}
