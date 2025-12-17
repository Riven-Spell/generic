package ptr_tools

import (
	"github.com/Riven-Spell/generic/internal"
)

func EnsureNotNil[T any](in *T, Default T) *T {
	return internal.EnsureNotNil(in, Default)
}
