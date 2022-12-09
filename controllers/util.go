package controllers

// hasAnnotationWithValue returns true if the key exists in a and has value desiredValue.
func hasAnnotationWithValue(a map[string]string, key, desiredValue string) bool {
	if actualV, ok := a[key]; ok {
		return actualV == desiredValue
	}

	return false
}
