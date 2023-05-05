package strutil

func StringFind(array []string, str string) int {
	if str == "" {
		return -1
	}
	for index, s := range array {
		if str == s {
			return index
		}
	}
	return -1
}

func StringIn(str string, array []string) bool {
	return StringFind(array, str) > -1
}
