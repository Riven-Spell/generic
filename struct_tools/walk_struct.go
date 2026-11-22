package struct_tools

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrSkipBranch = errors.New("skip branch")
	ErrStopNow    = errors.New("stop walk now")
)

// WalkStructType takes a reflect.Type referencing a struct, and walks the field(s) in all child structs.
// Pointer types are filled in, interfaces are not.
func WalkStructType(t reflect.Type, actor func(path string, field reflect.StructField) error) error {
	if t.Kind() != reflect.Struct {
		return nil
	}

	err := walkStructValueInternal(nil, reflect.New(t).Elem(), func(path string, field reflect.StructField, value reflect.Value) error {
		return actor(path, field)
	}, true)

	if errors.Is(err, ErrStopNow) {
		return nil
	}

	return err
}

// WalkStructValue walks a reflect.Value pointed at a struct, without filling in pointers and interfaces.
// Supplied interfaces' structs will be walked as well.
func WalkStructValue(val reflect.Value, actor func(path string, field reflect.StructField, value reflect.Value) error) error {
	if val.Kind() != reflect.Struct {
		return errors.New("WalkStructValue does not accept non-struct values")
	}

	err := walkStructValueInternal(nil, val, actor, false)

	if errors.Is(err, ErrStopNow) {
		return nil
	}

	return err
}

func walkStructValueInternal(path []string, val reflect.Value, actor func(path string, field reflect.StructField, value reflect.Value) error, rawType bool) error {
	for field, fVal := range val.Fields() {
		fPath := append(path, field.Name)

		err := actor(strings.Join(fPath, "."), field, fVal)
		if err != nil {
			if errors.Is(err, ErrSkipBranch) {
				err = nil
				continue
			}

			return err
		}

		for fVal.Kind() == reflect.Interface || fVal.Kind() == reflect.Pointer {
			if fVal.IsNil() {
				if fVal.Kind() == reflect.Pointer && rawType {
					fVal.Set(reflect.New(fVal.Type().Elem()))
					continue
				}

				break
			}

			fVal = fVal.Elem()
		}

		if fVal.Kind() == reflect.Struct {
			if err := walkStructValueInternal(fPath, fVal, actor, rawType); err != nil {
				return err
			}
		}
	}

	return nil
}
