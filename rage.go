package rage

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

type Entity struct {
	Id          int
	Name        string
	Description string
	Exits       map[string]*Entity
	Location    int
	Contents    []int
	Kind        string
}

func ParseLoop(entities []*Entity, player Entity) {
	fmt.Println("Welcome to the text adventure, type commands to play.")
	fmt.Println(entities[player.Location].Description)
	fmt.Println(ListExits(entities[player.Location]))
	fmt.Println(ListItems(entities[player.Location], entities))
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		input := scanner.Text()
		cmd, err := Parse(input)

		if err != nil {
			fmt.Println(err)
			fmt.Print("> ")
			continue
		}

		currentRoom := entities[player.Location]

		switch cmd.Action {
		case "look":
			fmt.Println(currentRoom.Description)
		case "quit":
			fmt.Println("Thanks for playing!")
			os.Exit(0)
		case "inventory":
			var things []string

			for _, thingId := range player.Contents {
				things = append(things, entities[thingId].Name)
			}

			fmt.Printf("You are carrying %s\n", FormatItems(things))
		case "take":
			for index, thingId := range currentRoom.Contents {
				if entities[thingId].Name == cmd.Noun {
					thing := entities[thingId]

					player.Contents = append(player.Contents, thingId)
					currentRoom.Contents = slices.Delete(currentRoom.Contents, index, index+1)

					thing.Location = player.Id

					fmt.Printf("You take the %s.\n", cmd.Noun)
					break
				}
			}
		default:
			newRoom, ok := currentRoom.Exits[input]
			if !ok {
				fmt.Println("Sorry I didn't understand.")
				break
			}
			player.Location = slices.Index(entities, newRoom)
			currentRoom = entities[player.Location]
			fmt.Println(currentRoom.Description)
			fmt.Println(ListExits(currentRoom))
			fmt.Println(ListItems(currentRoom, entities))
		}

		fmt.Print("> ")
	}
}

func FormatItems(input []string) string {
	switch len(input) {
	case 0:
		return ""

	case 1:
		return input[0]

	case 2:
		return input[0] + " and " + input[1]

	default:
		listString := strings.Join(input[:len(input)-1], ", ")
		listString += ", and " + input[len(input)-1]
		return listString
	}
}

func ListExits(room *Entity) string {
	exitString := "You can go"

	exits := maps.Keys(room.Exits)
	if len(exits) == 0 {
		return "There are no visible Exits."
	}
	if len(exits) == 1 {
		return exitString + " " + exits[0]
	}

	exitString += strings.Join(exits[:len(exits)-1], ",")
	exitString += " and " + exits[len(exits)-1]

	return exitString
}

func ListItems(room *Entity, entities []*Entity) string {
	var itemNames []string

	for _, itemID := range room.Contents {
		name := entities[itemID].Name
		if name != "" {
			itemNames = append(itemNames, entities[itemID].Name)
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
