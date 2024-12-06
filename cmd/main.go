package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alchemy",
	Short: "Alchemy CLI",
}

func main() {
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(AddCmd)
	rootCmd.AddCommand(RemoveCmd)

	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err)
		return
	}
}
