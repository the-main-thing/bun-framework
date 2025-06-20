package jsparser

func ReadCRPath(source *[]rune, crFunctionAlias string) string {
	start := -1
	firstRuneOfCRFunction := rune(crFunctionAlias[0])
	for i := 0; i < len(*source); i++ {
		char := (*source)[i]
		if char == firstRuneOfCRFunction {
			if i+len(crFunctionAlias) >= len(*source) {
				return ""
			}
			if string((*source)[i:i+len(crFunctionAlias)]) == crFunctionAlias {
				start = i + len(crFunctionAlias)
				start = skipWhitespace(source, start)
				if start < 0 {
					return ""
				}
				if start >= len(*source)-1 {
					return ""
				}
				if (*source)[start] != '(' {
					i = start
					start = -1
					continue
				}
			}
		}
	}
	if start < 0 {
		return ""
	}
	reading := false
	for i, char := range (*source)[start:] {
		if char == '(' {
			start += i + 1
			reading = true
			break
		}
	}
	if !reading {
		return ""
	}

	path, _ := GetStringWithQuotes(source, start)
	return path
}
