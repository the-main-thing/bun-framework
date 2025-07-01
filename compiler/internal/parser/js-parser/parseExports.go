package jsparser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"slices"

	"github.com/evanw/esbuild/pkg/api"
)

type Exports struct {
	Default bool
	Named   []string
}

type metafileResp struct {
	Outputs map[string]struct {
		Exports []string `json:"exports"`
	} `json:"outputs"`
}

func ParseExports(tsFilePath string, source string) Exports {
	ext := filepath.Ext(tsFilePath)
	var loader api.Loader
	if ext == ".tsx" {
		loader = api.LoaderTSX
	} else {
		loader = api.LoaderTS
	}

	outfilePath := filepath.Join(os.TempDir(), "bun-framework-compiler", strings.ReplaceAll(tsFilePath, string(filepath.Separator), "_")+"-out.js")

	res := api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   source,
			Loader:     loader,
			ResolveDir: filepath.Dir(tsFilePath),
			Sourcefile: filepath.Base(tsFilePath),
		},
		Bundle:   false,
		Write:    false,
		Metafile: true,
		Format:   api.FormatESModule,
		Platform: api.PlatformNeutral,
		Outfile:  outfilePath,
	})
	go func() {
		err := os.Remove(outfilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error removing temporary file that has been created for parsing exports by esbuild")
			panic(err)
		}
	}()
	if len(res.Errors) > 0 {
		return Exports{}
	}

	var metafile metafileResp
	if err := json.Unmarshal([]byte(res.Metafile), &metafile); err != nil {
		return Exports{}
	}

	outExports := Exports{
		Default: false,
		Named:   make([]string, 0, 8),
	}
	for _, out := range metafile.Outputs {
		for _, name := range out.Exports {
			if name == "default" {
				outExports.Default = true
			} else {
				// small arrays are much faster than maps/sets. Fits in the cache, no need to calculate the hash
				if !slices.Contains(outExports.Named, name) {
					outExports.Named = append(outExports.Named, name)
				}
			}
		}
	}

	if !outExports.Default && len(outExports.Named) == 0 {
		return Exports{}
	}

	return outExports
}
