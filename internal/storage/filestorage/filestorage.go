package filestorage

import (
	"encoding/json"
	"os"
)

func Save(filename string, data interface{}) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte(json)); err != nil {
		return err
	}
	if _, err := file.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}
