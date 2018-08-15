package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// GetInput prints the message to stdout and reads user input from stdin.
// The return values are the inputString and an optional io error.
func GetInput(msg string) (string, error) {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		err = errors.Wrapf(err, "Input error")
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}

// SetAdminPass sets the AdminPass variable. First it tries to get it from the
// environment. If not found user input will be required on stdin.
// May return an IO error.
func SetAdminPass() error {
	AdminPass = os.Getenv("ADMIN_PASSWORD")
	if AdminPass == "" {
		fmt.Printf("\nNo ADMIN_PASSWORD found in the environment\n")
		input, err := GetInput("Admin password : ")
		if err != nil {
			return err
		}
		AdminPass = input
		fmt.Println()
	}
	return nil
}
