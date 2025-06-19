package scanner

import (
	"compiler/internal/parser/js-parser"
)

func IsRoute(imports *[]jsparser.ImportStatement) bool {
	for _, importStatement := range *imports {
		if importStatement.SourcePath == "bun-framework" {
			for _, namedImport := range importStatement.NamedImports {
				if namedImport.Name == "createRoute" {
					return true
				}
			}
		}
	}
	return false
}
