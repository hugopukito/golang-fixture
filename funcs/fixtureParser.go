package funcs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hugopukito/golang-fixture/color"
	"github.com/hugopukito/golang-fixture/database"
)

var entityNames map[string]struct{}

func ParseFixture(yamlFixture Fixture, dbName string) {
	entityNames = make(map[string]struct{})
	err := createTables(yamlFixture, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}
	yamlFixture, err = createTableRefIDs(yamlFixture)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("\n")
	createEntities(yamlFixture)
}

func createTables(yamlFixture Fixture, dbName string) error {
	for structName := range yamlFixture.Entities {
		localStructOrdered, exist := GetLocalStructOrderByName(structName)
		if exist {
			err := ensureTableIsCreated(structName, localStructOrdered, dbName)
			if err != nil {
				return errors.New(color.Red + "failed creating table for structName: " + color.Orange + structName + " " + color.Red + err.Error() + color.Reset)
			}
		} else {
			return errors.New(color.Red + "Unknown struct ->" + color.Orange + structName + "...")
		}
	}
	return nil
}

func createTableRefIDs(yamlFixture Fixture) (Fixture, error) {
	fmt.Println(color.Purple + "Creating table ref IDs..." + color.Reset)

	refEntities := make(map[string]string, 0)

	for _, structValue := range yamlFixture.Entities {
		for _, v := range structValue {
			for k2, v2 := range v {
				if _, isString := v2.(string); isString {
					if refMatches := refRegex.FindAllStringSubmatch(v2.(string), -1); len(refMatches) > 0 {
						refMatch := refMatches[0][1]
						uuid := uuid.New().String()
						v[k2] = strings.Replace(v2.(string), refMatch, uuid, -1)
						refEntities[refMatch] = uuid
					}
				}
			}
		}
	}

	for ref, uuid := range refEntities {
		refFound := false
		for structName, structValue := range yamlFixture.Entities {
			if v, ok := structValue[ref]; ok {
				refFound = true
				if localStruct, exist := GetLocalStructByName(structName); exist {
					if localStruct["id"] == "uuid.UUID" {
						v["id"] = uuid
					} else {
						return Fixture{}, errors.New(color.Red + "This ref don't have id field: " + color.Orange + ref + color.Reset)
					}
				}
			}
		}
		if !refFound {
			return Fixture{}, errors.New(color.Red + "Ref not found: " + color.Orange + ref + color.Reset)
		}
	}

	return yamlFixture, nil
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
						if nbOfCreation == 1 && startNumber == 0 {
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

func ensureTableIsCreated(structName string, localStructOrdered []string, dbName string) error {
	exist, err := database.CheckTableExist(structName, dbName)
	if err != nil {
		return err
	}
	if !exist {
		err = database.CreateTable(structName, structMap, localStructOrdered)
		if err != nil {
			return err
		}
	}
	return nil
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
