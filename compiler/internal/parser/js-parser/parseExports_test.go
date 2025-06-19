package jsparser

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestAllExports(t *testing.T) {
	expectedOutput := Exports{
		Default: true,
		Named: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
	}

	filePath := filepath.Join(".", "parseExports_mock.ts")
	file, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	exports := ParseExports(filePath, string(file))

	if len(exports.Named) != len(expectedOutput.Named) {
		t.Fatalf("Expected %d named exports, got %d", len(expectedOutput.Named), len(exports.Named))
	}
	if exports.Default != expectedOutput.Default {
		t.Fatalf("Expected %t default export, got %t", expectedOutput.Default, exports.Default)
	}
	for _, named := range exports.Named {
		if !slices.Contains(expectedOutput.Named, named) {
			t.Fatalf("Unexpected named export, got %s", named)
			break
		}
	}
	for _, named := range expectedOutput.Named {
		if !slices.Contains(exports.Named, named) {
			t.Fatalf("Missing named export %s", named)
			break
		}
	}
}
