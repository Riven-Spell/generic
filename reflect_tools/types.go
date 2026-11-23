package reflect_tools

import "reflect"

func CanAssign(LHS, RHS reflect.Value) bool {
	return RHS.Type().AssignableTo(LHS.Type())
}

func CanAssignGenericLHS[LHS any](RHS reflect.Value) bool {
	return RHS.Type().AssignableTo(reflect.TypeFor[LHS]())
}

func CanAssignAnyGeneric[LHS any](RHS any) bool {
	return reflect.TypeOf(RHS).AssignableTo(reflect.TypeFor[LHS]())
}

// UnwrapType repeatedly calls Elem on t until it is no longer a pointer.
// Unlike reflect.Value, reflect.Type has no element type for interfaces
// (Type.Elem panics on Interface), so interface types are returned as-is.
func UnwrapType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
