package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/struckchure/go-alchemy"
	"gopkg.in/yaml.v3"
)

var CategoryMapping map[string]IAlchemyComponent = map[string]IAlchemyComponent{
	"authentication": NewAuthentication(),
}

type IConfigService interface {
	setupOrm(string, string) error
	provisionDatabase(string) error

	Init(InitArgs) error
	Add(string) error
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

func (c *ConfigService) provisionDatabase(databaseProvider string) error {
	color.Green("Creating %s with Docker Compose", databaseProvider)
	defer color.Green("Docker Compose file successfully generated")

	err := GenerateTmpl(GenerateTmplArgs{
		TmplPath:   "_templates/docker-compose.yaml.tmpl",
		OutputPath: "docker-compose.yaml",
		Values: map[string]interface{}{
			"ProjectName":      GetDirectoryName(),
			"DatabaseProvider": strings.ToLower(databaseProvider),
		},
		Funcs: map[string]any{"removeSigns": alchemy.RemoveNoneAlpha},
	})
	if err != nil {
		return err
	}

	return nil
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

	args.DatabaseUrl = lo.Ternary(
		args.DatabaseUrl != "",
		args.DatabaseUrl,
		fmt.Sprintf(`"postgresql://user:password@localhost:5432/%s?schema=public"`, alchemy.RemoveNoneAlpha(config.ProjectName)),
	)
	err = alchemy.WriteEnvVar(".env", "DATABASE_URL", args.DatabaseUrl)
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

// Adds single component to your project
//
//	`Authentication.Login` would be referring to only the login service
//	`Authentication` would add all available features in the authentication category
//
// The reference is case insenstive, we just like how the casing looks 🙂
func (c *ConfigService) Add(component string) error {
	component = strings.ToLower(component)

	var (
		categoryId  string
		componentId string
	)

	if len(strings.Split(component, ".")) > 1 {
		categoryId = strings.SplitN(component, ".", 2)[0]
		componentId = strings.SplitN(component, ".", 2)[1]
	} else {
		categoryId = component
		componentId = "all"
	}

	if !lo.HasKey(CategoryMapping, categoryId) {
		return errors.New("category is not available")
	}

	setupComponent := func(category IAlchemyComponent, id string) error {
		setup, err := category.Setup(id)
		if err != nil {
			return err
		}

		err = setup()
		if err != nil {
			color.Red("%s.%s setup failed", categoryId, id)
			return err
		}

		return nil
	}

	category := CategoryMapping[categoryId]
	category.Setup("no")

	if componentId == "all" {
		for _, componentId := range AuthenticationOptions {
			err := setupComponent(category, componentId)
			if err != nil {
				return err
			}
		}
	} else {
		err := setupComponent(category, componentId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ConfigService) Remove() error {
	return nil
}

func NewConfigService() IConfigService {
	return &ConfigService{}
}
