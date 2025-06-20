package codegen

import (
	"fmt"
	"path/filepath"
	"strings"
)

const NAMED_EXPORT_FUNCTION_TEMPLATE = `
export async function %s(request: BunRequest<%s>): Promise<ResponseType<typeof %s>> {
	let context: Context
	try {
		context = await createContext(request)
	} catch (error) {
		if (error instanceof Response) {
			return error
		}
		logger.error("Error creating context for method ", "%s", " for route: ", %s, "\n", error)
		return new Response("Internal server error", {
			status: 500,
		}) as never
	}
	try {
		const response = await %s(context)
		return context.cleanup(response) as never
	} catch (error) {
		if (error instanceof Response) {
			return context.cleanup(error) as never
		}
		logger.error("Error handling method ", "%s", " for route: ", %s, error)
		return new Response("Internal server error", {
			status: 500,
		}) as never
	}
}
`

const SHARED_IMPORTS = `

`

func routeWithNamedExports(route RouteInfo) string {
	filename := filepath.Base(route.FilePath)
	extension := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, extension)

	content := fmt.Sprintf(`import type { BunRequest } from "bun"
import { createContext, logger, type ResponseType, type Context } from "bun-framework"
import { %s } from "./%s"`, strings.Join(route.Methods, ", "), filename)

	for _, method := range route.Methods {
		content += fmt.Sprintf(NAMED_EXPORT_FUNCTION_TEMPLATE, method, route.Path, method, method, route.Path, method, method, route.Path)
	}

	return content
}
