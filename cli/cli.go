package main

import (
	"fmt"
	"log"

	"github.com/mikloslorinczi/infra-exec/cli/cmd"
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/spf13/viper"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("Error during execution : %v\n", err)
	}
}

func init() {
	// Create logs/ if not exists
	if err := common.CheckLogFolder(); err != nil {
		fmt.Printf("Error loading configFile cli.yaml %v\n", err)
	}
	// Load config from cli.yaml and ENV
	if err := common.ReadConfig("./", "cli", map[string]interface{}{
		"apiToken": "",
		"apiURL":   "",
	}); err != nil {
		fmt.Printf("Cannot set configuration %v\n", err)
	}
	// Setup global flags
	cmd.RootCmd.PersistentFlags().StringP("apiToken", "t", "", "API Token")
	cmd.RootCmd.PersistentFlags().StringP("apiURL", "u", "", "API URL")
	if err := viper.BindPFlag("apiToken", cmd.RootCmd.PersistentFlags().Lookup("apiToken")); err != nil {
		fmt.Printf("Cannot bind flag apiToken %v\n", err)
	}
	if err := viper.BindPFlag("apiURL", cmd.RootCmd.PersistentFlags().Lookup("apiURL")); err != nil {
		fmt.Printf("Cannot bind flag apiURL %v\n", err)
	}
}
