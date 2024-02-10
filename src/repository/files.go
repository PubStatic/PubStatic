package repository

import (
	"os"
)

func WriteFile(fileName string, content string) error {
	file, err := os.Create("config/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close() // Defer closing the file until the function exits

	// Write to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(fileName string) (string, error) {
	data, err := os.ReadFile("config/" + fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
