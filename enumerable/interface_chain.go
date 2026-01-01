package enumerable

import (
	"cmp"
)

// ExecChain1to1 is meant to be syntactic sugar for the other functions. Unfortunately, types cannot be ensured at compile time,
// so they are ensured at runtime. Validate that the first input of each ChainFunc matches the output of the prior ChainFunc
// (or I on the first ChainFunc, and that output on the last ChainFunc matches O).
// Use ChainFunc1to1, ChainFunc2to1, ChainFunc3to1 to quickly and easily define ChainFunc with some semblance of type insurance.
// If a ChainFunc has multiple outputs, they will be used as the first inputs of the next function.
// ExecChain1to1 expects one input, and one output.
func ExecChain1to1[I, O any](in Interface[I], funcs ...ChainFunc) O {
	out := executeChain([]any{in}, funcs...)
	return out[0].(O)
}

// ExecChain1to0 takes an interface in and returns nothing. Good for ending in ChainForEach.
func ExecChain1to0[I any](in Interface[I], funcs ...ChainFunc) {
	executeChain([]any{in}, funcs...)
}

// ExecChainArbitrary follows the same semantics as ExecChain1to1, but takes arbitrary inputs and arbitrary outputs.
// It is only really here to enable implementing custom arrangements of ChainXtoY, and isn't intended for direct use.
func ExecChainArbitrary(in []any, funcs ...ChainFunc) []any {
	return executeChain(in, funcs...)
}

// ChainZip makes no assurances about the length of either Interface.
// ChainZip will return false when either A or B returns false.
// ChainZip aligns two lists from their 0th index.
func ChainZip[A, B any](Right Interface[B]) ChainFunc {
	return ChainFuncRaw(Zip[A, B], Right)
}

// ChainMap converts an input value to an output value.
func ChainMap[I, O any](operator func(I) O) ChainFunc {
	return ChainFuncRaw(Map[I, O], operator)
}

// ChainFilter returns only the values in which operator returns true to.
func ChainFilter[I any](operator func(I) bool) ChainFunc {
	return ChainFuncRaw(Filter[I], operator)
}

// ChainFlatten reduces a incoming set of []Ts down to a single set of Ts.
func ChainFlatten[T any]() ChainFunc {
	return ChainFuncRaw(Flatten[T])
}

// ChainSortFunc sorts an incoming interface down to a set of single Ts. On the first Next(), this incurs a Collect() in order to sort.
// sortFunc(a, b) should return a negative number when a < b, a positive number when
// a > b and zero when a == b or a and b are incomparable in the sense of
// a strict weak ordering.
func ChainSortFunc[T any](sortFunc func(l, r T) int) ChainFunc {
	return ChainFuncRaw(SortFunc[T], sortFunc)
}

// ChainSort incurs the same semantics as SortFunc, but relies upon cmp.Ordered to provide the ordering.
func ChainSort[T cmp.Ordered]() ChainFunc {
	return ChainFuncRaw(Sort[T])
}

// ChainForEach fully executes an incoming Interface returning no values, and performing the requested operation on each value
func ChainForEach[T any](do func(T)) ChainFunc {
	return ChainFuncRaw(ForEach[T], do)
}

// ChainSum reduces an incoming Interface down to a single value, passing the next item and the current sum in on each call of operator.
func ChainSum[I, O any](operator func(I, O) O) ChainFunc {
	return ChainFuncRaw(Sum[I, O], operator)
}

// ChainCollect reduces an incoming interface to a []O.
func ChainCollect[O any]() ChainFunc {
	return ChainFuncRaw(Collect[O])
}
