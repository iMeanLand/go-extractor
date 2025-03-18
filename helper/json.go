package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func SaveToJsonFile(name string, data interface{}) error {
	now := time.Now()
	y := now.Year()
	m := now.Month()
	d := now.Day()

	filePath := fmt.Sprintf("json/%d/%d/%d", y, m, d)

	if err := os.MkdirAll(filepath.Dir(filePath), 0777); err != nil {
		return err
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(b)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", " ")

	file, err := os.Create(filePath + "/" + name + ".json")
	defer file.Close()

	if err != nil {
		return err
	}

	//err = file.Write(buffer)
	//if err != nil {
	//	return err
	//}

	return nil
}
