package funcs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func GetYamlStructs(pkgName string) ([]Fixture, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("Error getting current working directory: " + err.Error())
	}
	structDir := wd + "/" + pkgName

	files, err := os.ReadDir(structDir)
	if err != nil {
		return nil, errors.New("Error reading directory: " + err.Error())
	}

	var yamlFixtures []Fixture

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var yamlFixture Fixture

		filePath := filepath.Join(structDir, file.Name())

		yamlFile, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading YAML file:", err)
			continue
		}

		err = yaml.Unmarshal(yamlFile, &yamlFixture)
		if err != nil {
			fmt.Println("Error unmarshaling YAML:", err)
		}

		newMap := make(map[string]Entity, len(yamlFixture.Entities))

		for k, v := range yamlFixture.Entities {
			lowercaseKey := strings.ToLower(k)
			newMap[lowercaseKey] = v
		}
		yamlFixture.Entities = newMap

		yamlFixtures = append(yamlFixtures, yamlFixture)
	}

	return yamlFixtures, nil
}
