package jsparser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGettingStringValueWithQuotes(t *testing.T) {
	expectedOutput := []string{
		`"first"`,
		`"se\"cond"`,
		"`third`",
		"`fou\\`rth`",
		"'super\\\\\\n\\'escap\\\\\\\\\\\\\\'ed'",
	}

	filePath := filepath.Join(".", "getStringValue_mock.txt")
	file, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	runes := []rune(string(file))
	results := make([]string, 0, len(expectedOutput))
	for i := 0; i < len(runes); i++ {
		value, cursor := GetStringWithQuotes(&runes, i)
		if cursor <= i {
			break
		}
		i = cursor
		results = append(results, value)
	}

	if len(results) != len(expectedOutput) {
		t.Fatalf("Expected %d results, got %d\nexpected: %s\nreceived: %s\n", len(expectedOutput), len(results), strings.Join(expectedOutput, ", "), strings.Join(results, ", "))
	}

	for i, result := range results {
		expected := expectedOutput[i]
		if result != expected {
			t.Fatalf("Expected %s, got %s", expected, result)
		}
	}
}
