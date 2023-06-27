package funcs

import (
	"fixture/color"
	"fixture/database"
	"fmt"
	"log"
)

func ParseFixture(yamlFixture Fixture) {
	for structName, entityMap := range yamlFixture.Entities {
		localStruct, exist := GetLocalStructByName(structName)
		if exist {
			fmt.Println(color.Green+"Processing struct ->", color.Yellow, structName+"..."+color.Reset)
			for entityName, fieldsAndValues := range entityMap {
				fmt.Println(color.Cyan+"Adding entity ->", color.Yellow, entityName+"..."+color.Reset)
				if CheckEntityOfStructIsValid(structName, fieldsAndValues, entityName) {
					ensureTableIsCreated(structName, localStruct)
					database.InsertEntity(structName, fieldsAndValues)
				}
			}
		} else {
			fmt.Println(color.Red+"Unknown struct ->", color.Orange, structName+"...")
		}
		fmt.Println()
	}
}

func ensureTableIsCreated(structName string, localStruct map[string]string) {
	exist, err := database.CheckTableExist(structName)
	if err != nil {
		log.Panicln(color.Red+"failed creating table for structName: "+color.Orange+structName, color.Red+err.Error()+color.Reset)
	}
	if !exist {
		err = database.CreateTable(structName, localStruct)
		if err != nil {
			log.Panicln(color.Red+"failed creating table for structName: "+color.Orange+structName, color.Red+err.Error()+color.Reset)
		}
	}
}
