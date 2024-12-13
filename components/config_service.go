package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/struckchure/go-alchemy/internals"
)

var CategoryMapping map[string]IAlchemyComponent = map[string]IAlchemyComponent{
	"Authentication": NewAuthentication(),
}

type IConfigService interface {
	setupOrm(Config) error
	provisionDatabase(Config) error

	Init(InitArgs) error
	Add(AddArgs) error
	Remove() error
}

type ConfigService struct{}

func (c *ConfigService) setupPrisma(databaseProvider string, directory string) error {
	color.Green("Downloading go prisma client")
	cmd := exec.Command("go", "get", "github.com/steebchen/prisma-client-go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(err, errors.New(string(out)))
	}

	color.Green("Initializing new prisma project [%s]", directory)
	err := os.Chdir(directory)
	if err != nil {
		return err
	}

	cmd = exec.Command(
		"go", "run", "github.com/steebchen/prisma-client-go",
		"init", "--datasource-provider", strings.ToLower(databaseProvider),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(err, errors.New(string(out)))
	}

	return nil
}

func (c *ConfigService) setupGorm(_ string) error {
	return nil
}

func (c *ConfigService) setupOrm(cfg Config) error {
	switch cfg.Orm.Name {
	case "Prisma":
		color.Green("Using Prisma ORM")

		err := c.setupPrisma(cfg.Orm.DatabaseProvider, cfg.Root)
		if err != nil {
			return err
		}
	case "Gorm":
		color.Green("Using GORM")

		err := c.setupGorm(cfg.Orm.DatabaseProvider)
		if err != nil {
			return err
		}
	default:
		return errors.New("orm is not supported")
	}

	return nil
}

func (c *ConfigService) provisionDatabase(cfg Config) error {
	outputFilePath, err := JoinURLsOrPaths(cfg.Root, "docker-compose.yaml")
	if err != nil {
		return err
	}

	color.Green("Creating %s with Docker Compose", cfg.Orm.DatabaseProvider)
	defer color.Green("Docker Compose file successfully generated [%s]", outputFilePath)

	err = GenerateSingleTmpl(GenerateSingleTmplArgs{
		TmplPath:   "docker-compose.yaml",
		OutputPath: outputFilePath,
		Values: map[string]interface{}{
			"ProjectName":      GetDirectoryName(),
			"DatabaseProvider": strings.ToLower(cfg.Orm.DatabaseProvider),
		},
		Funcs: map[string]any{"removeSigns": internals.RemoveNoneAlpha},
	})
	if err != nil {
		return err
	}

	databaseUrl := fmt.Sprintf(`"postgresql://user:password@localhost:5432/%s?schema=public"`, internals.RemoveNoneAlpha(cfg.ProjectName))
	err = internals.WriteEnvVar(".env", "DATABASE_URL", databaseUrl)
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

	err := os.Chdir(config.Root)
	if err != nil {
		return err
	}

	if args.ShouldProvideDatabase {
		err := c.provisionDatabase(config)
		if err != nil {
			return err
		}
	}

	err = c.setupOrm(config)
	if err != nil {
		return err
	}

	err = internals.WriteYaml("alchemy.yaml", config)
	if err != nil {
		return err
	}

	color.Green("✨ Alchemy config has been generated!")

	color.Green("🛠️  Updating Go dependencies ...")
	cmd := exec.Command("go", "mod", "tidy")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(err, errors.New(string(out)))
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
$ go-alchemy add Authentication // this will add all components from the authentication module

Or add a specific component
$ go-alchemy add Authentication.Login
	`)

	return nil
}

type AddArgs struct {
	Component string
	Root      string
}

// Adds single component to your project
//
//	`Authentication.Login` would be referring to only the login service
//	`Authentication` would add all available features in the authentication category
//
// The reference is case insenstive, we just like how the casing looks 🙂
func (c *ConfigService) Add(args AddArgs) error {
	err := os.Chdir(args.Root)
	if err != nil {
		return err
	}

	var (
		categoryId  string
		componentId string
	)

	if len(strings.Split(args.Component, ".")) > 1 {
		categoryId = strings.SplitN(args.Component, ".", 2)[0]
		componentId = strings.SplitN(args.Component, ".", 2)[1]
	} else {
		categoryId = args.Component
		componentId = "all"
	}

	categoryId = lo.Capitalize(categoryId)
	componentId = lo.Capitalize(componentId)

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

	if componentId == "All" {
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
