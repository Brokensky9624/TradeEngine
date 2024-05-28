package tool

func FilterSlice[T any](sl []T, filterFunc func(el T) bool) []T {
	newSl := make([]T, 0)
	for _, el := range sl {
		if filterFunc(el) {
			newSl = append(newSl, el)
		}
	}
	return newSl
}

func MapSlice[T any](sl []T, mapFunc func(el T) T) []T {
	newSl := make([]T, len(sl))
	for i, el := range sl {
		newSl[i] = mapFunc(el)
	}
	return newSl
}
