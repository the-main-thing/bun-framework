import type { BunRequest } from 'bun'
import { getContext } from '~/server/context'
import handler, { pathName } from '__FILENAME__'

export type ResponseType = Awaited<ReturnType<typeof handler>>

export default async function __NAME__(
	request: BunRequest<typeof pathName>,
): Promise<ResponseType> {
	try {
		const context = await getContext(request)
		const response = await handler({ request, context }).catch((error) => {
			if (error instanceof Response) {
				return error as never
			}
			throw error
		})
		if (context.shouldDoWorkBeforeResponse) {
			return await context.doWorkBeforeResponse(response) as never
		}
		return response as never
	} catch (error) {
		console.error(
			'\n\n',
			new Date().toISOString(),
			`\nError handling request for ${pathName}\n`,
			error,
			'\n',
		)
		return new Response(error.message, { status: 500 }) as never
	}
}
