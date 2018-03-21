package testdeep

func ternRune(cond bool, a, b rune) rune {
	if cond {
		return a
	}
	return b
}

func ternStr(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
