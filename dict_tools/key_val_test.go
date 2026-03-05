package dict_tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyValues(t *testing.T) {
	a := assert.New(t)
	origin := map[int]int{}
	for x := range 10 {
		origin[x] = x
	}

	for _, v := range Keys(origin) {
		_, ok := origin[v]
		a.True(ok, "value %d not found", v)
	}

	for _, v := range Values(origin) {
		// because values are the same, we can abuse that fact.
		_, ok := origin[v]
		a.True(ok, "value %d not found", v)
	}
}
