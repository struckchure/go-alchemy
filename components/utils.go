package components

import (
	"bufio"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
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

func FormatGoCode(code string) (*string, error) {
	// Format the code using the go/format package
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return nil, fmt.Errorf("failed to format Go code: %w", err)
	}
	return lo.ToPtr(string(formatted)), nil
}

type GenerateTmplArgs struct {
	TmplPath   string
	OutputPath string
	Values     interface{}
	GoFormat   bool
}

func GenerateTmpl(args GenerateTmplArgs) error {
	tmplFileName := lo.Must(lo.Last(strings.Split(args.TmplPath, "/")))

	tmpl, err := template.New(tmplFileName).ParseFiles(args.TmplPath)
	if err != nil {
		return err
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
