package helper

import (
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

	filePath := fmt.Sprintf("json/%d/%02d/%02d", y, m, d)
	fileName := fmt.Sprintf("%s.json", name)

	if err := os.MkdirAll(filePath, 0777); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(filePath, fileName))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}
