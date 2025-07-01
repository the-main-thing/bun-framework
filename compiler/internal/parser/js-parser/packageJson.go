package jsparser

import (
	"encoding/json"
	"errors"
	"os"
)

type PackageJson struct {
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
}

func ReadPackageJson(filePath string) (PackageJson, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return PackageJson{}, errors.New("Error reading package.json: " + filePath + "\n" + err.Error())
	}

	if fileBytes == nil || len(fileBytes) == 0 {
		return PackageJson{}, errors.New("package.json file is empty: " + filePath)
	}

	var packageJson PackageJson
	err = json.Unmarshal(fileBytes, &packageJson)
	if err != nil {
		return PackageJson{}, errors.New("Can't unmarshal package.json after removing all the comments. Likely the file is corrupted")
	}

	return packageJson, nil
}

func GetPackagesList(filePath string) []string {
	packageJson, err := ReadPackageJson(filePath)
	if err != nil {
		return make([]string, 0, 0)
	}
	size := len(packageJson.Dependencies) + len(packageJson.DevDependencies) + len(packageJson.PeerDependencies)
	packages := make([]string, size, size)
	index := 0
	for dependencyName := range packageJson.Dependencies {
		packages[index] = dependencyName
		index += 1
	}
	for dependencyName := range packageJson.DevDependencies {
		packages[index] = dependencyName
		index += 1
	}
	for dependencyName := range packageJson.PeerDependencies {
		packages[index] = dependencyName
		index += 1
	}
	return packages
}
