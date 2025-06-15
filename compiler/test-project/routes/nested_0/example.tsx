'@framework route'

import { someFunc } from '~/server/someFunc'
import type { AuthRouteArgs } from '../routes/'
import { json } from '../json'

export const pathName = `/public/*`

export default async function example({
	request,
	context,
}: AuthRouteArgs<typeof pathName>) {
	const session = await context.getSession()
	if (!session) {
		throw new Response('Unauthorized', { status: 401 })
	}
	const someThing = await someFunc(request)
	return json({ someThing })
}
