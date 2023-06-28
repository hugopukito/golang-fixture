package funcs

import (
	"fixture/color"
	"fixture/database"
	"fmt"
	"log"
)

func ParseFixture(yamlFixture Fixture, dbName string) {
	for structName, entityMap := range yamlFixture.Entities {
		localStruct, exist := GetLocalStructByName(structName)
		if exist {
			fmt.Println(color.Green+"Processing struct ->", color.Yellow, structName+"..."+color.Reset)
			for entityName, fieldsAndValues := range entityMap {
				if CheckEntityOfStructIsValid(structName, fieldsAndValues, entityName) {
					ensureTableIsCreated(structName, localStruct, dbName)
					keysWithTypeUUID := retrieveKeysWithTypeUUID(localStruct)
					err := database.InsertEntity(structName, fieldsAndValues, keysWithTypeUUID)
					fmt.Println(color.Cyan+"Adding entity ->", color.Yellow, entityName+"..."+color.Reset)
					if err != nil {
						fmt.Println(color.Red+"failed creating entity: "+color.Orange+entityName, color.Red+err.Error()+color.Reset)
					}
				}
			}
		} else {
			fmt.Println(color.Red+"Unknown struct ->", color.Orange, structName+"...")
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

func retrieveKeysWithTypeUUID(localStruct map[string]string) []string {
	var keys []string
	for k, v := range localStruct {
		if v == "UUID" || v == "uuid.UUID" {
			keys = append(keys, k)
		}
	}
	return keys
}
