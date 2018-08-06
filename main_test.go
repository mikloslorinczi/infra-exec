package main

import "testing"

func TestCheckCommand_with_ls(t *testing.T) {
	res := checkCommand("ls")
	exp := true
	if res != exp {
		t.Errorf("checkCommand failed for testing ls, got: %t, want: %t.", res, exp)
	}
}

func TestCheckCommand_with_asd(t *testing.T) {
	res := checkCommand("asd")
	exp := false
	if res != exp {
		t.Errorf("checkCommand failed for testing asd, got: %t, want: %t.", res, exp)
	}
}

func TestCheckCommand_with_git(t *testing.T) {
	res := checkCommand("git")
	exp := true
	if res != exp {
		t.Errorf("checkCommand failed for testing git, got: %t, want: %t.", res, exp)
	}
}
