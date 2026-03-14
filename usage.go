package envstruct

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

func usageStruct(prefix string, rt reflect.Type, tw *tabwriter.Writer) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}

		envName := strings.ToUpper(f.Name)
		spec := parseTag(f, envName)
		if spec.Ignored {
			continue
		}

		ft := f.Type
		// Unwrap pointer for type display and struct check.
		isPtr := ft.Kind() == reflect.Ptr
		elemType := ft
		if isPtr {
			elemType = ft.Elem()
		}

		// Recurse into nested structs.
		if isStructField(ft) {
			nestedPrefix := spec.Name
			if f.Anonymous {
				nestedPrefix = prefix
			} else if prefix != "" {
				nestedPrefix = prefix + "_" + spec.Name
			}
			nestedPrefix = strings.ToUpper(nestedPrefix)
			if isPtr {
				usageStruct(nestedPrefix, elemType, tw)
			} else {
				usageStruct(nestedPrefix, ft, tw)
			}
			continue
		}

		// Build key.
		key := spec.Name
		if prefix != "" {
			key = prefix + "_" + spec.Name
		}
		key = strings.ToUpper(key)

		// Type name.
		typeName := ft.String()

		// Options column.
		var opts string
		if spec.Required {
			opts = "[required]"
		} else if spec.HasDefault {
			opts = fmt.Sprintf("[default: %s]", spec.DefaultValue)
		}

		fmt.Fprintf(tw, "  %s\t%s\t%s\t%s\n", key, typeName, opts, spec.Description)
	}
}

func writeUsage(prefix string, spec interface{}, out io.Writer) error {
	rv := reflect.ValueOf(spec)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("envstruct: spec must be a struct or pointer to struct")
	}
	tw := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	usageStruct(prefix, rv.Type(), tw)
	return tw.Flush()
}
