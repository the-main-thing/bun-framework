package tsconfig

import (
	jsParser "compiler/internal/parser/js-parser"
	"compiler/internal/parser/types"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetAliases(absoluteTsConfigFilePath string) ([]types.ResolvedImportAlias, error) {
	tsConfigFileBytes, err := os.ReadFile(absoluteTsConfigFilePath)
	if err != nil {
		return nil, errors.New("Error reading tsconfig file: " + absoluteTsConfigFilePath + "\n" + err.Error())
	}

	if tsConfigFileBytes == nil || len(tsConfigFileBytes) == 0 {
		return nil, errors.New("tsconfig file is empty: " + absoluteTsConfigFilePath)
	}

	runes := []rune(string(tsConfigFileBytes))
	jsParser.RemoveComments(&runes)

	var tsconfig types.TsConfig
	err = json.Unmarshal([]byte(string(runes)), &tsconfig)
	if err != nil {
		return nil, errors.New("Can't unmarshal tsconfig after removing all the comments. Likely the file is corrupted")
	}

	baseUrl := normalizePath(tsconfig.CompilerOptions.BaseUrl)
	baseUrl, err = filepath.Abs(baseUrl)
	if err != nil {
		tsconfigDir := filepath.Dir(absoluteTsConfigFilePath)
		baseUrl, err = filepath.Abs(filepath.Join(tsconfigDir, baseUrl))
		if err != nil {
			return nil, errors.New("Can't resolve baseUrl: " + baseUrl + "\n" + err.Error())
		}
		if !strings.HasPrefix(tsconfigDir, baseUrl) {
			fmt.Fprintln(os.Stderr, "Resolved baseUrl is outside of project's base directory: "+baseUrl+"\n")
		}
	}

	resolvedAliases := make([]types.ResolvedImportAlias, 0, len(tsconfig.CompilerOptions.Paths))
	for key, value := range tsconfig.CompilerOptions.Paths {
		if !strings.HasSuffix(key, "/*") {
			continue
		}
		alias := key[:len(key)-1]
		resolvedPaths := make([]string, 0, len(value))
		for _, path := range value {
			if !strings.HasSuffix(path, "/*") {
				continue
			}
			path = path[:len(path)-1]
			if strings.HasPrefix(path, ".") {
				absolutePath, err := filepath.Abs(filepath.Join(baseUrl, normalizePath(path)))
				if err != nil {
					fmt.Fprintln(os.Stderr, "Can't resolve path: "+path+"\n"+err.Error())
					continue
				}
				resolvedPaths = append(resolvedPaths, absolutePath)
				continue
			}
			absolutePath, err := filepath.Abs(filepath.Join(".", normalizePath(path)))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Can't resolve path: "+path+"\n"+err.Error())
				continue
			}
			resolvedPaths = append(resolvedPaths, absolutePath)
		}
		if len(resolvedPaths) == 0 {
			continue
		}
		resolvedAliases = append(resolvedAliases, types.ResolvedImportAlias{
			Alias: alias,
			Paths: resolvedPaths,
		})
	}
	return resolvedAliases, nil
}

func normalizePath(path string) string {
	segments := strings.Split(path, "/")
	return strings.Join(segments, string(filepath.Separator))
}
