package main

import "testing"

func TestUnknownCommandFails(t *testing.T) {
	command := "notfound"

	_, err := Parse(command)

	if err != ErrorInvalidCommand {
		t.Error(err)
	}
}

func TestInvalidCommandFails(t *testing.T) {
	command := "take key and run"

	_, err := Parse(command)

	if err != ErrorInvalidCommand {
		t.Error(err)
	}
}

func TestTakeKeySucceeds(t *testing.T) {
	command := "take key"

	cmd, err := Parse(command)

	if cmd.Action != "take" || cmd.Noun != "key" {
		t.Errorf("Command did not parse successfully")
	}

	if err != nil {
		t.Error(err)
	}
}
