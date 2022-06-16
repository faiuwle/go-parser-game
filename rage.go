package rage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

const (
	DefaultExitFailureMessage = "You cannot do that."
)

var ErrorExitRequirementNotMet = errors.New("exit requirement not met")

func (e *Entity) Contains(name string) bool {
	return slices.Contains(e.Contents, name)
}

func (e *Entity) ListContents() string {
	items := make([]string, 0, len(e.Contents))
	for _, item := range e.Contents {
		if item == "player" {
			continue
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return ""
	}
	return "You see here " + FormatItems(items) + "."
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

type Exit struct {
	Destination    string
	Requires       string
	FailureMessage string
}

func (e *Exit) GetFailureMessage() string {
	if e.FailureMessage == "" {
		return DefaultExitFailureMessage
	}

	return e.FailureMessage
}

func Start(game *Game) {
	fmt.Println("Welcome to the text adventure, type commands to play.")
	startRoom := game.PlayerLocation()
	fmt.Println(startRoom.Description)
	fmt.Println(ListExits(*startRoom))
	contents := startRoom.ListContents()
	if contents != "" {
		fmt.Println(contents)
	}
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

		err = game.Do(cmd)

		if errors.Is(err, ErrorInvalidCommand) {
			fmt.Println("Sorry I didn't understand.")
		} else if err != nil {
			fmt.Println(err)
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

type Game struct {
	Entities map[string]*Entity
	Output   io.Writer
	Player   *Entity
}

func (g *Game) Do(cmd Command) error {
	currentRoom := g.PlayerLocation()

	switch cmd.Action {
	case "look":
		g.Say(currentRoom.Description)
	case "quit":
		g.Say("Thanks for playing!")
		os.Exit(0)
	case "inventory":
		g.Say(g.ListInventory())
	case "take":
		g.Say(g.TakeItem(cmd.Noun))
	// case "north":
	// case "south":
	// case "east":
	// case "west":
	default:
		exit, ok := currentRoom.Exits[cmd.Action]
		if !ok {
			return ErrorInvalidCommand
		}

		err := g.SetPlayerLocation(exit)

		if errors.Is(err, ErrorExitRequirementNotMet) {
			g.Say(exit.GetFailureMessage())
			return nil
		}

		currentRoom = g.PlayerLocation()
		g.Say(currentRoom.Description)
		g.Say(ListExits(*currentRoom))
		contents := currentRoom.ListContents()
		if contents != "" {
			g.Say(contents)
		}
	}

	return nil
}

func (g *Game) GetEntity(entityName string) *Entity {
	entity, ok := g.Entities[entityName]
	if !ok {
		message := fmt.Sprintf("Inconsistent game data: %s not in entity list", entityName)
		panic(message)
	}
	return entity
}

func (g *Game) Say(message string) {
	fmt.Fprintln(g.Output, message)
}

func (g *Game) SetPlayerLocation(exit Exit) error {
	if exit.Requires != "" && !slices.Contains(g.Player.Contents, exit.Requires) {
		return ErrorExitRequirementNotMet
	}

	g.MoveEntity(g.Player.Name, exit.Destination)

	return nil
}

func (g *Game) PlayerLocation() *Entity {
	return g.GetEntity(g.Player.Location)
}

func (g *Game) ListInventory() string {
	return "You are carrying " + FormatItems(g.Player.Contents) + "."
}

func (g *Game) TakeItem(itemName string) string {
	currentRoom := g.PlayerLocation()
	if !currentRoom.Contains(itemName) {
		return "I can't see that here."
	}
	g.MoveEntity(itemName, g.Player.Name)
	return fmt.Sprintf("You take the %s.", itemName)
}

func (g *Game) MoveEntity(entityToMove string, destination string) {
	entity := g.GetEntity(entityToMove)
	location := g.GetEntity(entity.Location)
	idx := slices.Index(location.Contents, entityToMove)
	location.Contents = slices.Delete(location.Contents, idx, idx+1)

	d := g.GetEntity(destination)
	d.Contents = append(d.Contents, entityToMove)
	entity.Location = d.Name
}

type GameData map[string]Entity

func (gd GameData) Missing(entityName string) bool {
	_, ok := gd[entityName]
	return !ok
}

func NewGame(data GameData, output io.Writer) (*Game, error) {
	g := Game{
		Entities: map[string]*Entity{},
		Output:   output,
	}
	entities := g.Entities
	for key, val := range data {
		entity := val

		if entity.Kind == "Room" && entity.Location != "" {
			return nil, fmt.Errorf("rooms cannot have locations: %#v has a location", val)
		}
		if entity.Kind != "Room" && data.Missing(entity.Location) {
			return nil, fmt.Errorf("%#v has invalid location", val)
		}
		entities[key] = &entity
	}
	g.Player = g.Entities["player"]
	if g.Player == nil {
		return nil, errors.New("player missing from game data")
	}

	return &g, nil
}
