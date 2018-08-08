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

// ParseCommand takes the command string and split it with the " " separator.
// The first element will be the basecommand and the others will be the args.
func ParseCommand(command string) (string, []string) {
	commandSlice := strings.Split(command, " ")
	name := commandSlice[0]
	args := commandSlice[1:]
	return name, args
}

// GetWriteFile accepts a filename and tries to open it for Writing.
// It returns a pointer to the file, and an optional error (if it fails to get the required file).
func GetWriteFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		err = errors.Wrapf(err, "Cannot open outputFile %v for write", fileName)
	}
	return file, err
}

// ExecCommand accepts a string as a command (it will be parsed to basecommand and args)
// and an io.Writer where the result of the execution will be writen (both stdout and stderr).
// The function may return an error raised during the execution
func ExecCommand(command string, w io.Writer) (err error) {

	writer := bufio.NewWriter(w)
	defer func() {
		err = writer.Flush()
		if err != nil {
			err = errors.Wrap(err, "Cannot flush buffer")
			log.Fatal(err)
		}
	}()

	baseCommand, args := ParseCommand(command)
	cmd := exec.Command(baseCommand, args...) // #nosec
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
