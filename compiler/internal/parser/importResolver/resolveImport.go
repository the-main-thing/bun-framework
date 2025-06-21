package importResolver

import (
	"compiler/internal/parser/js-parser"
	"compiler/internal/parser/tsconfig"
	"compiler/internal/parser/types"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ResolvedImport struct {
	Default bool
	Name    string
	Alias   string
	Path    string
}

func ResolveImportPath(filePath string, importInfo jsparser.ImportStatement, importAliases []types.ResolvedImportAlias) string {
	for _, alias := range importAliases {
		if strings.HasPrefix(importInfo.SourcePath, alias.Alias) {
			break
		}
	}

	if strings.HasPrefix(importInfo.SourcePath, ".") {
		realPath, err := filepath.Abs(filepath.Join(filepath.Dir(filePath), importInfo.SourcePath))
		if err != nil {
			return realPath
		}
	}
	return filepath.Join("node_modules", importInfo.SourcePath)
}

// I'm not sure how it will work with wierd aliases that involve slashes
func NormalizeImportPath(importPath string) string {
	firstChar := rune(importPath[0])
	lastChar := rune(importPath[len(importPath)-1])
	if firstChar == '\'' || firstChar == '"' {
		if lastChar == firstChar {
			return NormalizeImportPath(importPath[1 : len(importPath)-1])
		}
		return ""
	}
	if firstChar == '.' {
		return strings.ReplaceAll(importPath, ".", string(filepath.Separator))
	}
	
}
