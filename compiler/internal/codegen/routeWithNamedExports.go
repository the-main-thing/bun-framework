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

func routeWithNamedExports(filePath string, routePath string, methods []string) string {
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, extension)

	content := fmt.Sprintf(`import type { BunRequest } from "bun"
import { createContext, logger, type ResponseType, type Context } from "bun-framework"
import { %s } from "./%s"`, strings.Join(methods, ", "), filename)

	for _, method := range methods {
		content += fmt.Sprintf(NAMED_EXPORT_FUNCTION_TEMPLATE, method, routePath, method, method, routePath, method, method, routePath)
	}

	return content
}
