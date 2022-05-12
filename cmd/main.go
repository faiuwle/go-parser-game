package main

import (
	"fmt"
	"os"

	"github.com/faiuwle/go-parser-game/rage"
)

func main() {
	data := map[string]*rage.Entity{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination:    "Bedroom",
					Requires:       "key",
					FailureMessage: "You rattle the handle, but the door seems locked.",
				},
			},
			Contents: []string{"key", "Shera"},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Exits: map[string]rage.Exit{
				"south": {
					Destination: "Living Room",
				},
			},
			Kind: "Room",
		},
		"key": {
			Name:     "key",
			Location: "Living Room",
			Kind:     "Thing",
		},
		"Shera": {
			Name:     "Shera",
			Location: "Living Room",
		},
	}

	game, err := rage.NewGame(data, "Shera", os.Stdout)
	if err != nil {
		fmt.Printf("%v", err)
	}

	rage.Start(game)
}
