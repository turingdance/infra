package slicekit

// Filter items that are not meet the condition
func Filter[T any](slice []T, condition func(item T, index int, slice []T) bool) []T {
	var filtered []T
	for index, item := range slice {
		if condition(slice[index], index, slice) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
