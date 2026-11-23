package reflect_tools

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrapTypeStripsPointers(t *testing.T) {
	tests := map[string]struct {
		in       reflect.Type
		expected reflect.Type
	}{
		"non-pointer":   {in: reflect.TypeOf(42), expected: reflect.TypeOf(42)},
		"pointer":       {in: reflect.TypeOf((*int)(nil)), expected: reflect.TypeOf(0)},
		"pointer chain": {in: reflect.TypeOf((***int)(nil)), expected: reflect.TypeOf(0)},
		"pointer to interface": {
			in:       reflect.TypeOf((*any)(nil)),
			expected: reflect.TypeOf((*any)(nil)).Elem(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, UnwrapType(test.in))
		})
	}
}

func TestUnwrapTypeLeavesInterfacesAlone(t *testing.T) {
	interfaceType := reflect.TypeOf((*any)(nil)).Elem()
	assert.Equal(t, interfaceType, UnwrapType(interfaceType))
}

func TestUnwrapTypeReturnsNilForNil(t *testing.T) {
	assert.Nil(t, UnwrapType(nil))
}
