package jsparser

import "errors"

type ArrayOfStrings struct {
	StartIndex int
	EndIndex   int
}

// Will only work on flat string arrays. No mixed types, no nesting.
func ReadArrayOfStrings(input *[]byte, offset int) ([]ArrayOfStrings, error) {
	reading := false
	result := []ArrayOfStrings{}
	for i := offset; i < len(*input); i++ {
		b := (*input)[i]
		if reading {
			if IsStringOpeningSymbol(b) {
				start, end, err := ReadString(input, i)
				if err != nil {
					return nil, errors.New("Error reading array of strings: " + err.Error())
				}
				result = append(result, ArrayOfStrings{StartIndex: start, EndIndex: end})
				i = end + 1
				continue
			}
			if b == ',' || b == '\n' || b == ' ' || b == '\t' {
				continue
			}
			if b == ']' {
				break
			}
			return nil, errors.New("Unexpected character in array of strings: " + string(b))
		}
	}

	return result, nil
}
