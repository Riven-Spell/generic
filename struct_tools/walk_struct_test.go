package struct_tools

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type walkLeaf struct {
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count"`
}

type walkMeta struct {
	Tag string `json:"tag"`
}

type walkOuter struct {
	Ptr   *walkLeaf `json:"ptr"`
	Meta  walkMeta  `json:"meta"`
	Label string    `json:"label"`
}

type walkHolder struct {
	Any any `json:"any"`
}

type walkDoublePointer struct {
	Ptr **walkLeaf `json:"ptr"`
}

type walkTriplePointer struct {
	Ptr ***walkLeaf `json:"ptr"`
}

func TestWalkStructTypeFillsPointers(t *testing.T) {
	var paths []string

	err := WalkStructType(reflect.TypeOf(walkOuter{}), func(path string, field reflect.StructField) error {
		paths = append(paths, path)

		switch path {
		case "Ptr":
			assert.Equal(t, reflect.Pointer, field.Type.Kind())
			assert.Equal(t, `json:"ptr"`, string(field.Tag))
		case "Ptr.Name":
			assert.Equal(t, "Name", field.Name)
			assert.Equal(t, reflect.TypeOf(""), field.Type)
			assert.Equal(t, `json:"name" validate:"required"`, string(field.Tag))
		case "Ptr.Count":
			assert.Equal(t, "Count", field.Name)
			assert.Equal(t, reflect.TypeOf(0), field.Type)
			assert.Equal(t, `json:"count"`, string(field.Tag))
		case "Meta":
			assert.Equal(t, `json:"meta"`, string(field.Tag))
		case "Meta.Tag":
			assert.Equal(t, `json:"tag"`, string(field.Tag))
		case "Label":
			assert.Equal(t, `json:"label"`, string(field.Tag))
		default:
			t.Errorf("unexpected path %q", path)
		}

		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"Ptr", "Ptr.Name", "Ptr.Count", "Meta", "Meta.Tag", "Label"}, paths)
}

func TestWalkStructTypeFillsPointerChains(t *testing.T) {
	tests := map[string]struct {
		typ      reflect.Type
		wantType reflect.Type
	}{
		"double pointer": {
			typ:      reflect.TypeOf(walkDoublePointer{}),
			wantType: reflect.TypeOf((**walkLeaf)(nil)),
		},
		"triple pointer": {
			typ:      reflect.TypeOf(walkTriplePointer{}),
			wantType: reflect.TypeOf((***walkLeaf)(nil)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var paths []string

			err := WalkStructType(test.typ, func(path string, field reflect.StructField) error {
				paths = append(paths, path)

				switch path {
				case "Ptr":
					assert.Equal(t, test.wantType, field.Type)
					assert.Equal(t, `json:"ptr"`, string(field.Tag))
				case "Ptr.Name":
					assert.Equal(t, "Name", field.Name)
					assert.Equal(t, `json:"name" validate:"required"`, string(field.Tag))
				case "Ptr.Count":
					assert.Equal(t, "Count", field.Name)
					assert.Equal(t, `json:"count"`, string(field.Tag))
				default:
					t.Errorf("unexpected path %q", path)
				}

				return nil
			})

			require.NoError(t, err)
			assert.Equal(t, []string{"Ptr", "Ptr.Name", "Ptr.Count"}, paths)
		})
	}
}

func TestWalkStructTypeDoesNotFillInterfaces(t *testing.T) {
	var paths []string

	err := WalkStructType(reflect.TypeOf(walkHolder{}), func(path string, field reflect.StructField) error {
		paths = append(paths, path)
		assert.Equal(t, reflect.Interface, field.Type.Kind())
		assert.Equal(t, `json:"any"`, string(field.Tag))
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"Any"}, paths)
}

func TestWalkStructValueFollowsInterfaceStructs(t *testing.T) {
	holder := walkHolder{Any: walkLeaf{Name: "leaf", Count: 3}}
	var paths []string
	var visited walkLeaf

	err := WalkStructValue(reflect.ValueOf(holder), func(path string, field reflect.StructField, value reflect.Value) error {
		paths = append(paths, path)

		switch path {
		case "Any":
			assert.Equal(t, reflect.Interface, field.Type.Kind())
			assert.Equal(t, `json:"any"`, string(field.Tag))
		case "Any.Name":
			assert.Equal(t, "Name", field.Name)
			assert.Equal(t, `json:"name" validate:"required"`, string(field.Tag))
			visited.Name = value.String()
		case "Any.Count":
			assert.Equal(t, "Count", field.Name)
			assert.Equal(t, `json:"count"`, string(field.Tag))
			visited.Count = int(value.Int())
		default:
			t.Errorf("unexpected path %q", path)
		}

		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"Any", "Any.Name", "Any.Count"}, paths)
	assert.Equal(t, walkLeaf{Name: "leaf", Count: 3}, visited)
}

func TestWalkStructValueScalarInterfaceIsNotFollowed(t *testing.T) {
	holder := walkHolder{Any: 42}
	var paths []string

	err := WalkStructValue(reflect.ValueOf(holder), func(path string, field reflect.StructField, value reflect.Value) error {
		paths = append(paths, path)
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"Any"}, paths)
}

