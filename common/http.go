package common

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// SendRequest creates a HTTP request with the given Method to the given URL with optional body.
// AdminPass will be set as custom HTTP header "adminpassword".
// The function returns the response body's byte-representation (JSON), and on optional error.
func SendRequest(method string, url string, body []byte) ([]byte, error) {
	response := []byte{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		err = errors.Wrapf(err, "Cannot create %v request to %v\n", method, url)
		return response, err
	}
	req.Header.Set("adminpassword", AdminPass)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "Cannot send request\n")
		return response, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			closeErr = errors.Wrap(closeErr, "Error closing response body\n")
			log.Fatal(closeErr)
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

// GetTasks sends a request to the APIURL/task/list and returns
// the fetched tasks as a slice, and on optional http error.
func GetTasks() ([]Task, error) {
	tasksJSON, err := SendRequest("GET", APIURL+"/task/list", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get task list")
	}
	tasks, err := JSONToTasks(tasksJSON)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get task list")
	}
	return tasks, nil
}
