package rage

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Entity struct {
	Id          int
	Name        string
	Description string
	Exits       map[string]Exit
	Location    string
	Contents    []string
	Kind        string
}

func (e *Entity) Contains(name string) bool {
	return slices.Contains(e.Contents, name)
}

func (e *Entity) ListContents() string {
	return "You see here " + FormatItems(e.Contents) + "."
}

type Exit struct {
	Destination int
	Requires    int
}

// func ParseLoop(entities []*Entity, player Entity) {
// 	fmt.Println("Welcome to the text adventure, type commands to play.")
// 	fmt.Println(entities[player.Location].Description)
// 	fmt.Println(ListExits(*entities[player.Location]))
// 	fmt.Println(ListItems(*entities[player.Location], entities))
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Print("> ")

// 	for scanner.Scan() {
// 		input := scanner.Text()
// 		cmd, err := Parse(input)
// 		if err != nil {
// 			fmt.Println(err)
// 			fmt.Print("> ")
// 			continue
// 		}

// 		currentRoom := entities[player.Location]

// 		switch cmd.Action {
// 		case "look":
// 			fmt.Println(currentRoom.Description)
// 		case "quit":
// 			fmt.Println("Thanks for playing!")
// 			os.Exit(0)
// 		case "inventory":
// 			var things []string

// 			for _, thingId := range player.Contents {
// 				things = append(things, entities[thingId].Name)
// 			}

// 			fmt.Printf("You are carrying %s\n", FormatItems(things))
// 		case "take":
// 			fmt.Println(TakeItem(entities, &player, currentRoom, cmd.Noun))
// 		default:
// 			exit, ok := currentRoom.Exits[input]
// 			if !ok {
// 				fmt.Println("Sorry I didn't understand.")
// 				break
// 			}
// 			player.Location = exit.Destination
// 			currentRoom = entities[player.Location]
// 			fmt.Println(currentRoom.Description)
// 			fmt.Println(ListExits(*currentRoom))
// 			fmt.Println(ListItems(*currentRoom, entities))
// 		}

// 		fmt.Print("> ")
// 	}
// }

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

type Game struct {
	Entities GameData
	Player   *Entity
}

func (g *Game) Do(string) error {
	return nil
}

func (g *Game) SetPlayerLocation(location string) error {
	destination, ok := g.Entities[location]
	if !ok {
		return fmt.Errorf("unknown location %q", location)
	}
	g.MoveEntity("player", destination)
	return nil
}

func (g *Game) PlayerLocation() string {
	return g.Entities[g.Player.Location].Name
}

func (g *Game) TakeItem(itemName string) string {
	currentRoom := g.Entities[g.PlayerLocation()]
	if !currentRoom.Contains(itemName) {
		return "I can't see that here."
	}
	g.MoveEntity(itemName, g.Player)
	return fmt.Sprintf("You take the %s.", itemName)
}

func (g *Game) MoveEntity(name string, destination *Entity) {
	e := g.Entities[name]
	location := g.Entities[e.Location]
	idx := slices.Index(location.Contents, name)
	location.Contents = slices.Delete(location.Contents, idx, idx+1)
	destination.Contents = append(destination.Contents, name)
	e.Location = destination.Name
}

type GameData map[string]*Entity

func NewGame(data GameData, startLocation string) (*Game, error) {
	return &Game{
		Entities: data,
		Player: &Entity{
			Name:     "Default Player",
			Location: startLocation,
		},
	}, nil
}
