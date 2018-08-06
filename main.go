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

	outFile, err := os.Create(resFile)
	if err != nil {
		quit(1, "Cannot open Output File\n")
	}
	defer func() {
		err = outFile.Close()
		if err != nil {
			quit(1, "Error closeing Output File")
		}
	}()

	fmt.Println("\nCommand Executor")
	fmt.Println("--------------------------------------")
	fmt.Println("Full Command:", command)
	fmt.Println("Output File:", resFile)
	fmt.Println("--------------------------------------")
	status, msg := execCommand(command, outFile)
	quit(status, msg)
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
	quit(1, "")
}

func execCommand(command string, outFile *os.File) (status int, msg string) {

	cmd := exec.Command("sh", "-c", command) // #nosec

	writer := bufio.NewWriter(outFile)
	defer func() {
		err := writer.Flush()
		if err != nil {
			status = 1
			msg = "Error flushing the buffer"
			return
		}
	}()

	cmd.Stdout = writer
	cmd.Stderr = writer

	err := cmd.Start()
	if err != nil {
		status = 1
		msg = "Error starting the execution"
		return
	}

	err = cmd.Wait()
	if err != nil {
		status = 1
		msg = "Error during execution :" + cmd.ProcessState.String()
		return
	}

	status = 0
	msg = "Execution completed"
	return

}

func quit(status int, msg string) {
	fmt.Println(msg)
	os.Exit(status)
}
