package main

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

type Entity struct {
	id          int
	name        string
	description string
	exits       map[string]*Entity
	location    int
	contents    []int
	kind        string
}

var entities []*Entity

func main() {
	outOfWorld := &Entity{
		id:   0,
		kind: "Room",
	}

	livingRoom := &Entity{
		id:          1,
		name:        "Living Room",
		description: "The living room",
		kind:        "Room",
		contents:    []int{3, 4},
	}

	bedroom := &Entity{
		id:          2,
		name:        "Bedroom",
		description: "The bedroom",
		exits:       map[string]*Entity{"south": livingRoom},
		kind:        "Room",
	}

	player := &Entity{
		id:       3,
		location: 1,
		kind:     "Character",
	}

	key := &Entity{
		id:       4,
		name:     "key",
		location: 1,
		kind:     "Thing",
	}

	livingRoom.exits = map[string]*Entity{"north": bedroom}

	entities = []*Entity{outOfWorld, livingRoom, bedroom, player, key}

	fmt.Println("Welcome to the text adventure, type commands to play.")
	fmt.Println(entities[player.location].description)
	fmt.Println(ListExits(entities[player.location]))
	fmt.Println(ListItems(entities[player.location]))
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		input := scanner.Text()

		currentRoom := entities[player.location]

		switch input {
		case "look":
			fmt.Println(currentRoom.description)
		case "quit":
			fmt.Println("Thanks for playing!")
			os.Exit(0)
		default:
			newRoom, ok := currentRoom.exits[input]
			if !ok {
				fmt.Println("Sorry I didn't understand.")
				break
			}
			player.location = slices.Index(entities, newRoom)
			currentRoom = entities[player.location]
			fmt.Println(currentRoom.description)
			fmt.Println(ListExits(currentRoom))
			fmt.Println(ListItems(currentRoom))
		}

		fmt.Print("> ")
	}
}

func ListExits(room *Entity) string {
	exitString := "You can go"

	exits := maps.Keys(room.exits)
	if len(exits) == 0 {
		return "There are no visible exits."
	}
	if len(exits) == 1 {
		return exitString + " " + exits[0]
	}

	exitString += strings.Join(exits[:len(exits)-1], ",")
	exitString += " and " + exits[len(exits)-1]

	return exitString
}

func ListItems(room *Entity) string {
	var itemNames []string

	for _, itemID := range room.contents {
		name := entities[itemID].name
		if name != "" {
			itemNames = append(itemNames, entities[itemID].name)
		}
	}

	if len(itemNames) == 0 {
		return ""
	}

	itemString := "You see here "

	if len(itemNames) == 1 {
		return itemString + itemNames[0] + "."
	}

	itemString += strings.Join(itemNames, ",")
	itemString += " and " + itemNames[len(itemNames)-1] + "."

	return itemString
}
