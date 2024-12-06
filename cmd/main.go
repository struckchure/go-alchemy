package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/struckchure/go-alchemy/components"
	"github.com/struckchure/go-alchemy/orms"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = time.Now().Format(time.DateTime)
)

var rootCmd = &cobra.Command{
	Use:   "alchemy",
	Short: "Alchemy CLI",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s\n", version, commit, buildDate)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new alchemy project",
	Run: func(cmd *cobra.Command, args []string) {
		var root string
		var orm string
		var databaseProvider string
		var databaseUrl string
		var provisionDatabase string

		err := survey.AskOne(&survey.Input{Message: "Provide alchemy component root: ", Default: "."}, &root)
		if err != nil {
			color.Red("%s", err)
			return
		}

		err = survey.AskOne(&survey.Select{Message: "Choose ORM: ", Options: orms.OrmOptions}, &orm)
		if err != nil {
			color.Red("%s", err)
			return
		}

		err = survey.AskOne(
			&survey.Select{
				Message: "Choose Database Provider: ",
				Options: orms.OrmMappings[orm],
			}, &databaseProvider,
		)
		if err != nil {
			color.Red("%s", err)
			return
		}

		err = survey.AskOne(
			&survey.Select{
				Message: "Provision Database with Docker Compose: ",
				Options: []string{"Yes", "No"},
				Default: "Yes",
			},
			&provisionDatabase,
		)
		if err != nil {
			color.Red("%s", err)
			return
		}

		if strings.ToLower(provisionDatabase) == "no" {
			err = survey.AskOne(&survey.Input{Message: "Provide database url: "}, &databaseUrl)
			if err != nil {
				color.Red("%s", err)
				return
			}
		}

		err = components.NewConfigService().Init(components.InitArgs{
			Root:                  root,
			Orm:                   orm,
			ShouldProvideDatabase: strings.ToLower(provisionDatabase) == "yes",
			DatabaseUrl:           databaseUrl,
			DatabaseProvider:      databaseProvider,
		})
		if err != nil {
			color.Red("%s", err)
			return
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new component",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			categoryId  string
			componentId string
		)

		if len(args) > 0 {
			id := args[0]
			idSplit := strings.SplitN(id, ".", 2)
			categoryId = strings.SplitN(id, ".", 2)[0]

			if len(idSplit) > 1 {
				componentId = lo.Must(lo.Last(strings.SplitN(id, ".", 2)))
			}
		}

		if categoryId == "" {
			err := survey.AskOne(&survey.Select{
				Message: "Component Category",
				Options: components.ComponentCategoryOptions,
			}, &categoryId)
			if err != nil {
				color.Red("%s", err)
				return
			}
		}

		if componentId == "" {
			err := survey.AskOne(&survey.Select{
				Message: "Select Components",
				Options: components.ComponentMapping[categoryId],
			}, &componentId)
			if err != nil {
				color.Red("%s", err)
				return
			}
		}

		err := components.NewConfigService().Add(fmt.Sprintf("%s.%s", categoryId, componentId))
		if err != nil {
			color.Red("%s", err)
			return
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Add new component",
	Args:  cobra.ExactArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {},
}

func main() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // TODO: use logger
		return
	}

	// client := db.NewClient()
	// if err := client.Prisma.Connect(); err != nil {
	// 	panic(err)
	// }

	// defer func() {
	// 	if err := client.Prisma.Disconnect(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// registerResult, err := services.NewAuthenticationService(client).Register(services.RegisterArgs{
	// 	Email:    "john@doe.com",
	// 	Password: "1224321434",
	// })
	// if err != nil {
	// 	color.Red("%s", err)
	// 	return
	// }
	// color.Green(registerResult.User.Id)
	// color.Green(registerResult.User.FirstName)
	// color.Green(registerResult.User.LastName)
	// color.Green(registerResult.User.Email)

	// loginResult, err := services.NewAuthenticationService(client).Login(services.LoginArgs{
	// 	Email:    "john@doe.com",
	// 	Password: "1224321434",
	// })
	// if err != nil {
	// 	color.Red("%s", err)
	// 	return
	// }

	// color.Green(loginResult.User.Id)
	// color.Green(loginResult.User.FirstName)
	// color.Green(loginResult.User.LastName)
	// color.Green(loginResult.User.Email)
}
