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

		cmd := exec.Command(command)

		outFile, err := os.Create(resFile)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		errorFile, err := os.Create(errFile)
		if err != nil {
			panic(err)
		}
		defer errorFile.Close()

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}

		outWriter := bufio.NewWriter(outFile)
		defer outWriter.Flush()

		errWriter := bufio.NewWriter(errorFile)
		defer errWriter.Flush()

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		go io.Copy(outWriter, stdoutPipe)
		go io.Copy(errWriter, stderrPipe)

		cmd.Wait()

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
