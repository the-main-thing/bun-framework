package types

type TsConfig struct {
	CompilerOptions CompilerOptions `json:"compilerOptions"`
}

type CompilerOptions struct {
	BaseUrl string              `json:"baseUrl"`
	Paths   map[string][]string `json:"paths"`
}

type ResolvedImportAlias struct {
	Alias string
	Paths []string
}

type Export struct {
	Default bool
	Named []string
}
