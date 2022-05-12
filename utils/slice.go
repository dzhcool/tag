package utils

func StringInSlice(arr []string, str string) bool {
	if len(arr) <= 0 {
		return false
	}
	for i := range arr {
		if arr[i] == str {
			return true
		}
	}
	return false
}
