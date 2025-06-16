package main

import (
	// jsParser "compiler/internal/parser/js-parser"
	"compiler/internal/parser/tsconfig"
	"fmt"
	"os"
)

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
	}

	fmt.Println(resolvedAliases)
}
