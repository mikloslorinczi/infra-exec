package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// RootCmd is the Cobra "entrypoint" icli
var RootCmd = &cobra.Command{
	Use:   "icli",
	Short: "Infra CLI is an infrastructure tool",
	Long: `
Infra CLI is an infrastructure tool for remote command execution

You can list, query and add Tasks for the Infra Server
Which distributes them to Infra Clients based on tag-matching.
Logfiles of executed tasks can be downloaded with the Infra CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			log.Fatalf("Error printing Usage %v\n", err)
		}
	},
}
