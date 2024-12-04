package internals

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type IConfigService interface {
	setupOrm(string, string) error
	provisionDatabase(string) error

	Load() (*Config, error)
	Init(InitArgs) error
	Add([]string) error
	Remove() error
}

type ConfigService struct{}

func (c *ConfigService) setupPrisma(databaseProvider string) error {
	color.Green("Downloading go prisma client")
	cmd := exec.Command("go", "get", "github.com/steebchen/prisma-client-go")
	if err := cmd.Run(); err != nil {
		color.Red("%s", err)

		return err
	}

	color.Green("Initializing new prisma project")
	cmd = exec.Command(
		"go", "run", "github.com/steebchen/prisma-client-go",
		"init", "--datasource-provider", strings.ToLower(databaseProvider),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		color.Red("%s", string(out))

		return err
	}

	return nil
}

func (c *ConfigService) setupOrm(orm string, databaseProvider string) error {
	switch strings.ToLower(orm) {
	case "prisma":
		color.Green("Using Prisma ORM")

		err := c.setupPrisma(databaseProvider)
		if err != nil {
			color.Red("%s", err)

			return err
		}
	case "gorm":
		color.Green("Using GORM")
	default:
		err := errors.New("orm is not supported")
		color.Red("%s", err)

		return err
	}

	return nil
}

type dockerComposeTmplArgs struct {
	ProjectName      string
	DatabaseProvider string
}

func (c *ConfigService) provisionDatabase(databaseProvider string) error {
	color.Green("Creating %s with Docker Compose", databaseProvider)

	dockerComposeTmplFile := "internals/_templates/docker-compose.yaml.tmpl"
	dockerComposeFile := "docker-compose.yaml"

	funcMap := template.FuncMap{
		"removeSigns": func(input string) string {
			re := regexp.MustCompile(`[^\w]+`) // Matches anything that's not a word character
			return re.ReplaceAllString(input, "")
		},
	}

	tmpl, err := template.New("docker-compose.yaml.tmpl").Funcs(funcMap).ParseFiles(dockerComposeTmplFile)
	if err != nil {
		color.Red("%s", err)

		return err
	}

	// Create or overwrite the output file
	file, err := os.Create(dockerComposeFile)
	if err != nil {
		color.Red("Failed to create output file: %s", err)
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, dockerComposeTmplArgs{
		ProjectName:      GetDirectoryName(),
		DatabaseProvider: strings.ToLower(databaseProvider),
	})
	if err != nil {
		color.Red("%s", err)

		return err
	}

	color.Green("Docker Compose file successfully written to %s", dockerComposeFile)

	return nil
}

func (c *ConfigService) Load() (*Config, error) {
	cfg := Config{}

	fileContent, _ := os.ReadFile("alchemy.yaml")
	err := yaml.Unmarshal([]byte(fileContent), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

type InitArgs struct {
	Root                  string
	Orm                   string
	ShouldProvideDatabase bool
	DatabaseUrl           string
	DatabaseProvider      string
}

func (c *ConfigService) Init(args InitArgs) error {
	config := Config{
		ProjectName: GetDirectoryName(),
		Root:        lo.Ternary(args.Root == "", ".", args.Root),
		Orm:         Orm{Name: args.Orm, DatabaseProvider: args.DatabaseProvider},
	}

	if args.ShouldProvideDatabase {
		err := c.provisionDatabase(config.Orm.DatabaseProvider)
		if err != nil {
			color.Red("%s", err)

			return err
		}
	}

	err := c.setupOrm(config.Orm.Name, config.Orm.DatabaseProvider)
	if err != nil {
		color.Red("%s", err)

		return err
	}

	err = ModifyOrCreateEnvVar(".env", "DATABASE_URL", args.DatabaseUrl)
	if err != nil {
		color.Red("%s", err)

		return err
	}

	fileName := "alchemy.yaml"
	file, err := os.Create(fileName)
	if err != nil {
		color.Red("%s", err)

		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(config)
	if err != nil {
		color.Red("%s", err)

		return err
	}

	color.Green("✨ Alchemy config has been generated!")

	color.Green("🛠️  Updating Go dependencies ...")
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		color.Red("%s", err)

		return err
	}

	color.Green("🥂 You're all set!")

	if args.ShouldProvideDatabase {
		color.Green(`
Start Database Service
$ docker compose up -d
		`)
	}

	color.Green(`
Interactively add component
$ alchemy add Authentication // this will add all components from the authentication module

Or add a specific component
$ alchemy add Authentication.Login
	`)

	return nil
}

func (c *ConfigService) Add(components []string) error {
	return nil
}

func (c *ConfigService) Remove() error {
	return nil
}

func NewConfigService() IConfigService {
	return &ConfigService{}
}
