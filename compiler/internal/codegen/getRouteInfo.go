package codegen

import (
	"compiler/internal/parser/js-parser"
	"os"
	"slices"
	"strings"
)

type RouteInfo struct {
	Path     string
	Default  bool
	Methods  []string
	FilePath string
}

func GetRouteInfo(filePath string) RouteInfo {
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	source := []rune(string(file))
	jsparser.RemoveComments(&source)
	imports := jsparser.ReadImports(&source, 0)
	createRouteFnAlias := ""
	for _, importStatement := range imports {
		for i := importStatement.StartIndex; i < importStatement.EndIndex; i++ {
			source[i] = ' '
		}
		if strings.Contains(importStatement.SourcePath, "bun-framework") {
			continue
		}
		for _, namedImport := range importStatement.NamedImports {
			if namedImport.Name == "createRoute" && createRouteFnAlias == "" {
				if namedImport.Alias == "" {
					createRouteFnAlias = "createRoute"
					continue
				}
				createRouteFnAlias = namedImport.Alias
				continue
			}
		}
	}
	if createRouteFnAlias == "" {
		return RouteInfo{}
	}

	path := jsparser.ReadCRPath(&source, createRouteFnAlias)
	if path == "" {
		return RouteInfo{}
	}
	var routeInfo RouteInfo
	routeInfo.Path = path
	routeInfo.FilePath = filePath
	exports := jsparser.ParseExports(filePath, string(source))
	if exports.Default {
		routeInfo.Default = true
		return routeInfo
	}

	methods := strings.Split(METHODS, ",")
	for _, named := range exports.Named {
		if !slices.Contains(methods, named) {
			continue
		}
		routeInfo.Methods = append(routeInfo.Methods, named)
	}
	if len(routeInfo.Methods) == 0 {
		return RouteInfo{}
	}
	return routeInfo
}
