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

var (
	ErrorExitRequirementNotMet = errors.New("exit requirement not met")
)

func (e *Entity) Contains(name string) bool {
	return slices.Contains(e.Contents, name)
}

func (e *Entity) ListContents() string {
	return "You see here " + FormatItems(e.Contents) + "."
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
	Entities GameData
	Output   io.Writer
	Player   *Entity
}

func (g *Game) Do(cmd Command) error {
	currentRoom := g.Entities[g.Player.Location]

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
	//case "north":
	//case "south":
	//case "east":
	//case "west":
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

		if err != nil {
			return err
		}

		currentRoom = g.Entities[g.Player.Location]
		g.Say(currentRoom.Description)
		g.Say(ListExits(*currentRoom))
		g.Say(currentRoom.ListContents())
	}

	return nil
}

func (g *Game) Say(message string) {
	fmt.Fprintln(g.Output, message)
}

func (g *Game) SetPlayerLocation(exit Exit) error {
	if exit.Requires != "" && !slices.Contains(g.Player.Contents, exit.Requires) {
		return ErrorExitRequirementNotMet
	}

	destination, ok := g.Entities[exit.Destination]
	if !ok {
		return fmt.Errorf("unknown location %q", exit.Destination)
	}

	g.MoveEntity(g.Player.Name, destination.Name)

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
	entity := g.Entities[entityToMove]
	location := g.Entities[entity.Location]
	idx := slices.Index(location.Contents, entityToMove)
	location.Contents = slices.Delete(location.Contents, idx, idx+1)

	d := g.Entities[destination]
	d.Contents = append(d.Contents, entityToMove)
	entity.Location = d.Name
}

type GameData map[string]*Entity

func NewGame(data GameData, startPlayer string, output io.Writer) (*Game, error) {
	return &Game{
		Entities: data,
		Output:   output,
		Player:   data[startPlayer],
	}, nil
}
