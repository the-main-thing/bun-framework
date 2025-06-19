package jsparser

import "unicode"

func isKeyword(source *[]rune, offset int, keyword string) bool {
	if offset+len(keyword) >= len(*source) {
		return false
	}
	if string((*source)[offset:offset+len(keyword)]) != keyword {
		return false
	}
	return unicode.IsSpace((*source)[offset+len(keyword)])
}
