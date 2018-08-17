package main

import (
	"fmt"
	"os"
	"strings"

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

func addTask() (common.ResponseMsg, error) {
	response := common.ResponseMsg{}
	commandObj, err := readCommand()
	if err != nil {
		return response, errors.Wrap(err, "Error reading commoand")
	}
	response, err = addNewTask(commandObj)
	if err != nil {
		return response, errors.Wrap(err, "Error adding new task")
	}
	return response, nil
}

func listTasks() ([]common.Task, error) {
	tasks, err := common.GetTasks()
	if err != nil {
		return tasks, errors.Wrap(err, "Error getting tasks")
	}
	return tasks, nil
}

func queryTask(query string) (common.Task, error) {
	task, err := getTaskByID(query)
	if err != nil {
		return task, errors.Wrap(err, "Error querying task")
	}
	return task, nil
}

func getLog(ID string) {
	task, err := queryTask(logs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.Contains(task.Status, "log available") {
		err := downloadLog(task.ID)
		if err != nil {
			fmt.Printf("Error downloading logfile %v.log\n%v\n", task.ID, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully downloaded logfile logs/%v.log\n", task.ID)
	} else {
		fmt.Printf("Log not available task status %v\n", task.Status)
	}
}
