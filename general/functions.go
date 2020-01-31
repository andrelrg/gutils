package general

import (
	"encoding/json"
	"os"
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

