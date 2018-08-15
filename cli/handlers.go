package main

import (
	"fmt"
	"os"

	"github.com/mikloslorinczi/infra-exec/infracli"
)

func listTasks() {
	tasks, err := infracli.ListTasks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, task := range tasks {
		fmt.Printf("\nID : %v\nNode : %v\nTags : %v\nStatus : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
	}
}

func addTask() {
	commandObj, err := infracli.ReadCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	response, err := infracli.AddTask(commandObj)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("New Task added with the ID %v\n", response.Msg)
}

func queryTask() {
	task, err := infracli.QueryTask(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Task\nID      : %v\nNode    : %v\nTags    : %v\nStatus  : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
}
