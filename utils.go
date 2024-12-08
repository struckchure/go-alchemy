package alchemy

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"errors"
	"go/format"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// WriteEnvVar modifies or creates an environment variable in a .env file.
func WriteEnvVar(envFilePath, key, value string) error {
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

func RemoveNoneAlpha(i string) string {
	re := regexp.MustCompile(`[^\w]+`) // Matches anything that's not a word character
	return re.ReplaceAllString(i, "")
}

func ReadYaml[T any](fileName string) (*T, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data T
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &data, nil
}

func WriteYaml(fileName string, data interface{}) error {
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

func FormatGoCode(code string) (*string, error) {
	// Format the code using the go/format package
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return nil, fmt.Errorf("failed to format Go code: %w", err)
	}
	return lo.ToPtr(string(formatted)), nil
}

func ReadRemoteFile(url string) (*string, error) {
	// Send an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the file: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch the file, status: %s", resp.Status)
	}

	// Read the content of the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file content: %w", err)
	}

	return lo.ToPtr(string(body)), nil
}

// IsRemoteURL determines whether a given string represents a remote URL or a local path.
// Returns true if the URL is remote (e.g., HTTP/HTTPS); false if it is local.
func IsRemoteURL(input string) (bool, error) {
	if input == "" {
		return false, errors.New("input string is empty")
	}

	// Parse the input string as a URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false, err
	}

	// Check if the URL has a recognized remote scheme
	// Typically remote URLs have schemes like "http", "https", etc.
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		return true, nil
	}

	// If the scheme is empty, it's likely a local file path
	if parsedURL.Scheme == "" && strings.HasPrefix(input, "/") {
		return false, nil
	}

	return false, nil
}

// JoinURLsOrPaths joins multiple URLs or paths dynamically.
// The first element determines whether the result will be an HTTP URL or a local path.
func JoinURLsOrPaths(base string, segments ...string) (string, error) {
	if base == "" {
		return "", errors.New("base URL/path is empty")
	}
	if len(segments) == 0 {
		return base, nil // If no segments, return the base as is
	}

	parsedBase, err := url.Parse(base)
	if err != nil {
		return "", errors.New("invalid base URL/path")
	}

	// If the base is a remote URL
	if parsedBase.Scheme == "http" || parsedBase.Scheme == "https" {
		joinedPath := parsedBase.Path
		for _, segment := range segments {
			joinedPath = path.Join(joinedPath, segment)
		}
		parsedBase.Path = joinedPath
		return parsedBase.String(), nil
	}

	// If the base is a local file path
	allPaths := append([]string{base}, segments...)
	return filepath.Join(allPaths...), nil
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
