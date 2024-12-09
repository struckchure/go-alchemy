package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/struckchure/go-alchemy/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "alchemy",
	Short: "Alchemy CLI",
}

func main() {
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.RemoveCmd)

	rootCmd.PersistentFlags().StringP("root", "r", ".", "Project root")

	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err)
		return
	}
}
