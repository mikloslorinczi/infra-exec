package cmd

import (
	"strings"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CompareTags checks if all the tags of the Task matches with the Node's tags.
func compareTags(taskTags, nodeTags string) bool {
	taskSplit := strings.Split(taskTags, " ")
	nodeSplit := strings.Split(nodeTags, " ")
	for _, taskTag := range taskSplit {
		found := false
		for _, nodeTag := range nodeSplit {
			if taskTag == nodeTag {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// getFirstMatch accepts a slice of Tasks and finds the first matching one
// based on the tags. It returns the found Task and a bool indicating the success.
func getFirstMatch(tasks []common.Task) (common.Task, bool) {
	for _, task := range tasks {
		if compareTags(task.Tags, viper.GetString("nodeTags")) {
			return task, true
		}
	}
	return common.Task{}, false
}

// getUnassigned accepts a slice of Tasks and returns a subslice
// containing the Tasks with no Node assegined to, and a bool indicating the success.
func getUnassigned(tasks []common.Task) ([]common.Task, bool) {
	found := false
	unassigned := []common.Task{}
	for _, task := range tasks {
		if task.Node == "None" {
			found = true
			unassigned = append(unassigned, task)
		}
	}
	return unassigned, found
}

// getTaskToExec polls the Infra Server for all the Tasks
// and using getUnassigned and getFirstMatch returns the
// first matching Task that can be executed, a bool indicating the succes
// and on optional error raised during the fetch.
func getTaskToExec() (common.Task, bool, error) {
	task := common.Task{}
	tasks, err := common.GetTasks()
	if err != nil {
		return task, false, errors.Wrap(err, "Cannot get tasks")
	}
	unassigned, _ := getUnassigned(tasks)
	taskToExec, ok := getFirstMatch(unassigned)
	if !ok {
		return task, false, nil
	}
	return taskToExec, true, nil
}
