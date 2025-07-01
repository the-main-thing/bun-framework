package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	abs, err := (filepath.Abs(filepath.Join(cwd, "/Users/pavelshevtsov/bun-framework/sandbox")))
	if err != nil {
		panic(err)
	}
	fmt.Println(abs)
}
