package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

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

func addTask() (common.ResponseMsg, error) {
	var (
		responseMsg  common.ResponseMsg
		commandJSON  []byte
		responseJSON []byte
	)
	commandObj, err := readCommand()
	if err != nil {
		return responseMsg, errors.Wrap(err, "Error reading commoand")
	}
	err = common.ToJSON(&commandObj, &commandJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	responseJSON, err = common.SendRequest("POST", common.APIURL+"/tasks", commandJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	err = common.FromJSON(&responseMsg, responseJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot add new task")
	}
	return responseMsg, nil
}

func listTasks() ([]common.Task, error) {
	tasks, err := common.GetTasks()
	if err != nil {
		return tasks, errors.Wrap(err, "Error getting tasks")
	}
	return tasks, nil
}

func queryTask(query string) (common.Task, error) {
	var task common.Task
	taskJSON, err := common.SendRequest("GET", common.APIURL+"/task/"+query, nil)
	if err != nil {
		return task, errors.Wrap(err, "Cannot query task")
	}
	err = common.FromJSON(&task, taskJSON)
	if err != nil {
		return task, errors.Wrap(err, "Cannot query task")
	}
	return task, nil
}

func getLog(ID string) error {
	task, err := queryTask(logs)
	if err != nil {
		return errors.Wrapf(err, "Error geting Task %v\n", ID)
	}
	if !strings.Contains(task.Status, "log available") {
		return fmt.Errorf("No log available for Task: %v Status: %v", task.ID, task.Status)
	}
	err = downloadLog(task.ID)
	if err != nil {
		return errors.Wrapf(err, "Error downloading logfile %v.log\n", task.ID)
	}
	return nil
}

func downloadLog(ID string) error {
	filepath := path.Join("logs/", ID+".log")
	outFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "Error opening logfile %v\n", filepath)
	}
	defer func() {
		closeErr := outFile.Close()
		if closeErr != nil {
			log.Fatalf("Error closing logfile %v\n%v\n", filepath, err)
		}
	}()

	req, err := http.NewRequest("GET", common.APIURL+"/log/"+ID, nil)
	if err != nil {
		return errors.Wrap(err, "Cannot upload logfile to server")
	}
	req.Header.Set("adminpassword", common.AdminPass)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Error requesting logfile")
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Fatalf("Error closing response body %v\n", closeErr)
		}
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Error writinng logfile %v\n", filepath)
	}

	return nil
}
