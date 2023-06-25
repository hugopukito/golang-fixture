package main

import (
	"fixture/structs"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	yamlFile, err := os.ReadFile("fixtures/dogs.yaml")
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	var data map[string]structs.EntityName

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}

	for entityPathName, entityPath := range data {
		fmt.Println(entityPathName)
		fmt.Println(entityPath.Entities)
	}
}
