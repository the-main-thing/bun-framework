package jsparser

import (
	"unicode"
)

type NamedImport struct {
	Alias string
	Name  string
}

type ImportStatement struct {
	StartIndex    int
	EndIndex      int
	SourcePath    string
	NamedImports  []NamedImport
	DefaultImport string
}

// we expect that all the comments are removed
func ParseImportStatement(source *[]rune, cursor int) (ImportStatement, int) {
	cursor = moveCursorToImportKeyword(source, cursor)
	startIndex := cursor
	if cursor == len(*source)-1 {
		return ImportStatement{}, cursor
	}

	cursor += len("import")
	cursor = skipWhitespace(source, cursor)

	if isKeyword(source, cursor, "type") {
		cursor = moveCursorPastTypeImport(source, cursor)
		for i := cursor; i < len(*source); i++ {
			// at this point there will be no imports
			// so the code after this loop will never be called
			// in case no inner conditions are met
			if i+1 >= len(*source) {
				return ImportStatement{}, cursor
			}
			char := (*source)[i]
			if char == ',' {
				cursor = i
				break
			}
			if isValidVariableFirstChar(char) {
				if !isKeyword(source, i, "from") {
					return ImportStatement{}, cursor
				}
				_, pathEndIndex := findPath(source, cursor+len("from"))
				if pathEndIndex >= cursor {
					cursor = min(pathEndIndex+1, len(*source))
				}
				return ImportStatement{}, cursor
			}
		}
	}

	var importStatement ImportStatement
	importStatement.StartIndex = startIndex
	// default import
	if isValidVariableFirstChar((*source)[cursor]) || isKeyword(source, cursor, "*") {

		if isKeyword(source, cursor, "*") {
			cursor += len("*")
			cursor = skipWhitespace(source, cursor)
		}
		if isKeyword(source, cursor, "as") {
			cursor += len("as")
			cursor = skipWhitespace(source, cursor)
		}
		endOfDefaultImportName := moveCursorPastVariableName(source, cursor)

		importStatement.DefaultImport = string((*source)[cursor:endOfDefaultImportName])
		if importStatement.DefaultImport == "" {
			return ImportStatement{}, cursor
		}

		for i := cursor; i < len(*source); i++ {
			if i+1 >= len(*source) {
				return ImportStatement{}, cursor
			}
			char := (*source)[i]
			if char == ',' {
				// named imports start
				cursor = min(i+1, len(*source))
				break
			}
			if isValidVariableFirstChar(char) {
				if isKeyword(source, i, "from") {
					cursor = moveCursorToPath(source, cursor)
					start, end := findPath(source, cursor)
					if start > 0 {
						path := string((*source)[start : end+1])
						importStatement.SourcePath = path
						importStatement.EndIndex = end
						return importStatement, end + 1
					}
					return ImportStatement{}, cursor
				}
			}
		}
	}

	if cursor >= len(*source) {
		return ImportStatement{}, len(*source) - 1
	}

	cursor = moveCursorToStartOfNamedImports(source, cursor)
	namedImports := make([]NamedImport, 0, 10)
	cursor = skipWhitespace(source, cursor)

	for i := cursor; i < len(*source); i++ {
		char := (*source)[i]
		if isValidVariableFirstChar(char) {
			if isKeyword(source, i, "type") {
				i = moveCursorPastTypeImport(source, i)
				continue
			}
			if isKeyword(source, i, "from") {
				cursor = i + len("from")
				break
			}
			namedImport, nextIndex := readNamedImport(source, i)
			i = nextIndex

			if namedImport.Name == "" {
				continue
			}
			namedImports = append(namedImports, namedImport)
		}
		if char == '}' {
			cursor = moveCursorToPath(source, i+1)
			break
		}
		if i+1 >= len(*source) {
			return ImportStatement{}, i
		}
	}

	importStatement.NamedImports = namedImports
	if importStatement.DefaultImport == "" && len(importStatement.NamedImports) == 0 {
		return ImportStatement{}, cursor
	}
	cursor = skipWhitespace(source, cursor)
	pathStart, pathEnd := findPath(source, cursor)
	if pathStart < 0 {
		return ImportStatement{}, cursor
	}
	path := string((*source)[pathStart : pathEnd+1])
	importStatement.SourcePath = path
	importStatement.EndIndex = pathEnd
	return importStatement, pathEnd + 1
}

