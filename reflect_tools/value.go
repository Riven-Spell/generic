package reflect_tools

import "reflect"

func NewWriteableValue[T any]() reflect.Value {
	return reflect.New(reflect.TypeFor[T]()).Elem()
}

func WriteableValueOf[T any](in *T) reflect.Value {
	return reflect.ValueOf(in).Elem()
}
