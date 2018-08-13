package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	list    bool
	add     bool
	logs    string
	envPass string
	apiURL  string
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	if err := setEnvPass(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Parse()

	if list {
		err := listTasks()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if add {
		err := addTask()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	flag.BoolVar(&list, "l", false, "List tasks")
	flag.BoolVar(&add, "a", false, "Add new task")
	flag.StringVar(&logs, "log", "", "Require logs of given task")
	flag.StringVar(&apiURL, "u", "http://localhost:7474/api", "API URL")
}

func printUsage() {
	fmt.Println("\nInfra CLI")
	fmt.Println()
	fmt.Println("Execute tasks remotely")

	fmt.Printf("Usage: %s [option] [command or taskID]", os.Args[0])
	fmt.Println("\n\nOptions:")
	flag.PrintDefaults()
}

func setEnvPass() error {
	envPass = os.Getenv("ADMIN_PASSWORD")
	if envPass == "" {
		fmt.Println("\nNo ADMIN_PASSWORD found in the environment\n")
		input, err := getInput("Admin password : ")
		if err != nil {
			return err
		}
		envPass = input
		fmt.Println()
	}
	return nil
}

func getInput(msg string) (string, error) {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		err = errors.Wrapf(err, "Input error")
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}
