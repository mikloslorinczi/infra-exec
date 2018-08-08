package executor_test

import (
	"bytes"
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

func (s *CommandExecutorTestSuite) TestParseCommand_empty_input() {
	name, args := executor.ParseCommand("")
	require.Equal(s.T(), "", name, "parseCommand should return empty name when an empty string is the input.")
	require.Equal(s.T(), []string{}, args, "parseCommand should return an empty array of string as args, when an empty string is the input.")
}

func (s *CommandExecutorTestSuite) TestExecCommand_empty_input() {
	testFile, err := executor.GetWriteFile("testFile")
	if err != nil {
		s.Error(err, "getWriteFile should create testFile")
	}
	defer testFile.Close()
	defer os.Remove("testFile")
	err = executor.ExecCommand(" ", testFile)
	if err != nil {
		s.Error(err, `./infra-exec -c " " should produce an empty outputFile and no error.`)
	}
}

func (s *CommandExecutorTestSuite) TestExecCommand_echo_cheese() {
	buf := new(bytes.Buffer)
	err := executor.ExecCommand("echo cheese", buf)
	if err != nil {
		s.Error(err, `./infra-exec -c "echo cheese" should produce outputFile with cheese in it, and no error.`)
	}
	require.Equal(s.T(), "cheese\n", buf.String(), "ExecCommand should create outputFile with \"cheese\\n\" in it. Got : %v", buf.String())
}

func TestCommandExecutor(t *testing.T) {
	suite.Run(t, &CommandExecutorTestSuite{})
}
