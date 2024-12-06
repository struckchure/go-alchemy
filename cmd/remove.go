package main

import "github.com/spf13/cobra"

var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Add new component",
	Args:  cobra.ExactArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {},
}
