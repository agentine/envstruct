package envstruct

import (
	"reflect"
	"strings"
)

// fieldSpec holds parsed struct tag information for a single field.
type fieldSpec struct {
	Name         string
	Required     bool
	DefaultValue string
	HasDefault   bool
	Ignored      bool
	Description  string
}

// parseTag extracts a fieldSpec from a struct field's tags.
// It checks env, then envconfig tags (for migration compatibility).
// The fieldName is the struct field name used as a fallback.
func parseTag(f reflect.StructField, fieldName string) fieldSpec {
	spec := fieldSpec{Name: fieldName}

	// Check env tag first, then envconfig for compat.
	tag, ok := f.Tag.Lookup("env")
	if !ok {
		tag, ok = f.Tag.Lookup("envconfig")
	}

	if ok {
		parts := strings.Split(tag, ",")
		name := parts[0]
		if name == "-" {
			spec.Ignored = true
			return spec
		}
		if name != "" {
			spec.Name = name
		}
		for _, opt := range parts[1:] {
			if opt == "required" {
				spec.Required = true
			}
		}
	}

	// Check default tag.
	if defVal, ok := f.Tag.Lookup("default"); ok {
		spec.DefaultValue = defVal
		spec.HasDefault = true
	}

	// Check desc tag.
	if desc, ok := f.Tag.Lookup("desc"); ok {
		spec.Description = desc
	}

	return spec
}
