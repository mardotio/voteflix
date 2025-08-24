package utils

func IsEqual[T comparable](v1 *T, v2 *T) bool {
	if v1 != nil && v2 != nil {
		return *v1 == *v2
	}

	return v1 == nil && v2 == nil
}
