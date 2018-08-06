package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	command string
	resFile string
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}
	flag.Parse()
	fmt.Println("\nCommand Executor")
	fmt.Println("--------------------------------------")
	fmt.Println("Full Command:", command)
	fmt.Println("Output File:", resFile)
	fmt.Println("--------------------------------------")
	execCommand(command)
}

func init() {
	flag.StringVar(&command, "c", "", "Command to execute")
	flag.StringVar(&resFile, "f", "resFile", "Otput File")
}

func printUsage() {
	fmt.Println("\nCommand Executor")
	fmt.Println()
	fmt.Printf("Usage: %s [option] [command or filename]\n\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	quit(1, "")
}

func execCommand(command string) {

	fmt.Println("Executing command ...")

	cmd := exec.Command("sh", "-c", command) // #nosec

	outFile, err := os.Create(resFile)
	if err != nil {
		quit(2, "Cannot create Output File")
	}
	defer func() {
		err = outFile.Close()
		if err != nil {
			quit(2, "Error closeing Output File")
		}
	}()

	writer := bufio.NewWriter(outFile)
	defer func() {
		err = writer.Flush()
		if err != nil {
			quit(3, "Error flushing file buffer")
		}
	}()

	cmd.Stdout = writer
	cmd.Stderr = writer

	err = cmd.Start()
	if err != nil {
		quit(3, "Error executing command")
	}

	err = cmd.Wait()
	if err != nil {
		if cmd.ProcessState.String() == "exit status 127" {
			quit(3, "Error: command not found")
		}
		quit(3, "Error during command execution, see resFile for details")
	}

	fmt.Println("end of execution ...")
	fmt.Println("output results of the command can be found in the resFile")
	fmt.Println()
}

func quit(status int, msg string) {
	fmt.Println(msg)
	os.Exit(status)
}
