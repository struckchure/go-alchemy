package main

import (
	"github.com/fatih/color"

	"github.com/struckchure/go-alchemy/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		color.Red("%s", err)
		return
	}
}
