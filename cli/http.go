package main

import (
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

func getTasks() ([]common.Task, error) {
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

func readCommand() (common.CommandObj, error) {
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

func addNewTask(commandObj common.CommandObj) (common.ResponseMsg, error) {
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

func getTaskByID(query string) (common.Task, error) {
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
