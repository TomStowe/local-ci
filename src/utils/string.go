package string_utils

// In determines whether a string is in an array
func In(v string, l []string) bool {
	for _, i := range l {
		if i == v {
			return true
		}
	}

	return false
}
