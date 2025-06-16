package jsparser

type ContentType int

const (
	CODE ContentType = iota
	SINGLE_LINE_COMMENT
	MULTILINE_COMMENT
	STRING
)

func RemoveComments(tsconfig *[]byte) {
	var state ContentType
	state = CODE
	var stringOpeningSymbol byte
	for i, b := range *tsconfig {
		switch state {
		case STRING:
			if b == stringOpeningSymbol {
				// i is 100% bigger than 0, because we can't enter STRING state before 0 index
				if (*tsconfig)[i-1] == '\\' {
					break
				}
				state = CODE
			}
			continue
		case MULTILINE_COMMENT:
			(*tsconfig)[i] = ' '
			if b == '*' {
				nextIndex := i + 1
				if nextIndex < len(*tsconfig) {
					nextChar := (*tsconfig)[nextIndex]
					if nextChar == '/' {
						(*tsconfig)[nextIndex] = ' '
						state = CODE
					}
				}
			}
			break

		case SINGLE_LINE_COMMENT:
			(*tsconfig)[i] = ' '
			if b == '\n' {
				state = CODE
			}
			break

		case CODE:
			if b == '/' {
				nextIndex := i + 1
				if nextIndex >= len(*tsconfig) {
					(*tsconfig)[i] = ' '
					state = SINGLE_LINE_COMMENT
					break
				}
				nextChar := (*tsconfig)[nextIndex]
				if nextChar == '/' {
					(*tsconfig)[i] = ' '
					state = SINGLE_LINE_COMMENT
					break
				}
				if nextChar == '*' {
					(*tsconfig)[i] = ' '
					state = MULTILINE_COMMENT
					break
				}
				break
			}
			if b == '\'' || b == '"' || b == '`' {
				stringOpeningSymbol = b
				state = STRING
			}
			break
		default:
			panic("Unknown state")
		}

	}
}
