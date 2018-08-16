package main

import (
	"fmt"
	"os"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

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

func addTask() {
	commandObj, err := readCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	response, err := addNewTask(commandObj)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("New Task added with the ID %v\n", response.Msg)
}

func listTasks() {
	tasks, err := common.GetTasks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, task := range tasks {
		fmt.Printf("\nID : %v\nNode : %v\nTags : %v\nStatus : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
	}
}

func queryTask() {
	task, err := getTaskByID(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Task\nID      : %v\nNode    : %v\nTags    : %v\nStatus  : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
}
