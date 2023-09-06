package funcs

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	// Will create a map of entities which have entities calling them with ref
	// Then give each ref an uuid
	for _, structValue := range yamlFixture.Entities {
		for _, v := range structValue {
			for k2, v2 := range v {
				if _, isString := v2.(string); isString {
					if refMatches := refRegex.FindAllStringSubmatch(v2.(string), -1); len(refMatches) > 0 {
						refMatch := refMatches[0][1]

						// Assign to the parent of ref the uuid to inject it later because we don't know if it exist yet
						// Check if we already create an uuid for this ref, if not create one
						if _, ok := refEntities[refMatch]; !ok {
							refEntities[refMatch] = uuid.New().String()
						}

						// Assign to the entities calling ref the uuid
						v[k2] = strings.Replace(v2.(string), refMatch, refEntities[refMatch], -1)
					}
				}
			}
		}
	}

	// Take the uuid given and inject it in the parent of ref
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

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					err = entitiesFactory(ctx, startNumber, nbOfCreation, structName, entityName, fieldsAndValues, localStruct)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		fmt.Print("\n")
	}
}

func entitiesFactory(ctx context.Context, startNumber, nbOfCreation int, structName, entityName string, fieldsAndValues map[string]any, localStruct map[string]string) error {
	for i := startNumber; i < nbOfCreation+startNumber; i++ {

		select {
		case <-ctx.Done():
			percentage := (i - startNumber + 1) * 100 / nbOfCreation
			fmt.Printf("\r%d%%", percentage)
			if i == nbOfCreation+startNumber-1 {
				fmt.Printf("\r\033[K")
			}
		default:
		}

		var entityNameNumberized string
		if nbOfCreation == 1 && startNumber == 0 {
			entityNameNumberized = entityName
		} else {
			entityNameNumberized = entityName + strconv.Itoa(i)
		}
		if !addAndCheckEntityName(entityNameNumberized) {
			return errors.New(color.Red + "entity name already taken: " + color.Orange + entityNameNumberized + color.Reset)
		}
		err := database.InsertEntity(structName, fieldsAndValues, localStruct, i)
		if err != nil {
			return errors.New(color.Red + "failed creating entity: " + color.Orange + entityNameNumberized + color.Red + err.Error() + color.Reset)
		}
	}
	return nil
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
