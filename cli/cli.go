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

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	if err := common.SetAdminPass(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Parse()

	if list {
		listTasks()
	}

	if add {
		addTask()
	}

	if query != "" {
		queryTask()
	}

}

func init() {
	flag.BoolVar(&list, "l", false, "List tasks")
	flag.BoolVar(&add, "a", false, "Add new task")
	flag.StringVar(&query, "q", "", "Query task by ID")
	flag.StringVar(&logs, "log", "", "Require logs of task by ID")
	flag.StringVar(&common.APIURL, "u", "http://localhost:7474/api", "URL address of the api")
}

func printUsage() {
	fmt.Println("\nInfra CLI")
	fmt.Println()
	fmt.Println("Execute tasks remotely")

	fmt.Printf("Usage: %s [option] [command or taskID]", os.Args[0])
	fmt.Println("\n\nOptions:")
	flag.PrintDefaults()
}
