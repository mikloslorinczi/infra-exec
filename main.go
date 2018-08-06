package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	command string
	resFile string
	errFile string
)

func main() {
	flag.Parse()
	splitCommand := strings.Split(command, " ")
	fmt.Println("\nCommand Executor")
	fmt.Println("Full Command:", command)
	fmt.Println("Base Command:", splitCommand[0])
	fmt.Println("Output File:", resFile)
	fmt.Println("Error File:", errFile)
	fmt.Println("Checking command ...")
	fmt.Println("Check command :", checkCommand(splitCommand[0]))

	if checkCommand(splitCommand[0]) {

		fmt.Println("Executing command ...")

		cmd := exec.Command(command) // #nosec

		outFile, err := os.Create(resFile)
		if err != nil {
			panic(err)
		}
		defer func() {
			err = outFile.Close()
			if err != nil {
				panic(err)
			}
		}()

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		errorFile, err := os.Create(errFile)
		if err != nil {
			panic(err)
		}
		defer func() {
			err = errorFile.Close()
			if err != nil {
				panic(err)
			}
		}()

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}

		outWriter := bufio.NewWriter(outFile)
		defer func() {
			err = outWriter.Flush()
			if err != nil {
				panic(err)
			}
		}()

		errWriter := bufio.NewWriter(errorFile)
		defer func() {
			err = errWriter.Flush()
			if err != nil {
				panic(err)
			}
		}()

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		go func() {
			_, writeErr := io.Copy(outWriter, stdoutPipe)
			if writeErr != nil {
				panic(writeErr)
			}
		}()

		go func() {
			_, writeErr := io.Copy(errWriter, stderrPipe)
			if writeErr != nil {
				panic(writeErr)
			}
		}()

		err = cmd.Wait()
		if err != nil {
			panic(err)
		}

		fmt.Println("end of execution ...")

	}
}

func init() {
	flag.StringVar(&command, "c", "", "Command to execute")
	flag.StringVar(&resFile, "f", "resFile", "Otput File")
	flag.StringVar(&errFile, "e", "errFile", "Error File")
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