func TestWalkStructValueFollowsInterfacePointerChains(t *testing.T) {
	leaf := &walkLeaf{Name: "leaf", Count: 3}
	double := &leaf

	tests := map[string]any{
		"pointer":        leaf,
		"double pointer": double,
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			holder := walkHolder{Any: value}
			var paths []string
			var visited walkLeaf

			err := WalkStructValue(reflect.ValueOf(holder), func(path string, field reflect.StructField, value reflect.Value) error {
				paths = append(paths, path)

				switch path {
				case "Any":
					assert.Equal(t, reflect.Interface, field.Type.Kind())
				case "Any.Name":
					visited.Name = value.String()
				case "Any.Count":
					visited.Count = int(value.Int())
				default:
					t.Errorf("unexpected path %q", path)
				}

				return nil
			})

			require.NoError(t, err)
			assert.Equal(t, []string{"Any", "Any.Name", "Any.Count"}, paths)
			assert.Equal(t, walkLeaf{Name: "leaf", Count: 3}, visited)
		})
	}
}

func TestWalkStructValueHandlesNilInterface(t *testing.T) {
	var holder walkHolder
	var paths []string

	err := WalkStructValue(reflect.ValueOf(holder), func(path string, field reflect.StructField, value reflect.Value) error {
		paths = append(paths, path)
		assert.True(t, value.IsNil())
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"Any"}, paths)
}

func TestWalkStructValueFollowsNonNilPointersOnly(t *testing.T) {
	t.Run("non-nil pointer", func(t *testing.T) {
		value := walkOuter{Ptr: &walkLeaf{Name: "leaf"}}
		var paths []string

		err := WalkStructValue(reflect.ValueOf(value), func(path string, field reflect.StructField, value reflect.Value) error {
			paths = append(paths, path)
			return nil
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"Ptr", "Ptr.Name", "Ptr.Count", "Meta", "Meta.Tag", "Label"}, paths)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var value walkOuter
		var paths []string

		err := WalkStructValue(reflect.ValueOf(value), func(path string, field reflect.StructField, value reflect.Value) error {
			paths = append(paths, path)
			return nil
		})

		require.NoError(t, err)
		assert.Nil(t, value.Ptr)
		assert.Equal(t, []string{"Ptr", "Meta", "Meta.Tag", "Label"}, paths)
	})

	t.Run("non-nil pointer chain", func(t *testing.T) {
		leaf := &walkLeaf{Name: "leaf", Count: 3}
		value := walkDoublePointer{Ptr: &leaf}
		var paths []string
		var visited walkLeaf

		err := WalkStructValue(reflect.ValueOf(value), func(path string, field reflect.StructField, value reflect.Value) error {
			paths = append(paths, path)

			switch path {
			case "Ptr.Name":
				visited.Name = value.String()
			case "Ptr.Count":
				visited.Count = int(value.Int())
			}

			return nil
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"Ptr", "Ptr.Name", "Ptr.Count"}, paths)
		assert.Equal(t, walkLeaf{Name: "leaf", Count: 3}, visited)
	})

	t.Run("nil pointer chain", func(t *testing.T) {
		var value walkDoublePointer
		var paths []string

		err := WalkStructValue(reflect.ValueOf(value), func(path string, field reflect.StructField, value reflect.Value) error {
			paths = append(paths, path)
			return nil
		})

		require.NoError(t, err)
		assert.Nil(t, value.Ptr)
		assert.Equal(t, []string{"Ptr"}, paths)
	})
}

func TestWalkSkipBranch(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		var paths []string

		err := WalkStructType(reflect.TypeOf(walkOuter{}), func(path string, field reflect.StructField) error {
			paths = append(paths, path)
			if path == "Ptr" {
				return ErrSkipBranch
			}
			return nil
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"Ptr", "Meta", "Meta.Tag", "Label"}, paths)
	})

	t.Run("value", func(t *testing.T) {
		value := walkOuter{Ptr: &walkLeaf{Name: "leaf"}, Meta: walkMeta{Tag: "meta"}}
		var paths []string

		err := WalkStructValue(reflect.ValueOf(value), func(path string, field reflect.StructField, value reflect.Value) error {
			paths = append(paths, path)
			if path == "Meta" {
				return ErrSkipBranch
			}
			return nil
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"Ptr", "Ptr.Name", "Ptr.Count", "Meta", "Label"}, paths)
	})
}

func TestWalkStopNow(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		var paths []string

		err := WalkStructType(reflect.TypeOf(walkOuter{}), func(path string, field reflect.StructField) error {
			paths = append(paths, path)
			return ErrStopNow
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"Ptr"}, paths)
	})

	t.Run("value", func(t *testing.T) {
		var paths []string

		err := WalkStructValue(reflect.ValueOf(walkOuter{}), func(path string, field reflect.StructField, value reflect.Value) error {
			paths = append(paths, path)
			return ErrStopNow
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"Ptr"}, paths)
	})
}

func TestWalkPropagatesActorError(t *testing.T) {
	errActor := errors.New("actor failed")
	var paths []string

	err := WalkStructValue(reflect.ValueOf(walkOuter{}), func(path string, field reflect.StructField, value reflect.Value) error {
		paths = append(paths, path)
		if path == "Meta.Tag" {
			return errActor
		}
		return nil
	})

	require.ErrorIs(t, err, errActor)
	assert.Equal(t, []string{"Ptr", "Meta", "Meta.Tag"}, paths)
}

func TestWalkRejectsNonStructs(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		called := false

		err := WalkStructValue(reflect.ValueOf(5), func(path string, field reflect.StructField, value reflect.Value) error {
			called = true
			return nil
		})

		require.Error(t, err)
		assert.False(t, called)
	})

	t.Run("pointer value", func(t *testing.T) {
		err := WalkStructValue(reflect.ValueOf(&walkOuter{}), func(path string, field reflect.StructField, value reflect.Value) error {
			return nil
		})

		require.Error(t, err)
	})

	t.Run("type", func(t *testing.T) {
		called := false

		err := WalkStructType(reflect.TypeOf(5), func(path string, field reflect.StructField) error {
			called = true
			return nil
		})

		require.NoError(t, err)
		assert.False(t, called)
	})
}
