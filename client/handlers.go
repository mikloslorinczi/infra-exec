package main

import (
	"log"
	"strings"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/mikloslorinczi/infra-exec/executor"
	"github.com/pkg/errors"
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

func getFirstMatch(tasks []common.Task) (common.Task, bool) {
	for _, task := range tasks {
		if compareTags(task.Tags, nodeTags) {
			return task, true
		}
	}
	return common.Task{}, false
}

func claimTask(task common.Task) (common.ResponseMsg, error) {
	responseMsg := common.ResponseMsg{}
	taskJSON, err := common.TaskToJSON(task)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	responseJSON, err := common.SendRequest("POST", common.APIURL+"/task/claim", taskJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	responseMsg, err = common.JSONToMsg(responseJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	return responseMsg, nil
}

func executeTask(task common.Task) (string, error) {
	outputFilePath := "logs/" + task.ID + ".log"
	outputFile, err := executor.NewWriteFile(outputFilePath)
	if err != nil {
		return outputFilePath, errors.Wrapf(err, "Cannot open outputFile %v\n", outputFilePath)
	}

	defer func() {
		err = outputFile.Close()
		if err != nil {
			log.Fatalf("Cannot close outputFile %v\nError :\n%v\n", outputFilePath, err)
		}
	}()

	command, commandArgs := executor.ParseCommand(task.Command)
	if err := executor.ExecCommand(command, commandArgs, outputFile); err != nil {
		return outputFilePath, errors.Wrapf(err, "Cannot execute command %v\n", task.Command)
	}

	return outputFilePath, nil
}
