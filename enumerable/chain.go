package enumerable

import (
	"fmt"
	"reflect"
)

// ChainFunc defines a func to be called in ExecChain1to1 or ExecChain1to0.
type ChainFunc struct {
	Func any
	Args []any
}

// ChainFuncRaw returns a ChainFunc formatted as expected. No type insurance is provided.
func ChainFuncRaw(f any, args ...any) ChainFunc {
	return ChainFunc{f, args}
}

func executeChain(in []any, funcs ...ChainFunc) []any {
	toReflectVals := func(i []any) []reflect.Value {
		return Collect(Map(FromList(i, false), func(i any) reflect.Value {
			return reflect.ValueOf(i)
		}))
	}

	cIn := toReflectVals(in)
	for chainStep, v := range funcs {
		fType := reflect.TypeOf(v.Func)

		if fType.Kind() != reflect.Func { // Make sure we're working with a func
			panic(fmt.Sprintf("on chain step (0-indexed) %d: expected a func, instead found %s",
				chainStep, fType.Kind()))
		}

		nIn := fType.NumIn() // Ensure we have the same number of inputs in total
		if argc := len(cIn) + len(v.Args); nIn != argc {
			panic(fmt.Sprintf("on chain step (0-indexed) %d: received %d (%d (last out) + %d (args provided)) values, but the next chain func expects %d args",
				chainStep, argc, len(cIn), len(v.Args), nIn))
		}

		inList := append(cIn, toReflectVals(v.Args)...)
		for k := range inList {
			if !inList[k].Type().AssignableTo(fType.In(k)) {
				panic(fmt.Sprintf("on chain step (0-indexed) %d, arg %d: received type %s, but expected %s",
					chainStep, k, inList[k].Type().String(), fType.In(k).String()))
			}
		}

		cIn = reflect.ValueOf(v.Func).Call(inList)
	}

	return Collect(Map(FromList(cIn, false), func(i reflect.Value) any {
		return i.Interface()
	}))
}
