package jsparser

import (
	"compiler/internal/parser/types"
	"os"
	"path/filepath"
	"strings"
)

type ImportResolverInfo struct {
	ResolvedAliases []types.ResolvedImportAlias
	// Sould be absolute
	FilePath   string
	ImportPath string
}

// Uses syscall so concider putting it inside of a go routine
func ResolveImport(importResolverInfo ImportResolverInfo) string {
	// relative
	if strings.HasPrefix(importResolverInfo.ImportPath, ".") {
		importPath := normalizePath(importResolverInfo.ImportPath)
		absoluteImportPath, err := filepath.Abs(filepath.Join(importResolverInfo.FilePath, importPath))
		if err != nil {
			return ""
		}
		_, err = os.Stat(absoluteImportPath)
		if err != nil {
			return ""
		}
		return absoluteImportPath
	}
	// aliased
	for _, alias := range importResolverInfo.ResolvedAliases {
		if !strings.HasPrefix(importResolverInfo.ImportPath, alias.Alias) {
			continue
		}
		importPath := strings.TrimPrefix(importResolverInfo.ImportPath, alias.Alias)
		normalizedImportPath := normalizePath(importPath)
		for _, path := range alias.Paths {
			absolutePath, err := filepath.Abs(filepath.Join(path, normalizedImportPath))
			if err != nil {
				continue
			}
			_, err = os.Stat(absolutePath)
			if err != nil {
				continue
			}
			return absolutePath
		}
	}
	// absolute
	importPath := normalizePath(importResolverInfo.ImportPath)
	absolutePath, err := filepath.Abs(importPath)
	if err != nil {
		return ""
	}
	_, err = os.Stat(absolutePath)
	if err != nil {
		return ""
	}
	return absolutePath
}

func normalizePath(path string) string {
	segments := strings.Split(path, "/")
	return strings.Join(segments, string(filepath.Separator))
}
