package main

import (
	jsParser "compiler/internal/parser/js-parser"
	"compiler/internal/parser/tsconfig"
	"fmt"
	"os"
)

func test(output []jsParser.ImportStatement) {
	fmt.Println("Testing...")
	expectedOutput := []jsParser.ImportStatement{
		{
			SourcePath: "bun-framework",
			NamedImports: []jsParser.NamedImport{
				{
					Name: "createRoute",
				},
			},
		},
		{
			NamedImports: []jsParser.NamedImport{
				{
					Name: "someFunction",
				},
				{
					Name: "someConstant",
				},
			},
			SourcePath: "../../some-function",
		},
		{
			DefaultImport: "defaultAsterisk",
			SourcePath:    "fs",
		},
		{
			NamedImports: []jsParser.NamedImport{
				{
					Name: "namedNonTypeImport",
				},
			},
			SourcePath: "mixed-type-and-non-type-imports",
		},
		{
			NamedImports: []jsParser.NamedImport{
				{
					Alias: "Alias",
					Name:  "aliasedImport",
				},
			},
			SourcePath: "aliased-import",
		},
		{
			NamedImports: []jsParser.NamedImport{
				{
					Name: "nonAliasedImport",
				},
				{
					Alias: "NonTypeAliasedImport",
					Name:  "nonTypeAliasedIimport",
				},
			},
			SourcePath: "aliased-import",
		},
		{
			DefaultImport: "path",
			NamedImports: []jsParser.NamedImport{
				{
					Name: "join",
				},
			},
			SourcePath: "path",
		},
		{
			SourcePath: "ts-morph",
			NamedImports: []jsParser.NamedImport{
				{
					Name: "Project",
				},
			},
		},
		{
			SourcePath: "as",
			NamedImports: []jsParser.NamedImport{
				{
					Name:  "join",
					Alias: "as",
				},
			},
		},
	}
	if len(output) != len(expectedOutput) {
		fmt.Println("Output length is not equal to expected length", len(expectedOutput), len(output))
		return
	}
	for i, statement := range output {
		expected := expectedOutput[i]
		if statement.SourcePath != expected.SourcePath {
			fmt.Println("SourcePath is not equal")
			return
		}
		if len(statement.NamedImports) != len(expected.NamedImports) {
			fmt.Println("NamedImports length is not equal")
			fmt.Println(expected)
			fmt.Println(statement)
			return
		}
		for j, expectedNamedImport := range expected.NamedImports {
			namedImport := statement.NamedImports[j]
			if namedImport.Alias != expectedNamedImport.Alias {
				fmt.Println("NamedImports alias is not equal", expectedNamedImport.Alias, namedImport)
				return
			}
			if namedImport.Name != expectedNamedImport.Name {
				fmt.Println("NamedImports name is not equal")
				return
			}
		}
		if statement.DefaultImport != expected.DefaultImport {
			fmt.Println("DefaultImport is not equal")
			return
		}
	}
	fmt.Println("Test passed")
}

func main() {
	basePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	resolvedAliases, err := tsconfig.GetAliases(tsconfig.GetAliasesProps{
		BasePath:         basePath,
		TsconfigFilePath: "./tsconfig.json",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading tsconfig", err)
		os.Exit(1)
	}

	if resolvedAliases == nil {
		fmt.Println("No aliases")
		os.Exit(1)
	}

	tsFile, err := os.ReadFile("./main.ts")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading main.ts", err)
		os.Exit(1)
	}

	runes := []rune(string(tsFile))

	fmt.Println("Removing comments...")

	jsParser.RemoveComments(&runes)

	fmt.Println("Reading imports...")
	importStatements := jsParser.ReadImports(&runes, 0)

	test(importStatements)

	// fmt.Println("-----------")
	//
	// for _, importStatement := range importStatements {
	// 	fmt.Println("DefaultImport", importStatement.DefaultImport)
	// 	fmt.Println("NamedImports")
	// 	for _, namedImport := range importStatement.NamedImports {
	// 		if namedImport.Alias == "" {
	// 			fmt.Println("\t", namedImport.Name)
	// 			continue
	// 		} else {
	// 			fmt.Print("\t", namedImport.Alias, " as ", namedImport.Name)
	// 		}
	// 	}
	// 	fmt.Println()
	// 	fmt.Println("SourcePath", importStatement.SourcePath)
	// 	fmt.Println("-----------")
	// }

	fmt.Println("Done")
}
