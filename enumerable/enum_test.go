package enumerable

import (
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
