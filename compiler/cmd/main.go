package main

import (
	"compiler/internal/parser/tsconfig"
	"fmt"
	"os"
	"path/filepath"
)

const DEFAULT_TS_CONFIG_FILE_PATH = "./tsconfig.json"

func main() {
	var tsConfigPath string
	if len(os.Args) >= 3 && os.Args[1] == "--tsconfig" && os.Args[2] != "" {
		tsConfigPath = os.Args[2]
	}
	if tsConfigPath == "" {
		tsConfigPath = DEFAULT_TS_CONFIG_FILE_PATH
	}

	basePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !filepath.IsAbs(tsConfigPath) {
		tsConfigPath, err = filepath.Abs(filepath.Join(basePath, tsConfigPath))
		if err != nil {
			panic(err)
		}
	}
	resolvedAliases, err := tsconfig.GetAliases(tsConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading tsconfig", err)
		os.Exit(1)
	}

	if resolvedAliases == nil {
		fmt.Println("No aliases")
		os.Exit(1)
	}
}

