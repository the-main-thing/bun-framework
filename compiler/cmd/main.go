package main

import (
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
		os.Exit(1)
	}

	if resolvedAliases == nil {
		fmt.Println("No aliases")
		os.Exit(1)
	}
}
