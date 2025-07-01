package scaner

import (
	"compiler/internal/constants"
	"compiler/internal/parser/js-parser"
	"compiler/internal/types"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const CREATE_ROUTE_SOURCE_RELATIVE_PATH = "./bun-framework/createRoute.ts"
const CREATE_ROUTE_SOURCE_RELATIVE_BARREL_PATH = "./bun-framework/index.ts"

func GetRoute(projectInfo ProjectInfo, filePath string) (types.RouteInfo, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return types.RouteInfo{}, err
	}
	fileString := string(fileBytes)
	fileRunes := []rune(fileString)
	imports := jsparser.ReadImports(&fileRunes, 0)

	createRouteSourcePaths := []string{
		filepath.Join(projectInfo.BasePath, CREATE_ROUTE_SOURCE_RELATIVE_PATH),
		filepath.Join(projectInfo.BasePath, CREATE_ROUTE_SOURCE_RELATIVE_BARREL_PATH),
	}

	var exports jsparser.Exports
	var path string

	for _, importStatement := range imports {
		for _, namedImport := range importStatement.NamedImports {
			if namedImport.Name != "createRoute" {
				continue
			}
			resolvedImportPath := jsparser.ResolveImportPath(jsparser.ResolveImportPathProps{
				FilePath:      filePath,
				BasePath:      projectInfo.BasePath,
				ImportInfo:    importStatement,
				ImportAliases: projectInfo.ImportAliases,
				PackagesList:  projectInfo.PackagesList,
			})
			if resolvedImportPath == "" {
				continue
			}
			if !slices.Contains(createRouteSourcePaths, resolvedImportPath) {
				continue
			}

			var createRouteFnAlias string
			if namedImport.Alias == "" {
				createRouteFnAlias = "createRoute"
				break
			}
			createRouteFnAlias = namedImport.Alias
			path = jsparser.ReadCRPath(&fileRunes, createRouteFnAlias)
			exports = jsparser.ParseExports(filePath, fileString)
			break
		}
	}

	if exports.Default == false && len(exports.Named) == 0 {
		return types.RouteInfo{}, nil
	}

	if exports.Default {
		return types.RouteInfo{
			Path:     path,
			Default:  true,
			Methods:  []string{},
			FilePath: filePath,
		}, nil
	}

	methodsList := strings.Split(constants.METHODS, ",")
	methods := make([]string, 0, len(exports.Named))
	for _, export := range exports.Named {
		if !slices.Contains(methodsList, export) {
			continue
		}
		methods = append(methods, export)
	}

	if len(methods) == 0 {
		return types.RouteInfo{}, nil
	}

	return types.RouteInfo{
		Path:     path,
		Default:  false,
		Methods:  methods,
		FilePath: filePath,
	}, nil
}
