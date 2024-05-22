package stl

func IndexElemInSliseString(s []string, elem string) (int, bool) {
	for i, e := range s {
		if e == elem {
			return i, true
		}
	}
	return -1, false
}

func DeleteElementByIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
