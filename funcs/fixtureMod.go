package funcs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CheckIfFixturesGotModified(fixtureDirName string) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, errors.New("Error getting current working directory: " + err.Error())
	}
	structDir := wd + "/" + fixtureDirName

	files, err := os.ReadDir(structDir)
	if err != nil {
		return false, errors.New("Error reading directory: " + err.Error())
	}

	fileModMap := make(map[string]time.Time)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(structDir, file.Name())

		modTime, err := getFileModTime(filePath)
		if err != nil {
			return false, err
		}

		if isYAMLFile(file.Name()) {
			fileModMap[file.Name()] = modTime
		}
	}

	filePath := filepath.Join(structDir, "modTimeFixture")
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = writeOrCreateFile(filePath, fileModMap)
			if err != nil {
				return false, err
			}
			return true, nil
		} else {
			return false, err
		}
	}

	retrievedMap, err := readFromFile(filePath)
	if err != nil {
		return false, err
	}
	if checkIfFilesGotModified(fileModMap, retrievedMap) {
		err = writeOrCreateFile(filePath, fileModMap)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func checkIfFilesGotModified(newModTime map[string]time.Time, oldModTime map[string]time.Time) bool {
	if len(newModTime) != len(oldModTime) {
		return true
	}

	for key, value1 := range newModTime {
		value2, exists := oldModTime[key]
		if !exists || !value1.Equal(value2) {
			return true
		}
	}

	return false
}

func getFileModTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

func writeOrCreateFile(filePath string, fileModMap map[string]time.Time) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range fileModMap {
		line := fmt.Sprintf("%s: %s\n", key, value)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func readFromFile(filePath string) (map[string]time.Time, error) {
	data := make(map[string]time.Time)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			key := parts[0]
			valueStr := parts[1]

			valueTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", valueStr)
			if err != nil {
				continue
			}

			data[key] = valueTime
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return data, nil
}

func isYAMLFile(filename string) bool {
	extension := filepath.Ext(filename)
	return strings.EqualFold(extension, ".yml") || strings.EqualFold(extension, ".yaml")
}
