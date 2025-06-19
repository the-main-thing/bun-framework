package jsparser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAllImports(t *testing.T) {
	expectedOutput := []ImportStatement{
		{
			SourcePath: "bun-framework",
			NamedImports: []NamedImport{
				{
					Name: "createRoute",
				},
			},
		},
		{
			NamedImports: []NamedImport{
				{
					Name: "someFunction",
				},
				{
					Name: "someConstant",
				},
			},
			SourcePath: "../../some-function",
		},
		{
			DefaultImport: "defaultAsterisk",
			SourcePath:    "fs",
		},
		{
			NamedImports: []NamedImport{
				{
					Name: "namedNonTypeImport",
				},
			},
			SourcePath: "mixed-type-and-non-type-imports",
		},
		{
			NamedImports: []NamedImport{
				{
					Alias: "Alias",
					Name:  "aliasedImport",
				},
			},
			SourcePath: "aliased-import",
		},
		{
			NamedImports: []NamedImport{
				{
					Name: "nonAliasedImport",
				},
				{
					Alias: "NonTypeAliasedImport",
					Name:  "nonTypeAliasedIimport",
				},
			},
			SourcePath: "aliased-import",
		},
		{
			DefaultImport: "path",
			NamedImports: []NamedImport{
				{
					Name: "join",
				},
			},
			SourcePath: "path",
		},
		{
			SourcePath: "ts-morph",
			NamedImports: []NamedImport{
				{
					Name: "Project",
				},
			},
		},
		{
			SourcePath: "as",
			NamedImports: []NamedImport{
				{
					Name:  "join",
					Alias: "as",
				},
			},
		},
	}

	file, err := os.ReadFile(filepath.Join(".", "parseImportStatement_mock.ts"))
	if err != nil {
		t.Fatal(err)
	}
	runes := []rune(string(file))
	RemoveComments(&runes)
	importStatements := ReadImports(&runes, 0)

	if len(importStatements) != len(expectedOutput) {
		t.Fatalf("Expected %d import statements, got %d", len(expectedOutput), len(importStatements))
	}

	for i, tc := range expectedOutput {
		got := importStatements[i]
		if got.SourcePath != tc.SourcePath {
			t.Errorf("SourcePath: expected\n%s, got\n%s", tc.SourcePath, got.SourcePath)
		}
		if len(got.NamedImports) != len(tc.NamedImports) {
			t.Errorf("NamedImports sizes are different: expected\n%d, got\n%d", len(tc.NamedImports), len(got.NamedImports))
		}
		if got.DefaultImport != tc.DefaultImport {
			t.Errorf("DefaultImport: expected\n%s, got\n%s", tc.DefaultImport, got.DefaultImport)
		}
		for j, namedImport := range got.NamedImports {
			if namedImport.Alias != tc.NamedImports[j].Alias {
				t.Errorf("NamedImports[%d].Alias: expected\n%s, got\n%s", j, tc.NamedImports[j].Alias, namedImport.Alias)
			}
			if namedImport.Name != tc.NamedImports[j].Name {
				t.Errorf("NamedImports[%d].Name: expected\n%s, got\n%s", j, tc.NamedImports[j].Name, namedImport.Name)
			}
		}

	}
}
