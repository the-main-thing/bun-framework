package jsparser

import "unicode"

func isValidVariableFirstChar(char rune) bool {
	return unicode.IsLetter(char) || char == '_' || char == '$'
}

func isValidVariableNonFirstChar(char rune) bool {
	return isValidVariableFirstChar(char) || unicode.IsDigit(char)
}
