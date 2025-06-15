export async function start() {
	const server = Bun.serve({
		routes: {},
	})

	return server
}
