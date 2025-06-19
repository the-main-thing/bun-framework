package jsparser

// Comments should be removed before calling this function
func ReadImports(source *[]rune, offset int) []ImportStatement {
	var imports []ImportStatement
	for i := offset; i < len(*source); i++ {
		importStatement, cursor := ParseImportStatement(source, i)
		i = cursor
		if importStatement.SourcePath == "" {
			continue
		}
		imports = append(imports, importStatement)
	}

	return imports
}
