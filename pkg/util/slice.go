package util

func IsContain[T1 comparable](slice []T1, item T1) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
