package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikloslorinczi/infra-exec/common"
)

var (
	list  bool
	add   bool
	logs  string
	query string
)

func main() {

	initCLI()

	flag.Parse()

	if list {
		tasks, err := listTasks()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, task := range tasks {
			fmt.Printf("\nID : %v\nNode : %v\nTags : %v\nStatus : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
		}
	}

	if add {
		response, err := addTask()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("New task added with the ID %v\n", response.Msg)
	}

	if query != "" {
		task, err := queryTask(query)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Task\nID      : %v\nNode    : %v\nTags    : %v\nStatus  : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
	}

	if logs != "" {
		getLog(logs)
	}

}

func init() {
	flag.BoolVar(&list, "l", false, "List tasks")
	flag.BoolVar(&add, "a", false, "Add new task")
	flag.StringVar(&query, "q", "", "Query task by ID")
	flag.StringVar(&logs, "log", "", "Require logs of task by ID")
	flag.StringVar(&common.APIURL, "u", "http://localhost:7474/api", "URL address of the api")
}

func initCLI() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}
	if err := common.SetAdminPass(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("\nInfra CLI")
	fmt.Println()
	fmt.Println("Execute tasks remotely")
	fmt.Println()
	fmt.Printf("Usage: %s [option] [command or taskID]", os.Args[0])
	fmt.Println("\n\nOptions:")
	flag.PrintDefaults()
}
