package main

import (
	"os"
	"testing"
)

func TestExecuteCommand_echo(t *testing.T) {
	outFile, err := os.Create(resFile)
	if err != nil {
		quit(1, "Cannot open Output File\n")
	}
	defer func() {
		err = outFile.Close()
		if err != nil {
			quit(1, "Error closeing Output File")
		}
	}()
	status, msg := execCommand("echo sajt", outFile)
	expStatus := 0
	expMsg := "Execution completed"
	if status != expStatus {
		t.Errorf("executeCommand failed for testing echo, got status: %v, expected status: %v.", status, expStatus)
	}
	if msg != expMsg {
		t.Errorf("executeCommand failed for testing echo, got msg: %v, expected msg: %v.", msg, expMsg)
	}
}
