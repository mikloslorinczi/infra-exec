package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

//ResponseMsg is a general response message from infra-server.
type ResponseMsg struct {
	Msg string `json:"msg"`
}

//CommandJSON is the JSON fromat to be sent to the infra-server.
type CommandJSON struct {
	Command string `json:"command"`
	Body    string `json:"body"`
}

//Task is the JSON structure of a single task.
type Task struct {
	ID      string `json:"id"`
	Node    string `json:"node"`
	Tags    string `json:"tags"`
	Status  string `json:"status"`
	Command string `json:"command"`
}

func commandToJSON(command string, body string) ([]byte, error) {
	commandJSON := CommandJSON{command, body}
	bytes, err := json.Marshal(commandJSON)
	if err != nil {
		err = errors.Wrapf(err, "Error encoding JSON command: %v body: %v\n", command, body)
		return bytes, err
	}
	return bytes, nil
}

func bodyToTasks(body []byte) ([]Task, error) {
	var tasks []Task
	err := json.Unmarshal(body, &tasks)
	if err != nil {
		err = errors.Wrapf(err, "Error decoding JSON body\n%v\n", body)
		return tasks, err
	}
	return tasks, nil
}

func bodyToMsg(body []byte) (ResponseMsg, error) {
	var response ResponseMsg
	err := json.Unmarshal(body, &response)
	if err != nil {
		err = errors.Wrapf(err, "Error decoding JSON body\n%v\n", body)
		return response, err
	}
	return response, nil
}

func postRequest(url string, body []byte) ([]byte, error) {
	response := []byte{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		err = errors.Wrapf(err, "Cannot create POST request to %v\n", url)
		return response, err
	}
	req.Header.Set("adminpassword", envPass)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "Cannot send request\n")
		return response, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			err = errors.Wrap(err, "Error closing response body\n")
			log.Fatal(err)
		}
	}()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading response body\n")
		return response, err
	}
	if resp.StatusCode != 200 {
		return response, errors.Errorf("Server answered with a non-200 status: %v\n", resp.StatusCode)
	}
	return response, nil
}
