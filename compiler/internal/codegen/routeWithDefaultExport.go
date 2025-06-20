package codegen

import (
	"fmt"
	"path/filepath"
	"strings"
)

const TEMPLATE = `import type { BunRequest } from "bun"
import { createContext, logger, type ResponseType, type Context } from "bun-framework"

import routeHandler from "./%s"

export default async function %s(request: BunRequest<%s>): Promise<ResponseType<typeof routeHandler>> {
	let context: Context
	try {
		context = await createContext(request)
	} catch (error) {
		if (error instanceof Response) {
			return error
		}
		logger.error("Error creating context for route: ", %s, "\n", error)
		return new Response("Internal server error", {
			status: 500,
		}) as never
	}
	try {
		const response = await routeHandler(context)
		return context.cleanup(response) as never
	} catch (error) {
		if (error instanceof Response) {
			return context.cleanup(error) as never
		}
		logger.error("Error handling route: ", %s, error)
		return new Response("Internal server error", {
			status: 500,
		}) as never
	}
}`

func routeWithDefaultExport(route RouteInfo) string {
	filename := filepath.Base(route.FilePath)
	extension := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, extension)
	functionNameSource := []rune(route.Path)
	for i, char := range functionNameSource {
		if char >= 'a' && char <= 'z' {
			continue
		}
		if char >= 'A' && char <= 'Z' {
			continue
		}
		functionNameSource[i] = '_'
	}
	functionName := string(functionNameSource[1:])

	return fmt.Sprintf(TEMPLATE, filename, functionName, route.Path, route.Path)
}
