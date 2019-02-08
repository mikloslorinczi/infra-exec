package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/mikloslorinczi/infra-exec/common"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <icli task add>",
	Short: "Add new Task to the Infra Server",
	Run: func(cmd *cobra.Command, args []string) {
		response, err := addTask()
		if err != nil {
			fmt.Println(response.Msg)
			log.Fatalf("Error adding new Task : %v\n", err)
		}
		fmt.Printf("New Task added to the Infra Server with the ID %v\n", response.Msg)
	},
}

func init() {
	taskCmd.AddCommand(addCmd)
}

func readCommand() (common.CommandObj, error) {
	commandObj := common.CommandObj{}
	command, err := common.GetInput("Enter command :")
	if err != nil {
		return commandObj, errors.Wrap(err, "Cannot add new task")
	}
	tags, err := common.GetInput("Enter tags :")
	if err != nil {
		return commandObj, errors.Wrap(err, "Cannot add new task")
	}
	return common.CommandObj{Command: command, Tags: tags}, nil
}

func addTask() (common.ResponseMsg, error) {
	var (
		responseMsg common.ResponseMsg
		commandJSON []byte
	)
	commandObj, err := readCommand()
	if err != nil {
		return responseMsg, errors.Wrap(err, "Error reading commoand")
	}
	if encodeErr := common.ToJSON(&commandObj, &commandJSON); encodeErr != nil {
		return responseMsg, errors.Wrap(encodeErr, "Cannot add new task")
	}
	responseJSON, err := common.SendRequest("POST", viper.GetString("apiURL")+"/tasks", commandJSON, common.HTTPHeaderJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	if err := common.FromJSON(&responseMsg, responseJSON); err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	return responseMsg, nil
}
