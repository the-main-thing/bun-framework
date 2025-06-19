package jsparser

// comments should be removed before calling this function
func RemoveImports(source *[]rune) {
	cursor := moveCursorToImportKeyword(source, 0)
	imports := ReadImports(source, cursor)
	for _, importStatement := range imports {
		if importStatement.SourcePath == "" {
			continue
		}
		for i := importStatement.StartIndex; i < importStatement.EndIndex; i++ {
			(*source)[i] = ' '
		}
	}
}
