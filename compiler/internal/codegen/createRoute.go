package codegen

import (
	"fmt"
	"path/filepath"
)

func CreateRoute(routeInfo RouteInfo) error {
	if routeInfo.Path == "" {
		return nil
	}
	filename := filepath.Base(routeInfo.FilePath)
	extension := filepath.Ext(routeInfo.FilePath)
	compiledPath := filepath.Join(filepath.Dir(routeInfo.FilePath), filename+"-compiled"+extension)

	var result string
	if routeInfo.Default {
		result = fmt.Sprintf("import type { BunRequest } from \"bun\"\nimport routeHandler from \"./%s\"\nexport default async function(request: BunRequest<\"%s\">) {\nreturn routeHandler({request})\n}", filename, routeInfo.Path)
	}
	
}
