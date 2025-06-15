package tsconfig

import (
	parser "compiler/internal/parser/js-parser"
	"errors"
	"strings"
)

// get list of paths from an input like this:
//
//	{
//	  "compilerOptions": {
//	    "paths": {
//	      "~/*": ["./*"],
//	      // valid tsconfig comment that is invalid JSON
//	      "bun-framework/*": ["./bun-framework/*"]
//	    }
//	  }
//	}
func readPaths(tsConfig *[]byte) ([]string, error) {
	pathsIndex := strings.Index(string(*tsConfig), "paths")

	if pathsIndex == -1 {
		return []string{}, nil
	}

	var result []string
	reading := false
	for i := pathsIndex + len("paths") + 1; i < len(*tsConfig); i++ {
		b := (*tsConfig)[i]
		if b == '{' {
			if reading {
				return nil, errors.New("Nested objects are not supported and actually it seems that tsconfig.json is not valid")
			}
			reading = true
		}
		if b == '}' {
			break
		}
		_, keyEnd, err := parser.ReadObjectKey(tsConfig, i)
		if err != nil {
			return nil, err
		}
		valuesIndexes, err := parser.ReadArrayOfStrings(tsConfig, keyEnd+1)
		if err != nil {
			return nil, err
		}
		for _, valueIndex := range valuesIndexes {
			value := string((*tsConfig)[valueIndex.StartIndex:valueIndex.EndIndex])
			result = append(result, value)
		}
		continue
	}
	return result, nil
}
