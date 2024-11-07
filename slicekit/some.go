package slicekit

func Some[T comparable](arr []T, fun func(arr []T, ele T) bool) bool {
	for _, v := range arr {
		if fun(arr, v) {
			return true
		}
	}
	return false
}
