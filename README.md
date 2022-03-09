go-jsonschema-generator [![Build Status](https://img.shields.io/github/workflow/status/ajzo90/go-jsonschema-generator/Test.svg)](https://github.com/ajzo90/go-jsonschema-generator/actions) [![GoDoc](http://godoc.org/github.com/ajzo90/go-jsonschema-generator?status.png)](https://pkg.go.dev/github.com/ajzo90/go-jsonschema-generator)
==============================
Basic [json-schema](http://json-schema.org/) generator based on Go types, for easy interchange of Go structures across languages.

Fork of github.com/mcuadros/go-jsonschema-generator. Thank you very much. 
Updated and refactored to handle null from pointer fields.
* uses `[]string` instead of `string` for the `type` field, and append "null" in case the field is a pointer.
* adds format for known sized numbers, `{"type":"integer", "format": "i16"}`. Supported formats: `u8,u16,u32,u64,i8,i16,i32,i64,f32,f64`.  

Installation
------------

The recommended way to install go-jsonschema-generator

```
go get github.com/ajzo90/go-jsonschema-generator
```

Examples
--------

A basic example:

```go
package main

import (
	"fmt"

	"github.com/ajzo90/go-jsonschema-generator"
)

type EmbeddedType struct {
	Zoo *string
}

type Item struct {
	Value string
}

type ExampleBasic struct {
	Foo bool   `json:"foo"`
	Bar string `json:",omitempty"`
	Qux *int8
	Baz []string
	EmbeddedType
	List []Item
}

func main() {
	fmt.Println(jsonschema.New(ExampleBasic{}))
}
```

```json
{"$schema":"http://json-schema.org/schema#","type":["object"],"properties":{"Bar":{"type":["string"]},"Baz":{"type":["array"],"items":{"type":["string"]}},"List":{"type":["array"],"items":{"type":["object"],"properties":{"Value":{"type":["string"]}},"required":["Value"]}},"Qux":{"type":["integer","null"],"format":"i8"},"Zoo":{"type":["string","null"]},"foo":{"type":["boolean"]}},"required":["foo","Qux","Baz","Zoo","List"]}

```

License
-------

MIT, see [LICENSE](LICENSE)
