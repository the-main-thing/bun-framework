package jsparser

func PrevByte(input *[]byte, index int) byte {
	if index == 0 {
		return 0
	}
	return (*input)[index-1]
}
