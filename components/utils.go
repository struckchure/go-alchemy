package components

import (
	"bufio"
	"errors"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"github.com/struckchure/go-alchemy"
)

func GetDirectoryName() string {
	return lo.Must(lo.Last(strings.Split(lo.Must(os.Getwd()), "/")))
}

// GetModuleName reads the go.mod file in the current directory
// and returns the module name.
func GetModuleName() (*string, error) {
	goModPath := "go.mod"

	// Open the go.mod file
	file, err := os.Open(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open go.mod file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// The module name is specified in the line starting with "module"
		if strings.HasPrefix(line, "module") {
			// Extract and return the module name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return &parts[1], nil
			}
			return nil, fmt.Errorf("invalid module declaration in go.mod")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading go.mod file: %w", err)
	}

	return nil, fmt.Errorf("module name not found in go.mod")
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

type GenerateTmplArgs struct {
	TmplPath   string
	OutputPath string
	Values     interface{}
	Funcs      map[string]any
	GoFormat   bool
}

func GenerateTmpl(args GenerateTmplArgs) error {
	var tmpl *template.Template

	baseUrlFromEnv := os.Getenv("ALCHEMY_TMPL_DIR")
	baseUrl := lo.Ternary(
		baseUrlFromEnv != "",
		baseUrlFromEnv,
		"https://raw.githubusercontent.com/struckchure/go-alchemy/refs/heads/main/",
	)

	tmplPath, err := JoinURLsOrPaths(baseUrl, args.TmplPath)
	if err != nil {
		return err
	}

	args.TmplPath = tmplPath

	isRemoteURL, err := IsRemoteURL(args.TmplPath)
	if err != nil {
		return err
	}

	if isRemoteURL {
		tmplFileName := lo.Must(lo.Last(strings.Split(args.TmplPath, "/")))
		content, err := ReadRemoteFile(args.TmplPath)
		if err != nil {
			return err
		}

		tmpl, err = template.New(tmplFileName).Funcs(template.FuncMap(args.Funcs)).Parse(*content)
		if err != nil {
			return err
		}
	} else {
		tmplFileName := lo.Must(lo.Last(strings.Split(args.TmplPath, "/")))

		tmpl, err = template.New(tmplFileName).Funcs(template.FuncMap(args.Funcs)).ParseFiles(args.TmplPath)
		if err != nil {
			return err
		}
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(args.OutputPath)
	err = os.MkdirAll(outputDir, 0755) // Create all missing directories with appropriate permissions
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Create or overwrite the output file
	file, err := os.Create(args.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", args.OutputPath, err)
	}
	defer file.Close()

	// Execute the template with the provided values
	err = tmpl.Execute(file, args.Values)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// If Go formatting is requested, format the file content
	if args.GoFormat {
		// Read the generated file content
		content, err := os.ReadFile(args.OutputPath)
		if err != nil {
			return fmt.Errorf("failed to read file for formatting: %w", err)
		}

		// Format the content using the FormatGoCode function
		formattedContent, err := FormatGoCode(string(content))
		if err != nil {
			return fmt.Errorf("failed to format Go code: %w", err)
		}

		// Write the formatted content back to the file
		err = os.WriteFile(args.OutputPath, []byte(*formattedContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write formatted content to file: %w", err)
		}
	}

	return nil
}

func UpdateComponentConfig(componentConfig Component) error {
	cfg, err := alchemy.ReadYaml[Config]("alchemy.yaml")
	if err != nil {
		return err
	}

	currentComponentConfig, componentExists := lo.Find(
		cfg.Components,
		func(c Component) bool { return c.Id == componentConfig.Id },
	)

	if componentExists {
		componentConfig.Models = lo.Uniq(append(componentConfig.Models, currentComponentConfig.Models...))
		componentConfig.Services = lo.Uniq(append(componentConfig.Services, currentComponentConfig.Services...))

		_, idx, ok := lo.FindIndexOf(
			cfg.Components,
			func(c Component) bool { return c.Id == componentConfig.Id },
		)
		if ok {
			cfg.Components[idx] = componentConfig
		}
	} else {
		cfg.Components = append(cfg.Components, componentConfig)
	}

	return nil
}
