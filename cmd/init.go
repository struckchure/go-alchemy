package main

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/struckchure/go-alchemy/components"
	"github.com/struckchure/go-alchemy/orms"
)

var InitCmd = &cobra.Command{
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
