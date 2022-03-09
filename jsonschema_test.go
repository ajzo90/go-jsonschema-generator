package jsonschema

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

const schemaURL = "http://json-schema.org/schema#"

func Test(t *testing.T) { TestingT(t) }

type propertySuite struct{}

var _ = Suite(&propertySuite{})

type ExampleJSONBasic struct {
	Omitted    string  `json:"-,omitempty"`
	Bool       bool    `json:",omitempty"`
	BoolP      *bool   `json:",omitempty"`
	Integer    int     `json:",omitempty"`
	Integer8   int8    `json:",omitempty"`
	Integer16  int16   `json:",omitempty"`
	Integer32  int32   `json:",omitempty"`
	Integer64  int64   `json:",omitempty"`
	UInteger   uint    `json:",omitempty"`
	UInteger8  uint8   `json:",omitempty"`
	UInteger16 uint16  `json:",omitempty"`
	UInteger32 uint32  `json:",omitempty"`
	UInteger64 uint64  `json:",omitempty"`
	String     string  `json:",omitempty"`
	Bytes      []byte  `json:",omitempty"`
	Float32    float32 `json:",omitempty"`
	Float64    float64
	Interface  interface{}
	InterfaceP *interface{}
	Timestamp  time.Time `json:",omitempty"`
}

func (suit *propertySuite) TestLoad(c *C) {
	j := Document{}
	j.Read(ExampleJSONBasic{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type:     []string{"object"},
			Required: []string{"Float64", "Interface", "InterfaceP"},
			Properties: map[string]*property{
				"Bool":       {Type: []string{"boolean"}},
				"BoolP":      {Type: []string{"boolean", "null"}},
				"Integer":    {Type: []string{"integer"}, Format: "i64"},
				"Integer8":   {Type: []string{"integer"}, Format: "i8"},
				"Integer16":  {Type: []string{"integer"}, Format: "i16"},
				"Integer32":  {Type: []string{"integer"}, Format: "i32"},
				"Integer64":  {Type: []string{"integer"}, Format: "i64"},
				"UInteger":   {Type: []string{"integer"}, Format: "u64"},
				"UInteger8":  {Type: []string{"integer"}, Format: "u8"},
				"UInteger16": {Type: []string{"integer"}, Format: "u16"},
				"UInteger32": {Type: []string{"integer"}, Format: "u32"},
				"UInteger64": {Type: []string{"integer"}, Format: "u64"},
				"String":     {Type: []string{"string"}},
				"Bytes":      {Type: []string{"string"}},
				"Float32":    {Type: []string{"number"}, Format: "f32"},
				"Float64":    {Type: []string{"number"}, Format: "f64"},
				"Interface":  {},
				"InterfaceP": {},
				"Timestamp":  {Type: []string{"string"}, Format: "date-time"},
			},
		},
	})
}

type ExampleJSONBasicWithTag struct {
	Bool bool `json:"test"`
}

func (suit *propertySuite) TestLoadWithTag(c *C) {
	j := Document{}
	j.Read(ExampleJSONBasicWithTag{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type:     []string{"object"},
			Required: []string{"test"},
			Properties: map[string]*property{
				"test": {Type: []string{"boolean"}},
			},
		},
	})
}

type SliceStruct struct {
	Value string
}

type ExampleJSONBasicSlices struct {
	Slice            []string      `json:",foo,omitempty"`
	SliceOfInterface []interface{} `json:",foo"`
	SliceOfStruct    []SliceStruct
}

func (suit *propertySuite) TestLoadSliceAndContains(c *C) {
	j := Document{}
	j.Read(ExampleJSONBasicSlices{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type: []string{"object"},
			Properties: map[string]*property{
				"Slice": {
					Type:  []string{"array"},
					Items: &property{Type: []string{"string"}},
				},
				"SliceOfInterface": {
					Type: []string{"array"},
				},
				"SliceOfStruct": {
					Type: []string{"array"},
					Items: &property{
						Type:     []string{"object"},
						Required: []string{"Value"},
						Properties: map[string]*property{
							"Value": {
								Type: []string{"string"},
							},
						},
					},
				},
			},

			Required: []string{"SliceOfInterface", "SliceOfStruct"},
		},
	})
}

type ExampleJSONNestedStruct struct {
	Struct struct {
		Foo string
	}
}

func (suit *propertySuite) TestLoadNested(c *C) {
	j := Document{}
	j.Read(ExampleJSONNestedStruct{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type: []string{"object"},
			Properties: map[string]*property{
				"Struct": {
					Type: []string{"object"},
					Properties: map[string]*property{
						"Foo": {Type: []string{"string"}},
					},
					Required: []string{"Foo"},
				},
			},
			Required: []string{"Struct"},
		},
	})
}

type EmbeddedStruct struct {
	Foo string
}

type ExampleJSONEmbeddedStruct struct {
	EmbeddedStruct
}

func (suit *propertySuite) TestLoadEmbedded(c *C) {
	j := Document{}
	j.Read(ExampleJSONEmbeddedStruct{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type: []string{"object"},
			Properties: map[string]*property{
				"Foo": {Type: []string{"string"}},
			},
			Required: []string{"Foo"},
		},
	})
}

type ExampleJSONBasicMaps struct {
	Maps           map[string]string `json:",omitempty"`
	MapOfInterface map[string]interface{}
}

func (suit *propertySuite) TestLoadMap(c *C) {
	j := Document{}
	j.Read(ExampleJSONBasicMaps{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type: []string{"object"},
			Properties: map[string]*property{
				"Maps": {
					Type: []string{"object"},
					Properties: map[string]*property{
						".*": {Type: []string{"string"}},
					},
					AdditionalProperties: false,
				},
				"MapOfInterface": {
					Type:                 []string{"object"},
					AdditionalProperties: true,
				},
			},
			Required: []string{"MapOfInterface"},
		},
	})
}

func (suit *propertySuite) TestLoadNonStruct(c *C) {
	j := Document{}
	j.Read([]string{})

	c.Assert(j, DeepEquals, Document{
		Schema: schemaURL,
		property: property{
			Type:  []string{"array"},
			Items: &property{Type: []string{"string"}},
		},
	})
}
