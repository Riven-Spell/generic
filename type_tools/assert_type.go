package type_tools

// AssertType is a no-op function intended to be placed inside an init function (or something uncalled),
// to assert that the input is assignable to type T. This is most useful with interfaces,
// or asserting that functions match their intended signature.
func AssertType[T any](in T) any { return nil }
