package tsconfig

import (
	jsParser "compiler/internal/parser/js-parser"
	"errors"
	"strings"
)

// get baseUrl from tsconfig.json
func readBasePath(tsConfig *[]byte) (string, error) {
	baseUrlIndex := strings.Index(string(*tsConfig), "baseUrl")
	if baseUrlIndex == -1 {
		return "", errors.New("baseUrl not found in tsconfig.json")
	}

	var baseUrl string
	_, keyEnd, err := jsParser.ReadObjectKey(tsConfig, baseUrlIndex)
	if err != nil {
		return "", err
	}
	valueStart, valueEnd, err := jsParser.ReadString(tsConfig, keyEnd+1)
	if err != nil {
		return "", err
	}
	baseUrl = string((*tsConfig)[valueStart:valueEnd])

	return baseUrl, nil
}
