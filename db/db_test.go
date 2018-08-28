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

// var (
// 	Equal   = require.Equal
// 	NoError = require.NoError
// )

type DBTestSuite struct {
	suite.Suite
	*require.Assertions
	TestDB db.Controller
}

func (s *DBTestSuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.TestDB = db.NewJSONDB("testDB.json")
	if err := s.TestDB.Load(); err != nil {
		log.Fatalf("Error loading testDB.josn %v\n", err)
	}
}

func (s *DBTestSuite) TearDownTest() {
	if err := os.Remove("testDB.json"); err != nil {
		log.Fatalf("Cannot remove test testDB.json %v\n", err)
	}
}

func (s *DBTestSuite) TestQueryAll_EmptyDB() {
	tasks, found := s.TestDB.QueryAll()
	s.Equal([]common.Task(nil), tasks, "QueryAll shoud return an empty slice of Tasks")
	s.False(found, "Found should be false in the case of an empty DB")
}

func (s *DBTestSuite) TestAdd_Empty_Task() {
	id, err := s.TestDB.Add(common.Task{})
	s.NoError(err, "Add shouldn't return an error when fed with an empty Task")
	s.Equal(20, len(id), "jsonDB should create a 20 character-log XID")
}

func (s *DBTestSuite) TestAdd_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	s.NoError(err, "Add shouldn't return an error when fed with a valid Task")
	s.Equal(20, len(id), "jsonDB should create a 20 character-log XID")
}

func (s *DBTestSuite) TestQuery_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	s.NoError(err, "Add shouldn't return an error when fed with a valid Task")
	task, ok := s.TestDB.Query(id)
	s.True(ok, "Query should return OK true, if Task is found")
	s.Equal("echo cheese", task.Command, "Query should return the correct Task")
}

func (s *DBTestSuite) TestRemove_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	s.NoError(err, "Add shouldn't return an error when fed with a valid Task")
	ok, err2 := s.TestDB.Remove(id)
	s.NoError(err2, "Remove shuldn't return an error, when removing an existing Task")
	s.True(ok, "Remove should retutn OK ture, when successfully removed Task")
	_, ok = s.TestDB.Query(id)
	s.False(ok, "Query should return OK false, when searching for removed Task")
}

func (s *DBTestSuite) TestUpdate_Task() {
	id, err := s.TestDB.Add(common.Task{Command: "echo cheese", Tags: "tag1 tag2 tag3"})
	s.NoError(err, "Add shouldn't return an error when fed with a valid Task")
	newTask := common.Task{Node: "MyNode", Command: "tree", Tags: "tag5"}
	ok, err := s.TestDB.Update(id, newTask)
	s.NoError(err, "Update shouldn't return an error when fed with a valid ID & Task")
	s.True(ok, "Update should retutn OK ture, when successfully updated Task")
	task, ok := s.TestDB.Query(id)
	s.True(ok, "Query should return OK true, when searching for an updated Task")
	s.Equal("tag5", task.Tags, "Update shoud change the Task accordingly")
}

func TestDB(t *testing.T) {
	suite.Run(t, &DBTestSuite{})
}
