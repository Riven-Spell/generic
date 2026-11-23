package reflect_tools

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnwrapValueStripsPointersAndInterfaces(t *testing.T) {
	t.Run("non-pointer", func(t *testing.T) {
		assert.Equal(t, 42, UnwrapValue(reflect.ValueOf(42)).Interface())
	})

	t.Run("pointer", func(t *testing.T) {
		value := 42
		assert.Equal(t, 42, UnwrapValue(reflect.ValueOf(&value)).Interface())
	})

	t.Run("pointer chain", func(t *testing.T) {
		value := 42
		ptr := &value
		ptrPtr := &ptr
		assert.Equal(t, 42, UnwrapValue(reflect.ValueOf(ptrPtr)).Interface())
	})

	t.Run("interface holding pointer chain", func(t *testing.T) {
		value := 42
		ptr := &value
		var wrapped any = &ptr
		assert.Equal(t, 42, UnwrapValue(reflect.ValueOf(wrapped)).Interface())
	})

	t.Run("interface holding value", func(t *testing.T) {
		var wrapped any = "hello"
		assert.Equal(t, "hello", UnwrapValue(reflect.ValueOf(wrapped)).Interface())
	})
}

func TestUnwrapValueReturnsInvalidForNil(t *testing.T) {
	t.Run("invalid value", func(t *testing.T) {
		assert.False(t, UnwrapValue(reflect.Value{}).IsValid())
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		assert.False(t, UnwrapValue(reflect.ValueOf(ptr)).IsValid())
	})

	t.Run("nil interface", func(t *testing.T) {
		holder := struct{ Any any }{}
		value := reflect.ValueOf(&holder).Elem().Field(0)
		require.Equal(t, reflect.Interface, value.Kind())
		require.True(t, value.IsNil())
		assert.False(t, UnwrapValue(value).IsValid())
	})
}

func TestUnwrapValuePreservesSettability(t *testing.T) {
	value := 42
	unwrapped := UnwrapValue(reflect.ValueOf(&value))
	require.True(t, unwrapped.CanSet())

	unwrapped.SetInt(7)
	assert.Equal(t, 7, value)
}
