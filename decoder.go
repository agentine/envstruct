package envstruct

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Decoder is implemented by types that can decode themselves from a string.
type Decoder interface {
	Decode(value string) error
}

// Setter is implemented by types that can set themselves from a string.
// This interface exists for compatibility with kelseyhightower/envconfig.
type Setter interface {
	Set(value string) error
}

var (
	decoderType        = reflect.TypeOf((*Decoder)(nil)).Elem()
	setterType         = reflect.TypeOf((*Setter)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	durationType       = reflect.TypeOf(time.Duration(0))
	urlType            = reflect.TypeOf(url.URL{})
)

// decode sets a reflect.Value from a string value.
// It returns a ParseError if decoding fails.
func decode(fv reflect.Value, val string, fieldName string, envVar string) error {
	// Handle pointer: allocate and decode the element.
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}
		return decode(fv.Elem(), val, fieldName, envVar)
	}

	// Check interfaces on pointer-to-value (to catch pointer receivers).
	pv := fv.Addr()
	if pv.Type().Implements(decoderType) {
		if err := pv.Interface().(Decoder).Decode(val); err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		return nil
	}
	if pv.Type().Implements(setterType) {
		if err := pv.Interface().(Setter).Set(val); err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		return nil
	}
	if pv.Type().Implements(textUnmarshalerType) {
		if err := pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(val)); err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		return nil
	}

	// Special types.
	ft := fv.Type()
	if ft == durationType {
		d, err := time.ParseDuration(val)
		if err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		fv.SetInt(int64(d))
		return nil
	}
	if ft == urlType {
		u, err := url.Parse(val)
		if err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		fv.Set(reflect.ValueOf(*u))
		return nil
	}

	// Slices.
	if ft.Kind() == reflect.Slice {
		return decodeSlice(fv, val, fieldName, envVar)
	}

	// Maps.
	if ft.Kind() == reflect.Map {
		return decodeMap(fv, val, fieldName, envVar)
	}

	// Scalar types.
	return decodeScalar(fv, val, fieldName, envVar)
}

func decodeScalar(fv reflect.Value, val string, fieldName string, envVar string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 0, fv.Type().Bits())
		if err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		fv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 0, fv.Type().Bits())
		if err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		fv.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, fv.Type().Bits())
		if err != nil {
			return parseErr(fieldName, envVar, val, fv, err)
		}
		fv.SetFloat(n)
	default:
		return parseErr(fieldName, envVar, val, fv, fmt.Errorf("unsupported type %s", fv.Type()))
	}
	return nil
}

func decodeSlice(fv reflect.Value, val string, fieldName string, envVar string) error {
	parts := strings.Split(val, ",")
	slice := reflect.MakeSlice(fv.Type(), len(parts), len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if err := decode(slice.Index(i), part, fieldName, envVar); err != nil {
			return err
		}
	}
	fv.Set(slice)
	return nil
}

func decodeMap(fv reflect.Value, val string, fieldName string, envVar string) error {
	m := reflect.MakeMap(fv.Type())
	pairs := strings.Split(val, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return parseErr(fieldName, envVar, val, fv,
				fmt.Errorf("expected key=value pair, got %q", pair))
		}
		key := reflect.New(fv.Type().Key()).Elem()
		if err := decode(key, strings.TrimSpace(kv[0]), fieldName, envVar); err != nil {
			return err
		}
		value := reflect.New(fv.Type().Elem()).Elem()
		if err := decode(value, strings.TrimSpace(kv[1]), fieldName, envVar); err != nil {
			return err
		}
		m.SetMapIndex(key, value)
	}
	fv.Set(m)
	return nil
}

func parseErr(fieldName, envVar, val string, fv reflect.Value, err error) *ParseError {
	return &ParseError{
		FieldName: fieldName,
		EnvVar:    envVar,
		Value:     val,
		TypeName:  fv.Type().String(),
		Err:       err,
	}
}
