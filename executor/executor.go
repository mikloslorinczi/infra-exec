package executor

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Command is parsed from the commandString,
// its the first part separated by space " " [0].
type Command string

// CommandArguments is an array of strings parsed from the commandString,
// separated by sapces " ", except for the first element [1:].
type CommandArguments []string

// ParseCommand takes the commandString and split it with the " " separator.
// The first element will be the Command itself, and the others will be its arguments.
func ParseCommand(commandString string) (Command, CommandArguments) {
	commandSlice := strings.Split(commandString, " ")
	command := Command(commandSlice[0])
	arguments := CommandArguments(commandSlice[1:])
	return command, arguments
}

// NewWriteFile accepts a filename and tries to open it for Writing.
// It returns a pointer to the file, and an optional error (if it fails to get the required file).
func NewWriteFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		err = errors.Wrapf(err, "Cannot open outputFile %v for write", path)
		file = nil
	}
	return file, err
}

// ExecCommand accepts executor's two own type Command and CommandArguments produced by ParseCommand,
// and an io.Writer where the result of the execution will be writen (both stdout and stderr).
// The function may return an error raised during the execution
func ExecCommand(command Command, commandArgs CommandArguments, w io.Writer) (err error) {

	writer := bufio.NewWriter(w)
	defer func() {
		err = writer.Flush()
		if err != nil {
			err = errors.Wrap(err, "Cannot flush buffer")
			log.Fatal(err)
		}
	}()

	cmd := exec.Command(string(command), commandArgs...) // #nosec
	cmd.Stdout = writer
	cmd.Stderr = writer

	err = cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "Cannot execute command %v", command)
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrapf(err, "Error during execution of command %v", command)
	}

	return

}
