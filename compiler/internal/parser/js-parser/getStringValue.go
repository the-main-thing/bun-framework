package jsparser

func GetStringWithQuotes(source *[]rune, cursor int) (string, int) {
	openingSymbol := '!'
	result := ""

	for i, char := range (*source)[cursor:] {
		if char == '\'' || char == '"' || char == '`' {
			if openingSymbol == '!' {
				openingSymbol = char
				result += string(char)
				continue
			}
			if openingSymbol != char {
				result += string(char)
				continue
			}
			result += string(char)
			if isEscaped(source, cursor+i) {
				continue
			}
			return result, min(cursor+i+1, len(*source)-1)
		}
		if openingSymbol != '!' {
			result += string(char)
		}
	}

	return "", len(*source) - 1
}

func isEscaped(source *[]rune, cursor int) bool {
	prevIndex := cursor - 1
	if prevIndex < 0 {
		return false
	}

	prevChar := (*source)[prevIndex]
	if prevChar == '\\' {
		return true && !isEscaped(source, prevIndex)
	}
	return false
}
