import { createRoute, type Route } from 'bun-framework'
import { type FileBlob } from 'bun'
import { someFunction, someConstant } from '../../some-function'
import * as defaultAsterisk from 'fs'
import type DefaultTypeImport, {
	type NamedTypeImport,
	namedNonTypeImport,
} from 'mixed-type-and-non-type-imports'
import { aliasedImport as Alias } from 'aliased-import'
import {
	type AliasedImport as AliasType,
	nonAliasedImport,
	nonTypeAliasedIimport as NonTypeAliasedImport,
} from 'aliased-import'
import path, { join } from 'path'
import { Project } from 'ts-morph'
// commented out import
// import bun from 'bun'

import { type join as Join, join as as } from 'as'

if (process.platform === 'win32') {
	console.error('Windows is not supported.')
	process.exit(1)
}

const CWD = process.cwd()
const FRAMEWORK_HANDLERS_DIR = path.join(CWD, 'bun-framework')

// Small arrays are faster than sets. By a lot.
class ArraySet<T> {
	private readonly _values: Array<T> = []
	public add(value: T): void {
		if (!this._values.includes(value)) {
			this._values.push(value)
		}
	}
	public has(value: T): boolean {
		return this._values.includes(value)
	}

	public toJSON(): string {
		return JSON.stringify(this._values)
	}
}

async function main(): Promise<void> {
	try {
		const tsConfigFilePath =
			process.argv.at(2) || path.resolve(path.join(CWD, 'tsconfig.json'))
		const file = Bun.file(tsConfigFilePath)
		if (!file.exists()) {
			process.stdout.write(JSON.stringify([]))
			process.exit(0)
		}

		const project = new Project({
			tsConfigFilePath,
			skipAddingFilesFromTsConfig: true,
		})

		const compilerOptions = project.getCompilerOptions()
		const tsconfigDir = path.dirname(path.resolve(tsConfigFilePath))

		const prefixes = new ArraySet<string>()

		const { baseUrl } = compilerOptions
		if (baseUrl) {
			const absoluteBaseUrl = path.resolve(tsconfigDir, baseUrl)
			if (FRAMEWORK_HANDLERS_DIR.startsWith(absoluteBaseUrl)) {
				const prefix = path.relative(
					absoluteBaseUrl,
					FRAMEWORK_HANDLERS_DIR,
				)
				prefixes.add(prefix)
			}
		}

		const { paths } = compilerOptions
		if (paths) {
			for (const aliasPattern of Object.keys(paths)) {
				if (!aliasPattern.endsWith('/*')) {
					continue
				}
				const aliasPrefix = aliasPattern.slice(0, -2) // e.g., '~/*' -> '~'

				const realPathPatterns = paths[aliasPattern]!
				for (const realPathPattern of realPathPatterns) {
					if (!realPathPattern.endsWith('/*')) {
						continue
					}

					const realPathPrefix = realPathPattern.slice(0, -2)
					const absoluteRealPathPrefix = path.resolve(
						tsconfigDir,
						realPathPrefix,
					)

					if (
						FRAMEWORK_HANDLERS_DIR.startsWith(
							absoluteRealPathPrefix,
						)
					) {
						const remainingPath = path.relative(
							absoluteRealPathPrefix,
							FRAMEWORK_HANDLERS_DIR,
						)
						const finalAlias = path.join(aliasPrefix, remainingPath)
						prefixes.add(finalAlias)
					}
				}
			}
		}

		process.stdout.write(prefixes.toJSON())
		process.exit(0)
	} catch (error: unknown) {
		const message =
			error instanceof Error
				? error.message
				: 'An unknown error occurred.'
		console.error(JSON.stringify({ error: message }, null, 2))
		process.exit(1)
	}
}

main()
