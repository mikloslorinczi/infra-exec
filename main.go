package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

var (
	command string
	outfile string
)

func main() {
	flag.Parse()
	splitCommand := strings.Split(command, " ")
	fmt.Println("\nCommand Executor\n")
	fmt.Println("Full Command:", command)
	fmt.Println("Base Command:", splitCommand[0])
	fmt.Println("File:", outfile)
	fmt.Println("Checking command ...")
	fmt.Println("Check command :", checkCommand(splitCommand[0]))
}

func init() {
	flag.StringVar(&command, "c", "", "Command to execute")
	flag.StringVar(&outfile, "f", "outfile", "Otput file")
}

func checkCommand(command string) bool {
	path, err := exec.LookPath(command)
	if err != nil {
		fmt.Println("Did not find command :", command)
		return false
	}
	fmt.Printf("The command found in %s'\n", path)
	return true
}
