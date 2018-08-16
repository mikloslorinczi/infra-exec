package main

import (
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

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
