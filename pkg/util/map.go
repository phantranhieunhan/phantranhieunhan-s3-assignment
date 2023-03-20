package	util

func MapValuesToSlice[T1 comparable, T2 any](m map[T1]T2) []T2 {
	s := make([]T2, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func MapKeysToSlice[T1 comparable, T2 any](m map[T1]T2) []T1 {
	s := make([]T1, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}
