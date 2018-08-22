package db_test

import (
	"log"
	"os"
	"testing"

	"github.com/mikloslorinczi/infra-exec/common"

	"github.com/mikloslorinczi/infra-exec/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	*require.Assertions
	TestDB db.I
}

func (s *DBTestSuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.TestDB = db.ConnectJSONDB("testDB.json")
	err := s.TestDB.Load()
	if err != nil {
		log.Fatalf("Error loading testDB.josn %v\n", err)
	}
}

func (s *DBTestSuite) TearDownTest() {
	err := os.Remove("testDB.json")
	if err != nil {
		log.Fatalf("Cannot remove test testDB.json %v\n", err)
	}
}

func (s *DBTestSuite) TestQueryAll_EmptyDB() {
	tasks, found := s.TestDB.QueryAll()
	require.Equal(s.T(), []common.Task(nil), tasks, "QueryAll shoud return an empty slice of Tasks")
	require.Equal(s.T(), false, found, "Found should be false in the case of an empty DB")
}

func (s *DBTestSuite) TestAdd_Empty_Task() {
	id, err := s.TestDB.Add(common.Task{})
	require.NoError(s.T(), err, "Add shouldn't return an error when fed with an empty Task")
	require.Equal(s.T(), 20, len(id), "jsonDB should create a 20 character-log XID")
}

func (s *DBTestSuite) TestAdd_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	require.NoError(s.T(), err, "Add shouldn't return an error when fed with a valid Task")
	require.Equal(s.T(), 20, len(id), "jsonDB should create a 20 character-log XID")
}

func (s *DBTestSuite) TestQuery_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	require.NoError(s.T(), err, "Add shouldn't return an error when fed with a valid Task")
	task, ok := s.TestDB.Query(id)
	require.Equal(s.T(), true, ok, "Query should return OK true, if Task is found")
	require.Equal(s.T(), "echo cheese", task.Command, "Query should return the correct Task")
}

func (s *DBTestSuite) TestRemove_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	require.NoError(s.T(), err, "Add shouldn't return an error when fed with a valid Task")
	ok, err2 := s.TestDB.Remove(id)
	require.NoError(s.T(), err2, "Remove shuldn't return an error, when removing an existing Task")
	require.Equal(s.T(), true, ok, "Remove should retutn OK ture, when successfully removed Task")
	_, ok = s.TestDB.Query(id)
	require.Equal(s.T(), false, ok, "Query should return OK false, when searching for removed Task")
}

func (s *DBTestSuite) TestUpdate_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	require.NoError(s.T(), err, "Add shouldn't return an error when fed with a valid Task")
	newTask := common.Task{Node: "MyNode", Command: "tree", Tags: "tag5"}
	ok, err := s.TestDB.Update(id, newTask)
	require.NoError(s.T(), err, "Update shouldn't return an error when fed with a valid ID & Task")
	require.Equal(s.T(), true, ok, "Update should retutn OK ture, when successfully updated Task")
	task, ok := s.TestDB.Query(id)
	require.Equal(s.T(), true, ok, "Query should return OK true, when searching for an updated Task")
	require.Equal(s.T(), "tag5", task.Tags, "Update shoud change the Task accordingly")
}

func TestDB(t *testing.T) {
	suite.Run(t, &DBTestSuite{})
}
