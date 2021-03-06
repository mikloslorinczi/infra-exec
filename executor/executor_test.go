package executor_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/mikloslorinczi/infra-exec/executor"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommandExecutorTestSuite struct {
	suite.Suite
	*require.Assertions
}

func (s *CommandExecutorTestSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *CommandExecutorTestSuite) TestParseCommand_Empty_Input() {
	command, commandArgs := executor.ParseCommand("")
	s.Equal(executor.Command(""), command, "parseCommand should return empty name when an empty string is the input.")
	s.Equal(executor.CommandArguments([]string{}), commandArgs, "parseCommand should return an empty array of string as args, when an empty string is the input.")
}

func (s *CommandExecutorTestSuite) TestExecCommand_Empty_Input() {
	testFile, err := executor.NewWriteFile("testFile")
	defer func() {
		if fileErr := testFile.Close(); fileErr != nil {
			log.Fatalf("Cannot close testFile %v\n", fileErr)
		}
	}()
	defer func() {
		if remErr := os.Remove("testFile"); remErr != nil {
			log.Fatalf("Cannot remove TestFile %v\n", remErr)
		}
	}()
	s.NoError(err, "getWriteFile should create testFile")
	err = executor.ExecCommand(" ", []string{}, testFile)
	s.NoError(err, `./infra-exec -c " " should produce an empty outputFile and no error.`)
}

func (s *CommandExecutorTestSuite) TestExecCommand_Echo_Cheese() {
	buf := new(bytes.Buffer)
	err := executor.ExecCommand("echo", []string{"cheese"}, buf)
	s.NoError(err, `./infra-exec -c "echo cheese" should produce outputFile with cheese in it, and no error.`)
	s.Equal("cheese\n", buf.String(), "ExecCommand should create outputFile with \"cheese\\n\" in it. Got : %v", buf.String())
}

func TestCommandExecutor(t *testing.T) {
	suite.Run(t, &CommandExecutorTestSuite{})
}
