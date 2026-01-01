package enumerable

import (
	"cmp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromList(t *testing.T) {
	a := assert.New(t)
	in := []string{"foo", "bar", "baz"}

	a.Equal(in, Collect(FromList(in, false)))
	a.Equal(in, Collect(FromList(in, true)))
}

func TestFromMap(t *testing.T) {
	a := assert.New(t)
	in := map[string]string{"foo": "bar", "asdf": "aoeu", "yourbases": "belongtous"}

	a.Equal(in, CollectMap(FromMap(in)))
}

func TestFromFunc(t *testing.T) {
	a := assert.New(t)
	result := []int{0, 1, 2, 3, 4}
	getFunc := func() func() (int, bool) {
		n := 0
		max := 5

		return func() (int, bool) {
			n++
			return n - 1, (n - 1) < max
		}
	}

	a.Equal(result, Collect(FromFunc(getFunc())))
}

func TestFlatten(t *testing.T) {
	a := assert.New(t)
	result := []int{1, 2, 3, 4, 5, 6}
	in := [][]int{{1, 2}, {}, {3}, {}, {4, 5}, {}, {6}}

	a.Equal(result, Collect(Flatten(FromList(in, false))))
}

func TestMap(t *testing.T) {
	a := assert.New(t)
	result := []float64{2, 4, 6, 8, 10, 12}
	in := []int{1, 2, 3, 4, 5, 6}

	a.Equal(result, Collect(Map(FromList(in, false), func(i int) float64 {
		return float64(i * 2)
	})))
}

func TestZip(t *testing.T) {
	a := assert.New(t)
	result := []AB[int, float64]{
		{1, 2},
		{2, 4},
		{3, 6},
	}

	inA := []int{1, 2, 3}
	inB := []float64{2, 4, 6}

	a.Equal(result, Collect(Zip(FromList(inA, true), FromList(inB, true))))
}

func TestFilter(t *testing.T) {
	a := assert.New(t)
	in := []int{1, 2, 3, 4, 5}
	out := []int{2, 4}

	a.Equal(out, Collect(Filter(FromList(in, false), func(i int) bool {
		return i%2 == 0
	})))
}

func invertSort[T cmp.Ordered](l, r T) int {
	if l > r {
		return -1 // inverse sort
	} else if l < r {
		return 1
	}

	return 0
}

func TestSort(t *testing.T) {
	a := assert.New(t)
	inOrdered := []int{10, 5, 15, 0}
	outOrdered := []int{0, 5, 10, 15}

	a.Equal(outOrdered, Collect(Sort(FromList(inOrdered, false))))

	outFunc := []int{15, 10, 5, 0}

	a.Equal(outFunc, Collect(SortFunc(FromList(inOrdered, false), invertSort)))
}

func TestChain(t *testing.T) {
	a := assert.New(t)
	in := []int{0, 2, 1, 4, 3}
	out := []float64{5, 10, 15, 20}

	// test a major set
	result := ExecChain1to1[int, []float64](
		FromList(in, false),
		ChainSort[int](),
		ChainFilter[int](func(i int) bool {
			return i != 0
		}),
		ChainZip[int, int](FromFunc(func() (int, bool) {
			return 5, true
		})),
		ChainMap[AB[int, int], int](func(i AB[int, int]) int {
			return i.A * i.B
		}),
		ChainMap[int, float64](func(i int) float64 {
			return float64(i)
		}),
		ChainCollect[float64](),
	)

	a.Equal(out, result)

	toFlatten := [][]int{{5, 10}, {2, 7}, {0, 8}}
	afterFlatten := 32

	// Test another subset
	fResult := ExecChain1to1[[]int, int](
		FromList(toFlatten, false),
		ChainFlatten[int](),
		ChainSortFunc[int](invertSort),
		ChainSum[int](func(i int, o int) int {
			o += i
			return o
		}),
	)

	a.Equal(afterFlatten, fResult)

	// Hit ForEach and 1to0
	toAdd := []int{1, 2, 3}
	afterForEach := make([]int, 0)

	ExecChain1to0(
		FromList(toAdd, false),
		ChainForEach(func(t int) {
			afterForEach = append(afterForEach, t)
		}))

	a.Equal(toAdd, afterForEach)

	// Test arbitrary
	arbitraryResult := ExecChainArbitrary(
		[]any{5},
		ChainFunc{Func: func(in int) (int, bool) {
			return in * 25, true
		}},
		ChainFunc{Func: func(in int, ok bool, newValue int) (int, bool) {
			return in / newValue, !ok
		}, Args: []any{5}})

	assert.Equal(t, []any{25, false}, arbitraryResult)
}
