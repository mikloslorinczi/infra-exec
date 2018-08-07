package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	command        string
	outputFileName string
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
	}

	flag.Parse()

	outputFile, err := getWriteFile(outputFileName)
	if err != nil {
		quit(1, err)
	}

	defer func() {
		err = outputFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = execCommand(command, outputFile)
	if err != nil {
		fmt.Printf("Error during execution %v", err)
		quit(1, err)
	}

}

func init() {
	flag.StringVar(&command, "c", "", "Command to execute")
	flag.StringVar(&outputFileName, "f", "outputFile", "Otput File")
}

func printUsage() {
	fmt.Println("\nCommand Executor")
	fmt.Println()
	fmt.Printf("Usage: %s [option] [command or filename]\n\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	quit(1, nil)
}

func getWriteFile(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func parseCommand(command string) (string, []string) {
	commandSlice := strings.Split(command, " ")
	name := commandSlice[0]
	args := []string{}
	if len(commandSlice) > 1 {
		args = commandSlice[1:]
	}
	return name, args
}

func execCommand(command string, outFile io.Writer) (err error) {

	writer := bufio.NewWriter(outFile)
	defer func() {
		err = writer.Flush()
		if err != nil {
			log.Fatal(err)
		}
	}()

	basecommand, args := parseCommand(command)
	cmd := exec.Command(basecommand, args...) // #nosec
	cmd.Stdout = writer
	cmd.Stderr = writer

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return

}

func quit(status int, err error) {
	fmt.Println(err.Error())
	os.Exit(status)
}
