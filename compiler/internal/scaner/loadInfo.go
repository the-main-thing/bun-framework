package scaner

import (
	"compiler/internal/parser/js-parser"
	"compiler/internal/parser/tsconfig"
	"compiler/internal/parser/types"
	"path/filepath"
)

type ProjectInfo struct {
	ImportAliases []types.ResolvedImportAlias
	PackagesList  []string
	BasePath      string
}

func LoadInfo(tsConfigFilePath string) (ProjectInfo, error) {
	absoluteTsConfigFilePath, err := filepath.Abs(tsConfigFilePath)
	if err != nil {
		return ProjectInfo{}, err
	}
	aliases, err := tsconfig.GetAliases(absoluteTsConfigFilePath)
	if err != nil {
		return ProjectInfo{}, err
	}
	basePath := filepath.Dir(absoluteTsConfigFilePath)
	return ProjectInfo{
		ImportAliases: aliases,
		BasePath:      basePath,
		PackagesList:  jsparser.GetPackagesList(filepath.Join(basePath, "package.json")),
	}, nil
}
