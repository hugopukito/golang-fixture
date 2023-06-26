package main

import (
	"fixture/funcs"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Fixture struct {
	Entities map[string]Entity `yaml:",inline"`
}

type Entity map[string]interface{}

func main() {

	funcs.InitLocalStructs("structs")

	yamlFile, err := os.ReadFile("fixtures/animals.yaml")
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	var yamlFixture Fixture

	err = yaml.Unmarshal(yamlFile, &yamlFixture)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}

	for structName, entityMap := range yamlFixture.Entities {
		structFields, exist := funcs.GetFieldsFromStructName(structName)
		if exist {
			fmt.Println(structName, structFields)
			fmt.Println()
		}
		for entityName, fieldsAndValues := range entityMap {
			fmt.Print(entityName + " ")
			fmt.Println(fieldsAndValues)
		}
		fmt.Println()
	}
}
