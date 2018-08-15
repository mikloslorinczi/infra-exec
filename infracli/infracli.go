package infracli

import (
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

// ListTasks queries the infra-server for all the Tasks.
// It return the slice of Tasks and an optional error related to fetch and decode JSON.
func ListTasks() ([]common.Task, error) {
	tasksJSON, err := common.SendRequest("GET", common.APIURL+"/task/list", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get task list")
	}
	tasks, err := common.JSONToTasks(tasksJSON)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get task list")
	}
	return tasks, nil
}

// ReadCommand reads a command and its tags from stdin and
// converts them to a CommandObj, which will be returned along with an optional io error.
func ReadCommand() (common.CommandObj, error) {
	commandObj := common.CommandObj{}
	command, err := common.GetInput("Enter command :")
	if err != nil {
		return commandObj, errors.Wrap(err, "Cannot add new task")
	}
	tags, err := common.GetInput("Enter tags :")
	if err != nil {
		return commandObj, errors.Wrap(err, "Cannot add new task")
	}
	return common.CommandObj{Command: command, Tags: tags}, nil
}

// AddTask sends a CommandObj containing the command and its tags to the
// infra-server. It returns the servers's response message (which will be
// the newly created Task's ID) and on optional http/decode error.
func AddTask(commandObj common.CommandObj) (common.ResponseMsg, error) {
	responseMsg := common.ResponseMsg{}
	commandJSON, err := common.CommandToJSON(commandObj)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	responseJSON, err := common.SendRequest("POST", common.APIURL+"/task/add", commandJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	responseMsg, err = common.JSONToMsg(responseJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	return responseMsg, nil
}

// QueryTask queries the infra-server for one Task by its ID.
// It returns the found Task or an optional http/decode error.
func QueryTask(query string) (common.Task, error) {
	task := common.Task{}
	taskJSON, err := common.SendRequest("GET", common.APIURL+"/task/query/"+query, nil)
	if err != nil {
		return task, errors.Wrap(err, "Cannot query task")
	}
	task, err = common.JSONToTask(taskJSON)
	if err != nil {
		return task, errors.Wrap(err, "Cannot query task")
	}
	return task, nil
}
