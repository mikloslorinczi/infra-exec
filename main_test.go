package main

import (
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

func TestParseCommand_echo_Lali(t *testing.T) {
	name, args := parseCommand("echo a")
	expName := "echo"
	expArgs := []string{"a"}
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
