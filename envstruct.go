// Package envstruct populates struct fields from environment variables.
//
// It is a drop-in replacement for kelseyhightower/envconfig with the same
// Process(prefix, &spec) function signature.
package envstruct

import (
	"errors"
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

// isStructField returns true if the field type is a struct that should be
// recursed into (not a special type like time.Duration or url.URL, and not
// implementing Decoder/Setter/TextUnmarshaler).
func isStructField(ft reflect.Type) bool {
	// Unwrap pointer.
	for ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	if ft.Kind() != reflect.Struct {
		return false
	}
	// Special types treated as scalar.
	if ft == durationType || ft == urlType {
		return false
	}
	// Check if it implements decode interfaces (use pointer type).
	pt := reflect.PointerTo(ft)
	if pt.Implements(decoderType) || pt.Implements(setterType) || pt.Implements(textUnmarshalerType) {
		return false
	}
	return true
}

func processStruct(prefix string, rv reflect.Value) error {
	rt := rv.Type()
	var errs []error
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

		// Check if this is a nested struct field.
		if isStructField(f.Type) {
			// Determine the nested prefix.
			nestedPrefix := spec.Name
			if f.Anonymous {
				// Embedded struct: flatten (no prefix added).
				nestedPrefix = prefix
			} else if prefix != "" {
				nestedPrefix = prefix + "_" + spec.Name
			}
			nestedPrefix = strings.ToUpper(nestedPrefix)

			if f.Type.Kind() == reflect.Ptr {
				// Pointer-to-struct: allocate temp, recurse, assign only if
				// at least one env var was set.
				tmp := reflect.New(f.Type.Elem())
				if err := processStruct(nestedPrefix, tmp.Elem()); err != nil {
					errs = append(errs, err)
				} else {
					fv.Set(tmp)
				}
			} else {
				if err := processStruct(nestedPrefix, fv); err != nil {
					errs = append(errs, err)
				}
			}
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
				errs = append(errs, &RequiredError{FieldName: f.Name, EnvVar: key})
				continue
			}
		}

		if !found {
			continue
		}

		// Decode and set the field value.
		if err := decode(fv, val, f.Name, key); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
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

