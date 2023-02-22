package helper

func ArrayContains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func ArrayContainsSome(arr []string, val ...string) bool {
	for _, v := range val {
		if ArrayContains(arr, v) {
			return true
		}
	}
	return false
}
