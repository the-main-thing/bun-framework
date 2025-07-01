package jsparser

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestGettingPackagesList(t *testing.T) {
	expectedOutput := []string{
		"@types/bun",
		"typescript",
		"ts-morph",
	}

	results := GetPackagesList(filepath.Join(".", "packageJson_mock.json"))

	if len(results) != len(expectedOutput) {
		t.Fatalf("Expected %d results, got %d\nexpected: %s\nreceived: %s\n", len(expectedOutput), len(results), strings.Join(expectedOutput, ", "), strings.Join(results, ", "))
	}

	for _, result := range results {
		if slices.Contains(expectedOutput, result) {
			continue
		}
		t.Fatal("Got unexpected value", result)
	}
}
