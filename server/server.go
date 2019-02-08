package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/mikloslorinczi/infra-exec/server/cmd"
)

func main() {
	// Run Server
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("Error during execution : %v\n", err)
	}
}

func init() {
	// Create logs/ if not exists
	if err := common.CheckLogFolder(); err != nil {
		fmt.Printf("Error loading configFile cli.yaml %v\n", err)
	}
	// Load config (defaults < server.yaml < ENV)
	if err := common.ReadConfig("./", "server", map[string]interface{}{
		"PORT":     8080,
		"APITOKEN": "",
	}); err != nil {
		log.Fatalf("Cannot load configfile server.yaml %v\n", err)
	}
	// Setup global flags
	cmd.RootCmd.PersistentFlags().StringP("apiToken", "t", "", "API Token")
	cmd.RootCmd.PersistentFlags().StringP("PORT", "p", "", "Server PORT")
	if err := viper.BindPFlag("apiToken", cmd.RootCmd.PersistentFlags().Lookup("apiToken")); err != nil {
		fmt.Printf("Cannot bind flag apiToken %v\n", err)
	}
	if err := viper.BindPFlag("PORT", cmd.RootCmd.PersistentFlags().Lookup("PORT")); err != nil {
		fmt.Printf("Cannot bind flag PORT %v\n", err)
	}
}
