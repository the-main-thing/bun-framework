package jsparser

import (
	"compiler/internal/parser/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type ResolvedImport struct {
	Default bool
	Name    string
	Alias   string
	Path    string
}

type ResolveImportPathProps struct {
	FilePath      string
	BasePath      string
	ImportInfo    ImportStatement
	ImportAliases []types.ResolvedImportAlias
	PackagesList  []string
}

func ResolveImportPath(props ResolveImportPathProps) string {
	importInfo := props.ImportInfo
	if strings.HasPrefix(importInfo.SourcePath, "./") || strings.HasPrefix(importInfo.SourcePath, "../") {
		normalizedPath := NormalizeImportPath(importInfo.SourcePath)
		absolutePath, err := filepath.Abs(filepath.Join(filepath.Dir(props.FilePath), normalizedPath))
		if err != nil {
			return ""
		}
		return absolutePath
	}

	for _, alias := range props.ImportAliases {
		if !strings.HasPrefix(importInfo.SourcePath, alias.Alias) {
			continue
		}
		pathWhithoutAlias := importInfo.SourcePath[len(alias.Alias):]
		normalizedPath := NormalizeImportPath(pathWhithoutAlias)

		wg := sync.WaitGroup{}
		var resolvedPath atomic.Value
		for _, path := range alias.Paths {
			var absolutePath string
			var err error
			if !filepath.IsAbs(normalizedPath) {
				absolutePath, err = filepath.Abs(filepath.Join(path, normalizedPath))
				if err != nil {
					continue
				}
			}
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				value := resolvedPath.Load()
				if value != nil {
					return
				}
				result := ResolveToFileIfExists(absolutePath)
				value = resolvedPath.Load()
				if value != nil || result == "" {
					return
				}
				resolvedPath.Store(result)
			}(path)
		}
		wg.Wait()
		value := resolvedPath.Load()
		if value == nil {
			// Cannot resolve the path for alias, so try to resolve via node_modules
			break
		}
		return value.(string)
	}

	normalizedPath := NormalizeImportPath(importInfo.SourcePath)
	for _, packageName := range props.PackagesList {
		if strings.HasPrefix(importInfo.SourcePath, packageName) {
			absoultePath, err := filepath.Abs(filepath.Join(props.BasePath, "node_modules", normalizedPath))
			if err != nil {
				return ""
			}
			return absoultePath
		}
	}

	if filepath.IsAbs(normalizedPath) {
		return normalizedPath
	}
	absoultePath, err := filepath.Abs(filepath.Join(props.BasePath, normalizedPath))
	if err != nil {
		return ""
	}
	return absoultePath
}

// I'm not sure how it will work with wierd aliases that involve slashes
func NormalizeImportPath(importPath string) string {
	firstChar := rune(importPath[0])
	if firstChar == '\'' || firstChar == '"' {
		lastChar := rune(importPath[len(importPath)-1])
		if lastChar == firstChar {
			return NormalizeImportPath(importPath[1 : len(importPath)-1])
		}
		return ""
	}
	if filepath.Separator == '/' {
		return importPath
	}
	return strings.ReplaceAll(importPath, "/", string(filepath.Separator))
}

func ResolveToFileIfExists(filePath string) string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fileInfo, err = os.Stat(filepath.Dir(filePath))
		if err != nil {
			return ""
		}
	}
	if !fileInfo.IsDir() {
		return filePath
	}
	filesList, err := os.ReadDir(filePath)
	if err != nil {
		return ""
	}

	fileName := filepath.Base(filePath)
	for _, file := range filesList {
		if strings.HasPrefix(file.Name(), fileName) {
			return filepath.Join(filePath, file.Name())
		}
	}

	for _, file := range filesList {
		if strings.HasPrefix(file.Name(), "index.") {
			return filepath.Join(filePath, file.Name())
		}
	}

	return ""
}
