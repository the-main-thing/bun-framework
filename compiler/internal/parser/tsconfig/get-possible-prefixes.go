package tsconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


func GetPossiblePrefixes() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}
	tsConfigFileBytes, err := os.ReadFile(filepath.Join(cwd, "tsconfig.json"))
	if err != nil {
		fmt.Println("Error reading tsconfig.json:", err)
		os.Exit(1)
	}

	if tsConfigFileBytes == nil || len(tsConfigFileBytes) == 0 {
		fmt.Println("Error: tsconfig.json is empty")
		os.Exit(1)
	}

	baseUrl, err := readBasePath(&tsConfigFileBytes)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error reading baseUrl from tsconfig.json:", err)
	}
	if baseUrl == "" {
		baseUrl = "."
	}
	paths, err := readPaths(&tsConfigFileBytes)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error reading paths from tsconfig.json:", err)
		paths = []string{}
	}

	absoluteBaseUrl := filepath.Join(cwd, baseUrl)
	for _, path := range paths {
	  if strings.HasSuffix(path, "*")	{

		}
	}


}
