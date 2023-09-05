package funcs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/hugopukito/golang-fixture/color"
	"gopkg.in/yaml.v3"
)

func GetYamlStructs(fixtureDirName string) (Fixture, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Fixture{}, errors.New("Error getting current working directory: " + err.Error())
	}
	structDir := wd + "/" + fixtureDirName

	files, err := os.ReadDir(structDir)
	if err != nil {
		return Fixture{}, errors.New("Error reading directory: " + err.Error())
	}

	var yamlFixtures Fixture
	yamlFixtures.Entities = make(map[string]Entity)

	for _, file := range files {
		if file.IsDir() || !isYAMLFile(file.Name()) {
			continue
		}

		var yamlFixture Fixture

		filePath := filepath.Join(structDir, file.Name())

		yamlFile, err := os.ReadFile(filePath)
		if err != nil {
			return Fixture{}, err
		}

		err = yaml.Unmarshal(yamlFile, &yamlFixture)
		if err != nil {
			return Fixture{}, err
		}

		newMap := make(map[string]Entity, len(yamlFixture.Entities))

		for k, v := range yamlFixture.Entities {
			lowercaseKey := strings.ToLower(k)
			newMap[lowercaseKey] = v
		}
		yamlFixture.Entities = newMap

		for entityKey, entity := range yamlFixture.Entities {
			if oldEntity, exists := yamlFixtures.Entities[entityKey]; exists {
				for k, v := range entity {
					// check if entity key exist
					if _, entityExist := oldEntity[k]; !entityExist {
						oldEntity[k] = v
					} else {
						return Fixture{}, errors.New("Error two entities same name key: " + color.Orange + k)
					}
				}
				yamlFixtures.Entities[entityKey] = oldEntity
			} else {
				yamlFixtures.Entities[entityKey] = entity
			}
		}
	}

	return yamlFixtures, nil
}
