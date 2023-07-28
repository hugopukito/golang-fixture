package funcs

import (
	"fixture/color"
	"fixture/database"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

var entityNames map[string]struct{}

func ParseFixture(yamlFixture Fixture, dbName string) {
	entityNames = make(map[string]struct{})
	createTables(yamlFixture, dbName)
	fmt.Print("\n")
	createEntities(yamlFixture)
}

func createTables(yamlFixture Fixture, dbName string) {
	for structName := range yamlFixture.Entities {
		localStruct, exist := GetLocalStructByName(structName)
		localStructOrdered, exist2 := GetLocalStructOrderByName(structName)
		if exist && exist2 {
			ensureTableIsCreated(structName, localStruct, localStructOrdered, dbName)
		} else {
			fmt.Println(color.Red+"Unknown struct ->", color.Orange, structName+"...")
		}
	}
}

func createEntities(yamlFixture Fixture) {
	for structName, entityMap := range yamlFixture.Entities {
		localStruct, exist := GetLocalStructByName(structName)
		if exist {
			fmt.Println(color.Green+"Processing struct ->", color.Yellow, structName+"..."+color.Reset)
			for entityName, fieldsAndValues := range entityMap {
				if CheckEntityOfStructIsValid(structName, fieldsAndValues, entityName) {
					digits, foundDigits := extractPatternDigits(entityName)
					nbOfCreation := digits[1] - digits[0] + 1
					startNumber := digits[0]
					if nbOfCreation <= 0 {
						fmt.Println(color.Red + "failed creating entity, bad {x..y} 0 or negative: " + color.Orange + entityName + color.Reset)
					}

					entityName = removePattern(entityName)

					if foundDigits {
						fmt.Println(color.Cyan+"Adding entities ->", color.Yellow, entityName, "from", digits[0], "to", digits[1], "..."+color.Reset)
					} else {
						fmt.Println(color.Cyan+"Adding entity ->", color.Yellow, entityName+"..."+color.Reset)
					}

					for i := startNumber; i < nbOfCreation+startNumber; i++ {
						var entityNameNumberized string
						if nbOfCreation == 1 && startNumber != 0 {
							entityNameNumberized = entityName
						} else {
							entityNameNumberized = entityName + strconv.Itoa(i)
						}
						if !addAndCheckEntityName(entityNameNumberized) {
							fmt.Println(color.Red + "entity name already taken: " + color.Orange + entityNameNumberized + color.Reset)
						}
						err := database.InsertEntity(structName, fieldsAndValues, localStruct, i)
						if err != nil {
							fmt.Println(color.Red+"failed creating entity: "+color.Orange+entityNameNumberized, color.Red+err.Error()+color.Reset)
						}
					}
				}
			}
		}
		fmt.Print("\n")
	}
}

func ensureTableIsCreated(structName string, localStruct map[string]string, localStructOrdered []string, dbName string) {
	exist, err := database.CheckTableExist(structName, dbName)
	if err != nil {
		log.Panicln(color.Red+"failed creating table for structName: "+color.Orange+structName, color.Red+err.Error()+color.Reset)
	}
	if !exist {
		err = database.CreateTable(structName, localStruct, localStructOrdered)
		if err != nil {
			log.Panicln(color.Red+"failed creating table for structName: "+color.Orange+structName, color.Red+err.Error()+color.Reset)
		}
	}
}

func addAndCheckEntityName(entityName string) bool {
	_, exists := entityNames[entityName]
	if !exists {
		entityNames[entityName] = struct{}{}
		return true
	}
	return false
}

func extractPatternDigits(input string) ([2]int, bool) {
	pattern := `\{(\d+)\.\.(\d+)\}`

	r := regexp.MustCompile(pattern)
	match := r.FindStringSubmatch(input)
	if match == nil {
		return [2]int{0, 0}, false
	}

	start, _ := strconv.Atoi(match[1])
	end, _ := strconv.Atoi(match[2])

	digits := [2]int{start, end}
	return digits, true
}

func removePattern(input string) string {
	pattern := `\{\d+\.\.\d+\}`

	r := regexp.MustCompile(pattern)
	output := r.ReplaceAllString(input, "")

	return output
}
