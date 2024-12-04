package alchemy

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// ModifyOrCreateEnvVar modifies or creates an environment variable in a .env file.
func ModifyOrCreateEnvVar(envFilePath, key, value string) error {
	// Open the .env file or create it if it doesn't exist
	file, err := os.OpenFile(envFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file content into memory
	var lines []string
	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+"=") {
			// Replace the line if the key already exists
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
			found = true
		} else {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Add the new key-value pair if it wasn't found
	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	// Write the modified content back to the file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func GetDirectoryName() string {
	return lo.Must(lo.Last(strings.Split(lo.Must(os.Getwd()), "/")))
}

type Yaml struct{}

func (y *Yaml) Read(fileName string, data interface{}) (*interface{}, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	err = yaml.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &data, nil
}

func (y *Yaml) Write(fileName string, data interface{}) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func NewYaml() *Yaml {
	return &Yaml{}
}