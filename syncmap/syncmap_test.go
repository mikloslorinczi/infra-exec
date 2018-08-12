package syncmap_test

import (
	"testing"

	"github.com/mikloslorinczi/infra-exec/syncmap"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncMapTestSuite struct {
	suite.Suite
	*require.Assertions
}

func (s *SyncMapTestSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *SyncMapTestSuite) TestGetMap() {
	var testMap = syncmap.NewSafeMap()
	emptyMap := make(map[string]string)
	myMap := testMap.GetMap()
	require.Equal(s.T(), emptyMap, myMap, "Excepted an empty string-map, got %v", myMap)
}

func (s *SyncMapTestSuite) TestSetKey() {
	var testMap = syncmap.NewSafeMap()
	testMap.SetKey("No.1", "Zeus")
	value := testMap.GetMap()["No.1"]
	require.Equal(s.T(), "Zeus", value, "Excepted \"Zeus\", got %v", value)
}

func (s *SyncMapTestSuite) TestGetKey_Nonexisting() {
	var testMap = syncmap.NewSafeMap()
	value := testMap.GetKey("No.2")
	require.Equal(s.T(), "", value, "Excepted \"\", got %v", value)
}

func (s *SyncMapTestSuite) TestGetKey_With_Preset() {
	var testMap = syncmap.NewSafeMap()
	testMap.SetKey("No.2", "Odin")
	value := testMap.GetKey("No.2")
	require.Equal(s.T(), "Odin", value, "Excepted \"Odin\", got %v", value)
}

func (s *SyncMapTestSuite) TestDeleteKey() {
	var testMap = syncmap.NewSafeMap()
	testMap.SetKey("No.2", "Odin")
	testMap.DeleteKey("No.2")
	value := testMap.GetKey("No.2")
	require.Equal(s.T(), "", value, "Excepted \"\", got %v", value)
}

func TestSyncMap(t *testing.T) {
	suite.Run(t, &SyncMapTestSuite{})
}
