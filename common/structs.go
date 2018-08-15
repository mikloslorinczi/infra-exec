package common

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	// AdminPass is the password used to authenticate the admin.
	AdminPass string
	// APIURL is the address of the infra-server's api. By default it will be "http://localhost:7474/api".
	APIURL string
)

// ResponseMsg is a general response message from infra-server, usually related to an error.
type ResponseMsg struct {
	Msg string `json:"msg"`
}

// CommandObj is the representation of a command ans its tags.
// It is ment to be sent to the infra-server so it can make a task out of it for the infra-client(s).
type CommandObj struct {
	Command string `json:"command"`
	Tags    string `json:"tags"`
}

// Task is the structure of a single task. ID is provided by github.com/rs/xid. It is unique and the
// time of creation can be parsed from it. Node is the name of the executing infra-client, initially it is "none".
// Tags hold the desired tags separated by space " ", only infra-clients with these tags can execute the task.
// Status is the current status of the task. Initially it is "Created".
// Command is the actual command and its argumentums separated by space " ".
type Task struct {
	ID      string `json:"id"`
	Node    string `json:"node"`
	Tags    string `json:"tags"`
	Status  string `json:"status"`
	Command string `json:"command"`
}

// CommandToJSON converts a CommandObj into JSON.
// It returns the byte representation of the JSON and an optional encoding error.
func CommandToJSON(command CommandObj) ([]byte, error) {
	bytes, err := json.Marshal(command)
	if err != nil {
		err = errors.Wrapf(err, "Error encoding JSON command: %v tags: %v\n", command.Command, command.Tags)
		return nil, err
	}
	return bytes, nil
}

// JSONToTasks converts a JSON into a slice of Tasks.
// It returns the slice and on optional decoding error.
func JSONToTasks(body []byte) ([]Task, error) {
	var tasks []Task
	err := json.Unmarshal(body, &tasks)
	if err != nil {
		err = errors.Wrapf(err, "Error decoding JSON body\n%s\n", body)
		return nil, err
	}
	return tasks, nil
}

// JSONToTask converts a JSON into a single Task.
// It returns the Task and an optional decoding error.
func JSONToTask(body []byte) (Task, error) {
	var task Task
	err := json.Unmarshal(body, &task)
	if err != nil {
		err = errors.Wrapf(err, "Error decoding JSON body\n%s\n", body)
		return task, err
	}
	return task, nil
}

// JSONToMsg converts a JSON into a ResponseMsg
// It returns the ResponseMsg and on optional decoding error.
func JSONToMsg(body []byte) (ResponseMsg, error) {
	var response ResponseMsg
	err := json.Unmarshal(body, &response)
	if err != nil {
		err = errors.Wrapf(err, "Error decoding JSON body\n%v\n", body)
		return response, err
	}
	return response, nil
}
