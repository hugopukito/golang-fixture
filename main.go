package main

import (
	"fixture/color"
	"fixture/database"
	"fixture/funcs"
	"fmt"
	"log"
)

func main() {

	fmt.Println(color.Pink + "Testing connection to your sql..." + color.Reset)
	database.InitDB("fixture")

	fmt.Println(color.Purple + "Parsing your local strucs..." + color.Reset)
	funcs.InitLocalStructs("structs")

	yamlFixtures, err := funcs.GetYamlStructs("fixtures")
	if err != nil {
		log.Panicln(color.Red + "GetYamlStructs err: " + err.Error() + color.Reset)
	}

	fmt.Println(color.Blue + "Parsing your fixtures... \n" + color.Reset)
	for _, yamlFixture := range yamlFixtures {
		funcs.ParseFixture(yamlFixture)
	}
}
