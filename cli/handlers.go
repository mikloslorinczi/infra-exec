package main

import (
	"fmt"
	"os"
)

func listTasks() {
	tasks, err := getTasks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, task := range tasks {
		fmt.Printf("\nID : %v\nNode : %v\nTags : %v\nStatus : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
	}
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

func queryTask() {
	task, err := getTaskByID(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Task\nID      : %v\nNode    : %v\nTags    : %v\nStatus  : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
}
