package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mikloslorinczi/infra-exec/executor"
)

var (
	commandString  string
	outputFileName string
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	outputFile, err := executor.NewWriteFile(outputFileName)
	if err != nil {
		fmt.Printf("Error opening outputFile %v\nError :\n%v\n", outputFileName, err)
		os.Exit(1)
	}

	defer func() {
		err = outputFile.Close()
		if err != nil {
			log.Fatalf("Cannot close outputFile %v\nError :\n%v\n", outputFileName, err)
		}
	}()

	command, commandArgs := executor.ParseCommand(commandString)
	if err := executor.ExecCommand(command, commandArgs, outputFile); err != nil {
		fmt.Printf("Cannot execute command (%v)\nError :\n%v\n", commandString, err)
		os.Exit(1)
	}

}

func init() {
	flag.StringVar(&commandString, "c", "", "Command to execute")
	flag.StringVar(&outputFileName, "f", "outputFile", "Otput File")
}

func printUsage() {
	fmt.Println("\nCommand Executor")
	fmt.Println()
	fmt.Printf("Usage: %s [option] [command or filename]", os.Args[0])
	fmt.Println("\n\nOptions:")
	flag.PrintDefaults()
}
