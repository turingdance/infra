package slicekit

// 泛型去重函数，支持所有可比较类型（comparable）
func Unique[T comparable](s []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
