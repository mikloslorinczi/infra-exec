package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

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

	req, err := http.NewRequest("GET", common.APIURL+"/log/download/"+ID, nil)
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
