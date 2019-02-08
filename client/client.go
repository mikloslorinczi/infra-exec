package main

import (
	"fmt"
	"log"

	"github.com/mikloslorinczi/infra-exec/client/cmd"
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
		fmt.Printf("Error loading client.yaml %v\n", err)
	}
	// Load config from client.yaml and ENV
	if err := common.ReadConfig("./", "client", map[string]interface{}{
		"apiToken": "",
		"apiURL":   "",
		"nodeName": "",
		"nodeTags": "",
		"period":   60,
	}); err != nil {
		fmt.Printf("Cannot set configuration %v\n", err)
	}
	// Setup flags
	cmd.RootCmd.PersistentFlags().StringP("apiToken", "t", "", "API Token")
	cmd.RootCmd.PersistentFlags().StringP("apiURL", "u", "", "API URL")
	cmd.RootCmd.PersistentFlags().StringP("nodeName", "n", "", "Node Name")
	cmd.RootCmd.PersistentFlags().StringP("nodeTags", "s", "", "Node Tags")
	cmd.RootCmd.PersistentFlags().IntP("period", "p", 60, "Node poll period")
	if err := viper.BindPFlag("apiToken", cmd.RootCmd.PersistentFlags().Lookup("apiToken")); err != nil {
		fmt.Printf("Cannot bind flag apiToken %v\n", err)
	}
	if err := viper.BindPFlag("apiURL", cmd.RootCmd.PersistentFlags().Lookup("apiURL")); err != nil {
		fmt.Printf("Cannot bind flag apiURL %v\n", err)
	}
	if err := viper.BindPFlag("nodeName", cmd.RootCmd.PersistentFlags().Lookup("nodeName")); err != nil {
		fmt.Printf("Cannot bind flag nodeName %v\n", err)
	}
	if err := viper.BindPFlag("nodeTags", cmd.RootCmd.PersistentFlags().Lookup("nodeTags")); err != nil {
		fmt.Printf("Cannot bind flag nodeTags %v\n", err)
	}
	if err := viper.BindPFlag("period", cmd.RootCmd.PersistentFlags().Lookup("period")); err != nil {
		fmt.Printf("Cannot bind flag period %v\n", err)
	}
}
