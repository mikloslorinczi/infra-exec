package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientCmdTestSuite struct {
	suite.Suite
	*require.Assertions
}

func (s *ClientCmdTestSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *ClientCmdTestSuite) TestCompareTags_Empty() {
	match := compareTags("", "")
	require.Equal(s.T(), true, match, "compareTags should return ture when fed with two empty strings")
}

func (s *ClientCmdTestSuite) TestCompareTags_Different() {
	match := compareTags("a", "b")
	require.Equal(s.T(), false, match, "compareTags should return false when fed with two completly different strings")
}

func (s *ClientCmdTestSuite) TestCompareTags_Matching() {
	match := compareTags("cat dog penguin elephant shark", "bird dog worm cat penguin crocodile shark elephant")
	require.Equal(s.T(), true, match, "compareTags should return ture when the second string contains all tag from the first string")
}

func TestHandler(t *testing.T) {
	suite.Run(t, &ClientCmdTestSuite{})
}
