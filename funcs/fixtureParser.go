package funcs

import (
	"fixture/color"
	"fixture/database"
	"fmt"
)

func ParseFixture(yamlFixture Fixture) {
	for structName, entityMap := range yamlFixture.Entities {
		if CheckLocalStructNameExist(structName) {
			fmt.Println(color.Green+"Processing struct ->", color.Yellow, structName+"..."+color.Reset)
			for entityName, fieldsAndValues := range entityMap {
				fmt.Println(color.Cyan+"Adding entity ->", color.Yellow, entityName+"..."+color.Reset)
				if CheckEntityOfStructIsValid(structName, fieldsAndValues, entityName) {
					database.InsertEntity(structName, fieldsAndValues)
				}
			}
		} else {
			fmt.Println(color.Red+"Unknown struct ->", color.Orange, structName+"...")
		}
		fmt.Println()
	}
}
