package genericutils

func MakeGenericWithDefault[T any]() T {
	var t T
	return t
}
