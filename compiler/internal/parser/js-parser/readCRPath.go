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
	for i, _ := range (*source)[start:] {
		char := (*source)[i]
		if char == '(' {
			start += i + 1
			reading = true
			break
		}
	}
	if !reading {
		return ""
	}

	openingSymbol := '!'
	for i, char := range (*source)[start:] {
		if char == ')' {
			return ""
		}
		if char == '\'' || char == '"' || char == '`' {
			if openingSymbol == '!' {
				openingSymbol = char
				continue
			}
			if openingSymbol != char {
				continue
			}
			if i == 1 {
				return ""
			}
			prev := (*source)[(start+i)-1]
			if prev == '\\' {
				continue
			}
			return string((*source)[start : start+i])
		}
	}

	return ""
}
