package general

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
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

// Contains is to verify if an iten or more exists in a slice or array.
// The first param is an Slice or Array of any type.
// The "second" is a variadic arguments to verify with the content of the first argument
// The return is true if all variadic arguments exists in your slice or array
func Contains(slc interface{}, verItems ...interface{}) (bool, error) {
	s := reflect.ValueOf(slc)
	if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
		return false, errors.New("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	itensToCheck := len(verItems)
	checks := 0
	for _, r := range ret {
		for _, v := range verItems {
			if v == r {
				checks++
				if checks == itensToCheck {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
