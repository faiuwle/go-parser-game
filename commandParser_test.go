package rage_test

import (
	"github.com/faiuwle/go-parser-game/rage"
	"testing"
)

/*
func TestUnknownCommandFails(t *testing.T) {
	command := "notfound"

	_, err := rage.Parse(command)

	if err != rage.ErrorInvalidCommand {
		t.Error(err)
	}
}
*/

func TestInvalidCommandFails(t *testing.T) {
	command := "take key and run"

	_, err := rage.Parse(command)

	if err != rage.ErrorInvalidCommand {
		t.Error(err)
	}
}

func TestTakeKeySucceeds(t *testing.T) {
	command := "take key"

	cmd, err := rage.Parse(command)

	if cmd.Action != "take" || cmd.Noun != "key" {
		t.Errorf("Command did not parse successfully")
	}

	if err != nil {
		t.Error(err)
	}
}
