package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/mikloslorinczi/infra-exec/executor"
	"github.com/pkg/errors"
)

func executeTask(task common.Task) (string, error) {
	outputFilePath := "logs/" + task.ID + ".log"
	outputFile, err := executor.NewWriteFile(outputFilePath)
	if err != nil {
		return outputFilePath, errors.Wrapf(err, "Cannot open outputFile %v\n", outputFilePath)
	}

	defer func() {
		if err := outputFile.Close(); err != nil {
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

	req, err := http.NewRequest("POST", viper.GetString("apiURL")+"/log/"+ID, file)
	if err != nil {
		return errors.Wrap(err, "Cannot upload logfile to server")
	}
	req.Header.Set("Content-Type", "binary/octet-stream")
	req.Header.Set("APITOKEN", viper.GetString("apiToken"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "Cannot send request\n")
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
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
