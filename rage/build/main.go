package main

import (
	"github.com/faiuwle/go-parser-game/rage"
	"os"
)

func main() {
	// TODO this should be building a binary, but we don't know where it's getting put

	dataPath := os.Args[1]
	err := rage.Compile(dataPath, "./adventure")
	if err != nil {
		panic(err)
	}
}
