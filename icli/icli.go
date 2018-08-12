package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	listTasks bool
	addTask   string
	logsTask  string
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}
	flag.Parse()
}

func init() {
	flag.BoolVar(&listTasks, "l", false, "List tasks")
	flag.StringVar(&addTask, "a", "", "Add new task")
	flag.StringVar(&logsTask, "log", "", "Require logs of given task")
}

func printUsage() {
	fmt.Println("\nInfra CLI")
	fmt.Println()
	fmt.Println("Execute tasks remotely")

	fmt.Printf("Usage: %s [option] [command or taskID]", os.Args[0])
	fmt.Println("\n\nOptions:")
	flag.PrintDefaults()
}
