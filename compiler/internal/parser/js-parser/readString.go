package jsparser

import "errors"

func ReadString(input *[]byte, offset int) (start int, end int, error error) {
	var openingSymbol byte
	var startIndex int
	var endIndex int
	stringClosed := false
	for i, b := range (*input)[offset:] {
		if IsStringOpeningSymbol(b) {
			if openingSymbol == b {
				prev := PrevByte(input, i+offset)
				if prev == '\\' {
					continue
				}
				stringClosed = true
				endIndex = i + offset
				break
			}
			openingSymbol = b
			startIndex = i + offset
			continue
		}
	}

	if len(*input) > 0 && openingSymbol == 0 {
		return startIndex, endIndex, errors.New("No string opening symbol found. Maybe this is a variable?")
	}

	if openingSymbol != 0 && !stringClosed {
		return startIndex, endIndex, errors.New("String is not closed whith the same opening symbol")
	}

	return startIndex, endIndex, nil
}

func IsStringOpeningSymbol(b byte) bool {
	return b == '\'' || b == '"' || b == '`'
}
