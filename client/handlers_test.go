package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientHandlerTestSuite struct {
	suite.Suite
	*require.Assertions
}

func (s *ClientHandlerTestSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *ClientHandlerTestSuite) TestCompareTags_Empty() {
	match := compareTags("", "")
	require.Equal(s.T(), true, match, "compareTags should return ture when fed with two empty strings")
}

func (s *ClientHandlerTestSuite) TestCompareTags_Different() {
	match := compareTags("a", "b")
	require.Equal(s.T(), false, match, "compareTags should return false when fed with two completly different strings")
}

func (s *ClientHandlerTestSuite) TestCompareTags_Matching() {
	match := compareTags("cat dog penguin elephant shark", "bird dog worm cat penguin crocodile shark elephant")
	require.Equal(s.T(), true, match, "compareTags should return ture when the second string contains all tag from the first string")
}

func TestHandler(t *testing.T) {
	suite.Run(t, &ClientHandlerTestSuite{})
}
