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
	fmt.Println(ListExits(*entities[player.Location]))
	fmt.Println(ListItems(*entities[player.Location], entities))
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
			fmt.Println(TakeItem(entities, &player, currentRoom, cmd.Noun))
		default:
			newRoom, ok := currentRoom.Exits[input]
			if !ok {
				fmt.Println("Sorry I didn't understand.")
				break
			}
			player.Location = slices.Index(entities, newRoom)
			currentRoom = entities[player.Location]
			fmt.Println(currentRoom.Description)
			fmt.Println(ListExits(*currentRoom))
			fmt.Println(ListItems(*currentRoom, entities))
		}

		fmt.Print("> ")
	}
}

func TakeItem(entities []*Entity, player *Entity, currentRoom *Entity, noun string) string {
	for index, thingId := range currentRoom.Contents {
		if entities[thingId].Name == noun {
			thing := entities[thingId]

			player.Contents = append(player.Contents, thingId)
			currentRoom.Contents = slices.Delete(currentRoom.Contents, index, index+1)

			thing.Location = player.Id

			return fmt.Sprintf("You take the %s.", noun)
		}
	}

	return "I can't see that here."
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

func ListExits(room Entity) string {
	exits := maps.Keys(room.Exits)
	exitList := FormatItems(exits)

	if exitList == "" {
		return "There are no visible Exits."
	}

	return "You can go " + exitList + " from here."
}

func ListItems(room Entity, entities []*Entity) string {
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

	return "You see here " + FormatItems(itemNames) + "."
}
