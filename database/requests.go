package database

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/hugopukito/golang-fixture/color"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

var specialTypes = make(map[string]any)
var currentRegex *regexp.Regexp
var randCommasRegex *regexp.Regexp
var randRangeRegex *regexp.Regexp
var refRegex *regexp.Regexp
var newRegex *regexp.Regexp
var hashRegex *regexp.Regexp

func init() {
	addFuncsToSpecialTypes()
	compileRegex()
}

func addFuncsToSpecialTypes() {
	generateUUID := func() any {
		return uuid.New()
	}
	specialTypes["uuid.UUID"] = generateUUID

	generateTimeStamp := func() any {
		return time.Now()
	}
	specialTypes["time.Time"] = generateTimeStamp
}

func compileRegex() {
	currentRegex, err = regexp.Compile(`\{current\}`)
	randCommasRegex, err = regexp.Compile(`\{random\{([^}]*)\}\}`)
	randRangeRegex, err = regexp.Compile(`\{random\{((?:\d+(?:\.\d+)?\.\.\d+(?:\.\d+)?))\}\}`)
	refRegex, err = regexp.Compile(`{ref{([^}]*)}}`)
	newRegex, err = regexp.Compile(`\{new\{\}\}`)
	hashRegex, err = regexp.Compile(`{hash{([^}]*)}}`)
	if err != nil {
		log.Fatalln("Failed to compile regular expression:", err)
	}
}

func InsertEntity(structName string, entity map[string]any, localStruct map[string]string, occurrence int) error {
	columns := make([]string, 0)
	values := make([]any, 0)
	placeholders := make([]string, 0)

	// Auto generate id if not given in yaml
	for k := range localStruct {
		if _, exist := entity[k]; !exist && k == "id" {
			columns = append(columns, "id")
			values = append(values, uuid.New().String())
			placeholders = append(placeholders, "?")
		}
	}

	// Fields from entity
	for column, value := range entity {
		_, isString := value.(string)

		if isString {
			refMatches := refRegex.FindAllStringSubmatch(value.(string), -1)
			if len(refMatches) > 0 {
				column += "_id"
			}
		}
		columns = append(columns, camelToSnake(column))

		if isString {
			if currentRegex.MatchString(value.(string)) {
				values = append(values, strings.ReplaceAll(value.(string), "{current}", strconv.Itoa(occurrence)))
			} else if randCommasMatches := randCommasRegex.FindAllStringSubmatch(value.(string), -1); len(randCommasMatches) > 0 {
				if randRangeMatches := randRangeRegex.FindAllStringSubmatch(value.(string), -1); len(randRangeMatches) > 0 {
					innerContent := randCommasMatches[0][1]
					randomValues := splitAndTrim(innerContent, "..")
					randomVal, err := generateRandom(randomValues, localStruct[column])
					if err != nil {
						return err
					}
					values = append(values, randomVal)
				} else if strings.Contains(randCommasMatches[0][1], ",") {
					innerContent := randCommasMatches[0][1]
					randomValues := splitAndTrim(innerContent, ",")
					//castValuesInGoodType(randomValues, localStruct[column])
					values = append(values, getRandomElement(randomValues))
				} else {
					randomVal, err := generateRandom(nil, localStruct[column])
					if err != nil {
						return err
					}
					values = append(values, randomVal)
				}
			} else if refMatches := refRegex.FindAllStringSubmatch(value.(string), -1); len(refMatches) > 0 {
				values = append(values, refMatches[0][1])
			} else if newRegex.MatchString(value.(string)) {
				if value, ok := specialTypes[localStruct[column]]; ok {
					values = append(values, value.(func() any)())
				} else {
					fmt.Println(color.Red + "can't generate new for type: " + color.Orange + localStruct[column] + color.Reset)
				}
			} else if hashMatches := hashRegex.FindAllStringSubmatch(value.(string), -1); len(hashMatches) > 0 {
				hashPwd, err := generateHashPassword(hashMatches[0][1])
				if err != nil {
					return err
				}
				values = append(values, hashPwd)
			} else {
				values = append(values, value)
			}
		} else {
			values = append(values, value)
		}
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", structName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err = sqlConn.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func CheckTableExist(tableName string, dbName string) (bool, error) {
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"

	var count int
	err := sqlConn.QueryRow(query, dbName, tableName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateTable(tableName string, localStructs map[string]map[string]string, localStructOrdered []string) error {
	fmt.Println(color.Purple+"Creating table ->", color.Yellow, tableName+"..."+color.Reset)
	columns := make([]string, 0, len(localStructs[tableName]))

	//fmt.Println(localStructs)

	for _, columnName := range localStructOrdered {
		columnType := localStructs[tableName][columnName]

		sqlType, exist := GoSQLTypeMap[columnType]
		if !exist {
			structType, structTypeExist := localStructs[strings.ToLower(columnType)]
			if structTypeExist {
				structTypeKeys := make([]string, 0, len(structType))
				for structTypekey := range structType {
					structTypeKeys = append(structTypeKeys, structTypekey)
				}
				if primaryColumn := getIfExistPrimaryColumn(structTypeKeys); primaryColumn != "" {
					if structFieldType, ok := GoSQLTypeMap[structType[primaryColumn]]; ok {
						columns = append(columns, fmt.Sprintf("%s %s", columnName+"_id", structFieldType))
						continue
					} else {
						return errors.New(color.Red + "sql type for type: " + color.Orange + columnType + color.Red + " doesn't exist")
					}
				} else {
					return errors.New(color.Red + "no primary key: id" + " for struct: " + color.Orange + columnName)
				}
			} else {
				return errors.New(color.Red + "sql type for type: " + color.Orange + columnType + color.Red + " doesn't exist")
			}
		}
		columns = append(columns, fmt.Sprintf("%s %s", camelToSnake(columnName), sqlType))
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))

	_, err := sqlConn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func splitAndTrim(input, sep string) []string {
	var result []string
	parts := strings.Split(input, sep)
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getRandomElement(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(strings))
	randomElement := strings[randomIndex]

	return randomElement
}

func generateRandom(values []string, targetType string) (any, error) {
	if targetType == "bool" {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		return rand.Intn(2) == 1, nil
	} else if targetType == "int" || targetType == "float64" {
		return generateRandomNumber(values, targetType)
	}

	return nil, errors.New("can't generate random for type: " + targetType)
}

func generateRandomNumber(values []string, targetType string) (any, error) {
	if len(values) != 2 {
		return nil, errors.New("need 2 values for {random{}}")
	}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	if targetType == "int" {
		first, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, err
		}
		second, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}
		return rand.Intn(second) + first, nil
	} else if targetType == "float64" {
		first, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return nil, err
		}
		second, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return nil, err
		}
		return math.Round((rand.Float64()*(second-first)+first)*100) / 100, nil
	}

	return nil, errors.New("problem with your random{x..y}")
}

func getIfExistPrimaryColumn(slice []string) string {
	for _, v1 := range slice {
		if v1 == "id" {
			return v1
		}
	}
	return ""
}

func camelToSnake(camel string) string {
	var result []rune
	for i, r := range camel {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
