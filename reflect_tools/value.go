package reflect_tools

import "reflect"

func NewWriteableValue[T any]() reflect.Value {
	return reflect.New(reflect.TypeFor[T]()).Elem()
}

func WriteableValueOf[T any](in *T) reflect.Value {
	return reflect.ValueOf(in).Elem()
}

// UnwrapValue repeatedly calls Elem on v until it is no longer a pointer or
// interface, and returns what remains. An invalid Value is returned for nil
// pointers or interfaces.
func UnwrapValue(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		v = v.Elem()
	}
	return v
}
