package common

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// LoadEnv accepts an envfile (with lines like: key=value) as parameter,
// and tries to load it to the environment. May return an io error.
func LoadEnv(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return errors.Wrapf(err, "Cannot open envFile %v\n", path)
	}
	defer func() {
		err2 := file.Close()
		if err2 != nil {
			log.Fatalf("Error closing envFile %v\n%v\n", path, err2)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) == 2 {
			err3 := os.Setenv(line[0], line[1])
			if err3 != nil {
				return errors.Wrapf(err3, "Error reading env from envFile %v\n", path)
			}
		}
	}
	if err4 := scanner.Err(); err4 != nil {
		log.Fatalf("Error reading env from envFile %v\n%v\n", path, err4)
	}

	return nil
}
