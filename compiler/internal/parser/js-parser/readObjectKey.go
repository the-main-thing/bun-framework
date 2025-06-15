package jsparser

import "errors"

func ReadObjectKey(input *[]byte, offset int) (start int, end int, eror error) {
	openingBytes := []byte{'[', '"', '\'', '`'}
	for _, b := range (*input)[offset:] {
		for index, openingByte := range openingBytes {
			if b == openingByte {
				start, end, err := ReadString(input, index+offset)
				if err != nil {
					return start, end, errors.New("Error reading object key: " + err.Error())
				}
				return start, end, nil
			}
		}

	}

	return 0, 0, errors.New("No object key found. Maybe this is a variable?")
}
