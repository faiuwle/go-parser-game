package main

import (
	"github.com/faiuwle/go-parser-game/rage"
)

var entities []*rage.Entity

func main() {
	outOfWorld := &rage.Entity{
		Id:   0,
		Kind: "Room",
	}

	livingRoom := &rage.Entity{
		Id:          1,
		Name:        "Living Room",
		Description: "The living room",
		Kind:        "Room",
		Contents:    []int{3, 4},
	}

	bedroom := &rage.Entity{
		Id:          2,
		Name:        "Bedroom",
		Description: "The bedroom",
		Exits:       map[string]*rage.Entity{"south": livingRoom},
		Kind:        "Room",
	}

	player := &rage.Entity{
		Id:       3,
		Location: 1,
		Kind:     "Character",
	}

	key := &rage.Entity{
		Id:       4,
		Name:     "key",
		Location: 1,
		Kind:     "Thing",
	}

	livingRoom.Exits = map[string]*rage.Entity{"north": bedroom}

	entities = []*rage.Entity{outOfWorld, livingRoom, bedroom, player, key}

	rage.ParseLoop(entities, *player)
}
