package funcs

import (
	"fixture/color"
	"fixture/database"
	"fmt"
	"log"
)

func ParseFixture(yamlFixture Fixture, dbName string) {
	for structName := range yamlFixture.Entities {
		localStruct, exist := GetLocalStructByName(structName)
		if exist {
			ensureTableIsCreated(structName, localStruct, dbName)
		} else {
			fmt.Println(color.Red+"Unknown struct ->", color.Orange, structName+"...")
		}
	}
	fmt.Println()
	for structName, entityMap := range yamlFixture.Entities {
		localStruct, exist := GetLocalStructByName(structName)
		if exist {
			fmt.Println(color.Green+"Processing struct ->", color.Yellow, structName+"..."+color.Reset)
			for entityName, fieldsAndValues := range entityMap {
				if CheckEntityOfStructIsValid(structName, fieldsAndValues, entityName) {
					err := database.InsertEntity(structName, fieldsAndValues, localStruct)
					fmt.Println(color.Cyan+"Adding entity ->", color.Yellow, entityName+"..."+color.Reset)
					if err != nil {
						fmt.Println(color.Red+"failed creating entity: "+color.Orange+entityName, color.Red+err.Error()+color.Reset)
					}
				}
			}
		}
		fmt.Println()
	}
}

func ensureTableIsCreated(structName string, localStruct map[string]string, dbName string) {
	exist, err := database.CheckTableExist(structName, dbName)
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
