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
