package dict_tools

func SafeAdd[A comparable, B any](target map[A]B, key A, value B) map[A]B {
	if target == nil {
		target = map[A]B{
			key: value,
		}

		return target
	}

	target[key] = value
	return target
}
