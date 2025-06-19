import { createRoute } from 'bun-framework'

const DEFAULT = createRoute('/api/v1/test', async (request) => {
	return new Response('Hello World!')
})

export const GET = createRoute('/api/v1/test', async (request) => {
	return new Response('Hello World!')
})

export const POST = createRoute('/api/v1/test', async (request) => {
	return new Response('Hello World!')
})

const PUT = createRoute('/api/v1/test', async (request) => {
	return new Response('Hello World!')
})

type Put = typeof PUT

export { PUT, type Put as PutHandler }

const d = createRoute('/api/v1/test', async (request) => {
	return new Response('Hello World!')
})

export { d as DELETE }

export default DEFAULT
