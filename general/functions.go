package general

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
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

// ConvertIntInterfaceToString tries to convert an interface value to string
func ConvertIntInterfaceToString(value interface{}) (string, error) {
	var err error
	var stringResult string
	var intResult int
	var longResult int64

	switch i := value.(type) {
	case int:
		intResult, _ = value.(int)
		stringResult = strconv.Itoa(intResult)
		return stringResult, nil
	case int64:
		longResult, _ = value.(int64)
		return strconv.FormatInt(longResult, 10), nil
	case float32:
		return strconv.FormatInt(int64(i), 10), nil
	case float64:
		return strconv.FormatInt(int64(i), 10), nil
	}

	longResult, err = strconv.ParseInt(fmt.Sprint(value), 10, 64)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return strconv.FormatInt(longResult, 10), nil
}

// ConvertInterfaceToInt64 converts an interface value to int64
func ConvertInterfaceToInt64(value interface{}) int64 {
	switch i := value.(type) {
	case int:
		return int64(i)
	case int64:
		return value.(int64)
	case string:
		long, _ := strconv.ParseInt(fmt.Sprint(value), 10, 64)
		return long
	case float32:
		return int64(value.(float32))
	case float64:
		return int64(value.(float64))
	}

	return 0
}

// InArray checks if a certain value exists inside the slice/array
func InArray(value interface{}, array interface{}) (bool, int) {
	exists := false
	index := -1

	if reflect.TypeOf(array).Kind() == reflect.Slice {
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				index = i
				exists = true
				return exists, index
			}
		}
	}
	return exists, index
}