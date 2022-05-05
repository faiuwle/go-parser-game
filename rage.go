package rage

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Entity struct {
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
	Destination string
	Requires    string
}

func Start(game *Game) {
	fmt.Println("Welcome to the text adventure, type commands to play.")
	fmt.Println(game.Entities[game.Player.Location].Description)
	fmt.Println(ListExits(*game.Entities[game.Player.Location]))
	fmt.Println(game.Entities[game.Player.Location].ListContents())
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

		currentRoom := game.Entities[game.Player.Location]

		switch cmd.Action {
		case "look":
			fmt.Println(currentRoom.Description)
		case "quit":
			fmt.Println("Thanks for playing!")
			os.Exit(0)
		case "inventory":
			fmt.Println(game.ListInventory())
		case "take":
			fmt.Println(game.TakeItem(cmd.Noun))
		default:
			exit, ok := currentRoom.Exits[input]
			if !ok {
				fmt.Println("Sorry I didn't understand.")
				break
			}

			err = game.SetPlayerLocation(exit.Destination)
			if err != nil {
				fmt.Println(err)
				fmt.Print("> ")
				continue
			}
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

func ListExits(room Entity) string {
	exits := maps.Keys(room.Exits)
	slices.Sort(exits)
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
	g.MoveEntity(g.Player.Name, destination.Name)

	fmt.Println(destination.Description)
	fmt.Println(ListExits(*destination))
	fmt.Println(destination.ListContents())

	return nil
}

func (g *Game) PlayerLocation() string {
	return g.Entities[g.Player.Location].Name
}

func (g *Game) ListInventory() string {
	return "You are carrying " + FormatItems(g.Player.Contents) + "."
}

func (g *Game) TakeItem(itemName string) string {
	currentRoom := g.Entities[g.PlayerLocation()]
	if !currentRoom.Contains(itemName) {
		return "I can't see that here."
	}
	g.MoveEntity(itemName, g.Player.Name)
	return fmt.Sprintf("You take the %s.", itemName)
}

func (g *Game) MoveEntity(entityToMove string, destination string) {
	e := g.Entities[entityToMove]
	location := g.Entities[e.Location]
	idx := slices.Index(location.Contents, entityToMove)
	location.Contents = slices.Delete(location.Contents, idx, idx+1)

	d := g.Entities[destination]
	d.Contents = append(d.Contents, entityToMove)
	e.Location = d.Name
}

type GameData map[string]*Entity

func NewGame(data GameData, startPlayer string) (*Game, error) {
	return &Game{
		Entities: data,
		Player:   data[startPlayer],
	}, nil
}
