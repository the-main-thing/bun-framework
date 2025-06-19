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
				endIndex = max(i + offset - 1, 0)
				break
			}
			openingSymbol = b
			startIndex = min(i + offset + 1, len(*input)-1)
			continue
		}
	}

	if len(*input) > 0 && openingSymbol == 0 {
		return 0, 0, errors.New("No string opening symbol found. Maybe this is a variable?")
	}

	if openingSymbol != 0 && !stringClosed {
		return 0, 0, errors.New("String is not closed whith the same opening symbol")
	}

	if endIndex < startIndex {
		panic("Logical error in ReadString: endIndex is smaller than startIndex")
	}

	return startIndex, endIndex, nil
}

func IsStringOpeningSymbol(b byte) bool {
	return b == '\'' || b == '"' || b == '`'
}
