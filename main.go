package main

import (
	"fixture/funcs"
	"fixture/structs"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {

	funcs.InitLocalStructs("structs")

	yamlFile, err := os.ReadFile("fixtures/animals.yaml")
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	var yamlFixture structs.Fixture

	err = yaml.Unmarshal(yamlFile, &yamlFixture)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}

	// structuredFixture := structs.Fixture{
	// 	Entities: make(map[string]structs.Entity),
	// }

	for k := range yamlFixture.Entities {
		// fmt.Println("entity struct name: ", k)
		// for k, v := range v {
		// 	fmt.Println("entity name: ", k)
		// 	fmt.Println("entity fields and values", v)
		// }
		funcs.StructAssign(k)
	}
}
