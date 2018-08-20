package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func getFirstMatch(tasks []common.Task) (common.Task, bool) {
	for _, task := range tasks {
		if compareTags(task.Tags, nodeTags) {
			return task, true
		}
	}
	return common.Task{}, false
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

func claimTask(task common.Task) (common.ResponseMsg, error) {
	var (
		responseMsg  common.ResponseMsg
		responseJSON []byte
		taskJSON     []byte
	)
	err := common.ToJSON(&task, &taskJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	responseJSON, err = common.SendRequest("POST", common.APIURL+"/task/claim", taskJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	err = common.FromJSON(&responseMsg, responseJSON)
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

func uploadLog(path, ID string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return errors.Wrap(err, "Cannot open logfile")
	}
	// File already closed. Why?
	// defer func() {
	// 	err = file.Close()
	// 	if err != nil {
	// 		log.Fatalf("Error closing file %v\n%v", path, err)
	// 	}
	// }()

	req, err := http.NewRequest("POST", common.APIURL+"/log/upload/"+ID, file)
	if err != nil {
		return errors.Wrap(err, "Cannot upload logfile to server")
	}
	req.Header.Set("Content-Type", "binary/octet-stream")
	req.Header.Set("adminpassword", common.AdminPass)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "Cannot send request\n")
		return err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			closeErr = errors.Wrap(closeErr, "Error closing response body\n")
			log.Fatal(closeErr)
		}
	}()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading response body\n")
		return err
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("Server answered with a non-200 status: %v\n", resp.StatusCode)
	}

	return nil
}