func findPath(source *[]rune, offset int) (int, int) {
	pathStartIndex := skipWhitespace(source, offset)

	if pathStartIndex < 0 {
		return -1, -1
	}
	openingChar := (*source)[pathStartIndex]
	if openingChar != '\'' && openingChar != '"' {
		return -1, -1
	}

	pathEndIndex := skipToWhitespace(source, pathStartIndex) - 1
	if pathEndIndex < 0 {
		return -1, -1
	}
	closingChar := (*source)[pathEndIndex]

	if openingChar != closingChar {
		return -1, -1
	}

	start := pathStartIndex + 1
	end := pathEndIndex - 1

	if start >= end {
		return -1, -1
	}

	return start, end
}

func skipWhitespace(source *[]rune, cursor int) int {
	if cursor >= len(*source) {
		return cursor
	}
	if !unicode.IsSpace((*source)[cursor]) {
		return cursor
	}
	for i, char := range (*source)[cursor:] {
		if !unicode.IsSpace(char) {
			return i + cursor
		}
	}

	return len(*source) - 1
}

func skipToWhitespace(source *[]rune, offset int) int {
	if offset >= len(*source) {
		return -1
	}
	if unicode.IsSpace((*source)[offset]) {
		return offset
	}
	for i, char := range (*source)[offset:] {
		if unicode.IsSpace(char) {
			return i + offset
		}
	}
	return -1
}

func moveCursorPastTypeImport(source *[]rune, cursor int) int {
	cursor = cursor + len("type")
	for i := cursor; i < len(*source); i++ {
		char := (*source)[i]
		if isValidVariableFirstChar(char) {
			i = moveCursorPastVariableName(source, i)
			i = skipWhitespace(source, i)
			if !isKeyword(source, i, "as") {
				return i
			}
			cursor = moveCursorPastVariableName(source, i)
			cursor = skipWhitespace(source, cursor)
			return moveCursorPastVariableName(source, cursor)
		}
	}

	return len(*source) - 1
}

func moveCursorPastVariableName(source *[]rune, cursor int) int {
	for i, char := range (*source)[cursor:] {
		if !isValidVariableNonFirstChar(char) {
			return cursor + i
		}
	}
	return len(*source) - 1
}

func moveCursorToVariable(source *[]rune, cursor int) int {
	for i, char := range (*source)[cursor:] {
		if isValidVariableFirstChar(char) {
			return i + cursor
		}
	}
	return len(*source)
}

func moveCursorToImportKeyword(source *[]rune, cursor int) int {
	for i := cursor; i < len(*source); i++ {
		if isKeyword(source, i, "import") {
			return i
		}
	}
	return len(*source) - 1
}

func moveCursorToStartOfNamedImports(source *[]rune, cursor int) int {
	for i, char := range (*source)[cursor:] {
		if char == '{' {
			return min(i+cursor+1, len(*source))
		}
		if char == 'f' {
			if isKeyword(source, i, "from") {
				return len(*source) - 1
			}
		}
	}

	return len(*source) - 1
}

func readNamedImport(source *[]rune, cursor int) (NamedImport, int) {
	endOfNameIndex := moveCursorPastVariableName(source, cursor) - 1
	name := string((*source)[cursor : endOfNameIndex+1])
	cursor = skipWhitespace(source, endOfNameIndex+1)
	if !isKeyword(source, cursor, "as") {
		return NamedImport{
			Name: name,
		}, cursor
	}

	aliasStartIndex := moveCursorToVariable(source, cursor+len("as"))
	aliasEndIndex := moveCursorPastVariableName(source, aliasStartIndex) - 1
	alias := string((*source)[aliasStartIndex : aliasEndIndex+1])
	return NamedImport{
		Alias: alias,
		Name:  name,
	}, aliasEndIndex + 1

}

func moveCursorToPath(source *[]rune, cursor int) int {
	for i := cursor; i < len(*source); i++ {
		char := (*source)[i]
		if char == 'f' {
			if isKeyword(source, i, "from") {
				i = skipWhitespace(source, i+len("from"))
				return i
			}
		}
	}
	return len(*source) - 1
}
