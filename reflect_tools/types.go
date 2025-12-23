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
