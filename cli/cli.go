package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

var (
	list  bool
	add   bool
	logs  string
	query string
)

// This function has complexity 13 according to gocyclo, despite it is "perfectly" flat
func main() { // nolint: gocyclo

	ok, status, err := initConfig()
	if !ok || err != nil {
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(status)
	}

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
		err := getLog(logs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Succesfully downloaded logs/%v.log\n", logs)
		}
	}

}

func init() {
	flag.BoolVar(&list, "l", false, "List tasks")
	flag.BoolVar(&add, "a", false, "Add new task")
	flag.StringVar(&query, "q", "", "Query task by ID")
	flag.StringVar(&logs, "log", "", "Require logs of task by ID")
	flag.StringVar(&common.AdminPass, "pass", "", "Admin Password")
	flag.StringVar(&common.APIURL, "u", "http://localhost:7474/api", "URL address of the api")
}

func initConfig() (bool, int, error) {
	if len(os.Args) < 2 {
		printUsage()
		return false, 0, nil
	}
	common.AdminPass = os.Getenv("ADMIN_PASSWORD")
	flag.Parse()
	if common.AdminPass == "" {
		fmt.Printf("\nNo ADMIN_PASSWORD found in the environment / flag\n")
		input, err := common.GetInput("Admin password : ")
		if err != nil {
			return false, 1, errors.Wrap(err, "Error reading input")
		}
		common.AdminPass = input
		fmt.Println()
	}
	return true, 0, nil
}

func printUsage() {
	fmt.Println("\nInfra CLI")
	fmt.Println()
	fmt.Println("User interface for the Infra Server")
	fmt.Println("You can list, query and add new Tasks.")
	fmt.Println("Also download the logs of already executed ones.")
	fmt.Println()
	fmt.Printf("Usage: %s [option] [command or taskID]", os.Args[0])
	fmt.Println("\n\nOptions:")
	flag.PrintDefaults()
}
