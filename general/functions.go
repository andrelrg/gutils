package general

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

//ReadConfig Reads Settings file.
func ReadConfigJson(configStruct interface{}, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configStruct)
	if err != nil {
		return err
	}

	return nil
}

// RemoveDuplicate remove duplicate entry in string array
func RemoveDuplicate(slice []string, verbose bool) []string {
	now := time.Now()
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	if verbose == true {
		elapsed := time.Now().Sub(now)
		log.Println("[removeDuplicate] total run time: ", fmt.Sprint(elapsed))
	}
	return list
}