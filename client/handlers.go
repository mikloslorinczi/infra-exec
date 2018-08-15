package main

import (
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

func fetchTasks() ([]common.Task, error) {
	tasks := []common.Task{}
	tasksJSON, err := common.SendRequest("GET", common.APIURL+"/task/list", nil)
	if err != nil {
		err = errors.Wrap(err, "Cannot fetch Tasks")
		return tasks, err
	}
	tasks, err = common.JSONToTasks(tasksJSON)
	if err != nil {
		err = errors.Wrap(err, "Cannot fetch Tasks")
		return tasks, err
	}
	return tasks, nil
}
