package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func listTasks() error {
	commandJSON, err := commandToJSON("tasks", "")
	if err != nil {
		return errors.Wrap(err, "Cannot get task list")
	}
	tasksJSON, err := postRequest(apiURL+"/tasks/list", commandJSON)
	if err != nil {
		return errors.Wrap(err, "Cannot get task list")
	}
	tasks, err := bodyToTasks(tasksJSON)
	if err != nil {
		return errors.Wrap(err, "Cannot get task list")
	}
	fmt.Println("List of Tasks")
	for _, task := range tasks {
		fmt.Printf("\nID : %v\nNode : %v\nTags : %v\nStatus : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
	}
	fmt.Println()
	return nil
}

func addTask() error {
	command, err := getInput("Enter command :")
	if err != nil {
		return errors.Wrap(err, "Cannot add new task")
	}
	tags, err := getInput("Enter tags :")
	if err != nil {
		return errors.Wrap(err, "Cannot add new task")
	}
	fmt.Printf("Command : %v\nTags : %v\n", command, tags)
	commandJSON, err := commandToJSON(command, tags)
	if err != nil {
		return errors.Wrap(err, "Cannot add new task")
	}
	responseJSON, err := postRequest(apiURL+"/tasks/add", commandJSON)
	if err != nil {
		return errors.Wrap(err, "Cannot add new task")
	}
	response, err := bodyToMsg(responseJSON)
	if err != nil {
		return errors.Wrap(err, "Cannot add new task")
	}
	fmt.Printf("Succesfully added new task with the %v\n", response.Msg)
	return nil
}
