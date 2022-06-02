package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/faiuwle/go-parser-game/rage"
)

//go:embed "game.json"
var standardData string

func main() {
	data, err := rage.ReadConfig(strings.NewReader(standardData))
	if err != nil {
		fmt.Printf("Error reading game data: %s", err)
		os.Exit(1)
	}

	game, err := rage.NewGame(*data, os.Stdout)
	if err != nil {
		fmt.Printf("%v", err)
	}

	rage.Start(game)
}
