# bun-framework

My custom framework for bun

The goal is to learn golang while also inventing a new complier for fun while
also create a nice tool for building my pet apps with bun.

This thing is not even compiles yet.

The rules will be the following:

1. You choose a directory for the source code
2. At the root of the directory there should be a `server.ts|tsx` with routes
   definition in any place, i.e.:

```ts
Bun.serve({
	routes: {}, // this is what framework tools will be looking for
})
```

You can edit this file as you wish, just keep routes definition inside of
`Bun.serve({ routes: {} })` i.e.:

```ts
import { myCustomRouteHandler } from './any-path/my-custom-route-handler'

export async function startServer() {
	// Do whatever you like here

	const server = Bun.serve({
		routes: {
			'/*': new Response('Hello world!'),
			'/hello': myCustomRouteHandler,
		},
	})

	return server
}
```

3. The build tool will create `server-compiled.ts|tsx` file. Make sure to use
   this for starting the server
4. To create a route file you import `createRoute` function, call it and export
   it as default or as `GET|PUT|POST|DELETE|OPTIONS|HEAD`-named export, i.e.:

```ts
import { createRoute } from 'bun-framework'
import { someFunction } from './some-function'

export default createRoute(
	'/api/:id',
	async function myCustomRouteHandler({ request, context }) {
		if (request.method !== 'GET') {
			// Instead of returning, thow it, so ts will ifer types that should be expected by the client
			throw new Response('Method not allowed', 405)
		}
		// Calling the `context.getUser()` will trigger the need to update the cookies or tokens or whatever
		const user = await context.getUser()

		if (!user) {
			throw new Response('Unauthorized', 401)
		}

		const data = someFunction(request.params.id, user)
		// framework will handle updating the cookies or tokens for you as well as generating client types, adding CORS, etc.
		return TypedResponse.json(data)
	},
)
```

or

```ts
import { createRoute } from '~/bun-framework'
import { someFunction } from './some-function'

export const GET = createRote('/api/:id', async (args) => {
	const { request, context } = args
	const user = await context.getUser()

	if (!user) {
		throw new Response('Unauthorized', 401)
	}

	const data = someFunction(request.params.id, user)
	return TypedResponse.json(data)
})
```

5. Add all the relevant handlers for the context (TODO)
6. Run the build tool or start in dev mode: TODO
