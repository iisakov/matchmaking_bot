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

func GetNumberMatches(s1, s2 []string) (result int) {
	var biggerS, smallerS []string
	if len(s1) > len(s2) {
		biggerS = s1
		smallerS = s2
	} else {
		biggerS = s2
		smallerS = s1
	}

	for _, vS := range smallerS {
		for _, vB := range biggerS {
			if vS == vB {
				result += 1
			}
		}
	}
	return
}

func CreateSlicePositiveInt(arg ...int) (result []int) {
	switch len(arg) {
	case 1:
		for i := 1; i < arg[0]; i++ {
			result = append(result, i)
		}
	case 2:

		if arg[0] < 0 {
			arg[0] = 1
		}
		for i := arg[0]; i < arg[1]; i++ {
			result = append(result, i)
		}
	case 3:
		if arg[0] < 0 {
			arg[0] = 1
		}
		for i := arg[0]; i < arg[1]; i += arg[2] {
			result = append(result, i)
		}
	}
	return result
}
