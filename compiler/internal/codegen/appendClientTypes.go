package codegen

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"compiler/internal/types"
)

const DEFAULT_IMPORT_TEMPLATE = "import type %s from '%s'"
const NAMED_IMPORT_TEMPALTE = "import type { %s as %s } from '%s'"

type PathType struct {
	Path string
	Type string
}

type DefaultImport struct {
	FunctionName string
	TypeName string
	ImportStatement string
}

type Artifact struct {
	File string
	Path string
	DefaultImport DefaultImport
	NamedImports string
	Types string
	Router string
	Next *Artifact
}

type ClientTypes struct {
	DefaultImports []string
	NamedImports   []string
	RoutesTypes    []PathType
	Types          []string
}

func pathToJsImport(path string) string {
	return strings.ReplaceAll(path, string(filepath.Separator), "/")
}

func appendClientTypes(artifacts *Artifact, routeInfo types.RouteInfo) {
	var next Artifact

	if routeInfo.Default {
		functionName := strings.ReplaceAll(routeInfo.FilePath, string(filepath.Separator), "_")
		firstLetter := string(unicode.ToUpper(rune(functionName[0])))
		typeName := firstLetter + functionName[1:]

		defaultImport := fmt.Sprintf(DEFAULT_IMPORT_TEMPLATE, functionName, pathToJsImport(routeInfo.Path))
		typeDef := fmt.Sprintf("type %s = Awaited<ReturnType<typeof %s>>['type']", typeName, functionName)
		next.DefaultImport = DefaultImport{
			TypeName: typeName,
			FunctionName: functionName,
			ImportStatement: defaultImport,
		}

		return
	}

	methodsHandlers := make([]string, len(routeInfo.Methods))
	namedImports := make([]string, len(routeInfo.Methods))
	typesToAppend := make([]string, len(routeInfo.Methods))

	for i, method := range routeInfo.Methods {
		handlerName := fmt.Sprintf("%s_handler", method)
		methodsHandlers[i] = handlerName
		namedImports[i] = fmt.Sprintf("%s as %s", method, handlerName)
		typesToAppend[i] = fmt.Sprintf("type %s = Awaited<ReturnType<typeof %s>>['type']", method, handlerName)
	}

	types.NamedImports = append(types.NamedImports, fmt.Sprintf(NAMED_IMPORT_TEMPALTE, strings.Join(namedImports, ", "), pathToJsImport(routeInfo.Path)))
	types.Types = append(types.Types, typesToAppend...)
	types.RoutesTypes = append(types.RoutesTypes, PathType{
		Path: fmt.Sprintf("'%s'", routeInfo.Path),
		Type: strings.Join(routeInfo.Methods, ": "),
	})

	return
}

func serializeClientTypes(types *ClientTypes) string {
	defaultImports := strings.Join(types.DefaultImports, "\n")
	namedImports := strings.Join(types.NamedImports, "\n")
	handlersReturnTypes := strings.Join(types.Types, "\n")

	result := strings.Join([]string{defaultImports, namedImports, handlersReturnTypes, "export type Routes = {"}, "\n")

	for _, route := range types.RoutesTypes {
		result += fmt.Sprintf("\t%s: %s,", route.Path, route.Type)
	}
	result += "}\n"
	return result
}


