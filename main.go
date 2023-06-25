package main

import (
	"fixture/structs"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	yamlFile, err := os.ReadFile("fixtures/animals.yaml")
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	var fixture structs.Fixture

	err = yaml.Unmarshal(yamlFile, &fixture)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}
	for k, v := range fixture.Entities {
		fmt.Println("entity struct name: ", k)
		for k, v := range v {
			fmt.Println("entity name: ", k)
			fmt.Println("entity fields and values", v)
		}
	}
}
