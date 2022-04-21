package rage

import (
	"errors"
	"strings"
)

var (
	ErrorInvalidCommand = errors.New("invalid command")
)

type Command struct {
	Action string
	Noun   string
}

func Parse(command string) (Command, error) {
	commandParts := strings.Split(command, " ")

	if len(commandParts) == 1 {
		return Command{
			Action: commandParts[0],
		}, nil
	} else if len(commandParts) == 2 {
		return Command{
			Action: commandParts[0],
			Noun:   commandParts[1],
		}, nil
	} else {
		return Command{}, ErrorInvalidCommand
	}
}
