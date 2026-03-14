// Package envstruct populates struct fields from environment variables.
//
// It is a drop-in replacement for kelseyhightower/envconfig with the same
// Process(prefix, &spec) function signature.
package envstruct

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// Process populates the struct pointed to by spec with values from
// environment variables. The prefix is prepended to each field name
// (or tag override) when looking up environment variables.
func Process(prefix string, spec interface{}) error {
	rv := reflect.ValueOf(spec)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("envstruct: spec must be a non-nil pointer to a struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("envstruct: spec must be a pointer to a struct")
	}
	return processStruct(prefix, rv)
}

func processStruct(prefix string, rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fv := rv.Field(i)

		// Skip unexported fields.
		if !f.IsExported() {
			continue
		}

		// Build the env var name component from the field name.
		envName := strings.ToUpper(f.Name)
		spec := parseTag(f, envName)
		if spec.Ignored {
			continue
		}

		// Build full env var key.
		key := spec.Name
		if prefix != "" {
			key = prefix + "_" + spec.Name
		}
		key = strings.ToUpper(key)

		// Look up value.
		val, found := os.LookupEnv(key)

		if !found {
			if spec.HasDefault {
				val = spec.DefaultValue
				found = true
			} else if spec.Required {
				return &RequiredError{FieldName: f.Name, EnvVar: key}
			}
		}

		if !found {
			continue
		}

		// Decode and set the field value.
		if err := decode(fv, val, f.Name, key); err != nil {
			return err
		}
	}
	return nil
}

// MustProcess is like Process but panics on error.
func MustProcess(prefix string, spec interface{}) {
	if err := Process(prefix, spec); err != nil {
		panic(err)
	}
}

// Usage writes a usage message describing the environment variables
// expected by spec to the given writer.
func Usage(prefix string, spec interface{}, out io.Writer) error {
	return nil
}
