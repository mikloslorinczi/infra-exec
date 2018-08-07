package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseCommand_empty(t *testing.T) {
	name, args := parseCommand("")
	expName := ""
	expArgs := []string{}
	if name != expName {
		t.Errorf("parseCommand failed for testing echo, expected name: %v, got name: %v.", expName, name)
	}
	if !reflect.DeepEqual(args, expArgs) {
		t.Errorf("ParseCommand failed for testing echo, expected args: %v, got args: %v.", expArgs, args)
	}
}

func TestParseCommand_echo(t *testing.T) {
	name, args := parseCommand("echo")
	expName := "echo"
	expArgs := []string{}
	if name != expName {
		t.Errorf("parseCommand failed for testing echo, expected name: %v, got name: %v.", expName, name)
	}
	if !reflect.DeepEqual(args, expArgs) {
		t.Errorf("ParseCommand failed for testing echo, expected args: %v, got args: %v.", expArgs, args)
	}
}

func TestParseCommand_echo_sajt(t *testing.T) {
	name, args := parseCommand("echo sajt")
	expName := "echo"
	expArgs := []string{"sajt"}
	if name != expName {
		t.Errorf("parseCommand failed for testing echo, expected name: %v, got name: %v.", expName, name)
	}
	if !reflect.DeepEqual(args, expArgs) {
		t.Errorf("ParseCommand failed for testing echo, expected args: %v, got args: %v.", expArgs, args)
	}
}

func TestParseCommand_docker(t *testing.T) {
	name, args := parseCommand("docker run -i -t -d -p 80:8080 tutum/hello-world")
	expName := "docker"
	expArgs := []string{"run", "-i", "-t", "-d", "-p", "80:8080", "tutum/hello-world"}
	if name != expName {
		t.Errorf("parseCommand failed for testing echo, expected name: %v, got name: %v.", expName, name)
	}
	if !reflect.DeepEqual(args, expArgs) {
		t.Errorf("ParseCommand failed for testing echo, expected args: %v, got args: %v.", expArgs, args)
	}
}

func TestGetWriteFile(t *testing.T) {
	_, err := getWriteFile("testFile")
	if err != nil {
		t.Errorf("getWriteFile failed geting WriteFile %v", err)
	}
	_, err = os.Stat("testFile")
	if os.IsNotExist(err) {
		t.Errorf("getWriteFile failed to create testFile %v", err)
	}
	err = os.Remove("testFile")
	if err != nil {
		t.Errorf("Error removeing testFile")
	}
}

// if the previous test runs without error, we can assume that getWriteFile works properly, and wont throw an error
func TestExecCommand_space(t *testing.T) {
	testFile, _ := getWriteFile("testFile")
	defer testFile.Close()
	defer os.Remove("testFile")
	err := execCommand(" ", testFile)
	if err != nil {
		t.Errorf(`-c " " should produce an empty outputFile and no error, got: %v`, err)
	}
}
