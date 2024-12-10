package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "go-alchemy",
	Short: "Golang Alchemy CLI",
}

func init() {
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(RemoveCmd)

	RootCmd.PersistentFlags().StringP("root", "r", ".", "Project root")
}
