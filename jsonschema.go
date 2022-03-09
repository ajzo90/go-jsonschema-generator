/*
Basic json-schema generator based on Go types, for easy interchange of Go
structures between diferent languages.
*/
package jsonschema

import (
	"encoding/json"
	"reflect"
	"strings"
)

const DEFAULT_SCHEMA = "http://json-schema.org/schema#"

type Document struct {
	Schema string `json:"$schema,omitempty"`
	property
}

// Reads the variable structure into the JSON-Schema Document
func (d *Document) Read(variable interface{}) {
	if d.Schema == "" {
		d.Schema = DEFAULT_SCHEMA
	}

	value := reflect.ValueOf(variable)
	d.read(value.Type(), "")
}

// Marshal returns the JSON encoding of the Document
func (d *Document) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "    ")
}

// String return the JSON encoding of the Document as a string
func (d *Document) String() string {
	json, _ := d.Marshal()
	return string(json)
}

type property struct {
	typee                string
	Type                 []string             `json:"type,omitempty"`
	Format               string               `json:"format,omitempty"`
	Items                *property            `json:"items,omitempty"`
	Properties           map[string]*property `json:"properties,omitempty"`
	Required             []string             `json:"required,omitempty"`
	AdditionalProperties bool                 `json:"additionalProperties,omitempty"`
}

func (p *property) read(t reflect.Type, opts tagOptions) {
	jsType, format, kind := getTypeFromMapping(t)
	if jsType != "" {
		p.typee = jsType
	}
	if format != "" {
		p.Format = format
	}

	switch kind {
	case reflect.Slice:
		p.readFromSlice(t)
	case reflect.Map:
		p.readFromMap(t)
	case reflect.Struct:
		p.readFromStruct(t)
	case reflect.Ptr:
		p.read(t.Elem(), opts)
		if len(p.Type) > 0 {
			p.Type = append(p.Type, "null")
		}
		return
	}

	if len(p.typee) == 0 {
		return
	}
	p.Type = append(p.Type[:0], p.typee)
	p.typee = ""
}

func (p *property) readFromSlice(t reflect.Type) {
	jsType, _, kind := getTypeFromMapping(t.Elem())
	if kind == reflect.Uint8 {
		p.typee = "string"
	} else if jsType != "" {
		p.Items = &property{}
		p.Items.read(t.Elem(), tagOptions(""))
	}
}

func (p *property) readFromMap(t reflect.Type) {
	jsType, _, _ := getTypeFromMapping(t.Elem())

	if jsType != "" {
		p.Properties = make(map[string]*property, 0)
		var newProp property
		newProp.read(t.Elem(), "")
		p.Properties[".*"] = &newProp
	} else {
		p.AdditionalProperties = true
	}
}

func (p *property) readFromStruct(t reflect.Type) {
	p.typee = "object"
	p.Properties = make(map[string]*property, 0)
	p.AdditionalProperties = false

	count := t.NumField()
	for i := 0; i < count; i++ {
		field := t.Field(i)

		tag := field.Tag.Get("json")
		name, opts := parseTag(tag)
		if name == "" {
			name = field.Name
		}
		if name == "-" {
			continue
		}

		if field.Anonymous {
			embeddedProperty := &property{}
			embeddedProperty.read(field.Type, opts)

			for name, property := range embeddedProperty.Properties {
				p.Properties[name] = property
			}
			p.Required = append(p.Required, embeddedProperty.Required...)

			continue
		}

		p.Properties[name] = &property{}
		p.Properties[name].read(field.Type, opts)

		if !opts.Contains("omitempty") {
			p.Required = append(p.Required, name)
		}
	}
}

func getTypeFromMapping(t reflect.Type) (string, string, reflect.Kind) {

	switch v := t.Kind(); v {
	case reflect.Ptr:
		return "", "", v
	case reflect.Slice:
		return "array", "", v
	case reflect.Map, reflect.Struct:
		switch t.String() {
		case "time.Time":
			return "string", "date-time", reflect.String
		}
		return "object", "", v
	case reflect.String:
		return "string", "", v
	case reflect.Int8:
		return "integer", "i8", v
	case reflect.Int16:
		return "integer", "i16", v
	case reflect.Int32:
		return "integer", "i32", v
	case reflect.Int64, reflect.Int:
		return "integer", "i64", v
	case reflect.Uint8:
		return "integer", "u8", v
	case reflect.Uint16:
		return "integer", "u16", v
	case reflect.Uint32:
		return "integer", "u32", v
	case reflect.Uint64, reflect.Uint:
		return "integer", "u64", v
	case reflect.Float32:
		return "number", "f32", v
	case reflect.Float64:
		return "number", "f64", v
	case reflect.Bool:
		return "boolean", "", v
	case reflect.Interface:
		return "", "", v
	default:
		return "", "", v
	}
}

type tagOptions string

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}

	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
