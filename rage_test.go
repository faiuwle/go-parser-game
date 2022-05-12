package rage_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/faiuwle/go-parser-game/rage"
)

func TestFormatItems(t *testing.T) {
	testcases := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "Three items",
			items: []string{"one", "two", "three"},
			want:  "one, two, and three",
		},
		{
			name:  "Two items",
			items: []string{"one", "two"},
			want:  "one and two",
		},
		{
			name:  "One item",
			items: []string{"one"},
			want:  "one",
		},
		{
			name:  "No items",
			items: []string{},
			want:  "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			list := rage.FormatItems(tc.items)

			if list != tc.want {
				t.Errorf("Wanted %q, got %q", tc.want, list)
			}
		})
	}
}

func TestListExits(t *testing.T) {
	room := rage.Entity{
		Exits: map[string]rage.Exit{
			"north": {},
			"east":  {},
		},
	}

	exitString := rage.ListExits(room)
	want := "You can go east and north from here."

	if exitString != want {
		t.Errorf("Wanted %q, got %q", want, exitString)
	}
}

func TestListContents(t *testing.T) {
	room := rage.Entity{
		Contents: []string{"key", "phone", "chocolate"},
	}

	got := room.ListContents()
	want := "You see here key, phone, and chocolate."

	if got != want {
		t.Errorf("Wanted %q, got %q", want, got)
	}
}

func TestTakeItemSucceedsIfItemIsPresent(t *testing.T) {
	entities := rage.GameData{
		"Room": {
			Name:     "Room",
			Contents: []string{"key", "phone", "player"},
		},
		"key": {
			Name:     "key",
			Location: "Room",
		},
		"phone": {
			Name:     "phone",
			Location: "Room",
		},
		"player": {
			Name:     "player",
			Location: "Room",
		},
	}
	g, err := rage.NewGame(entities, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	response := g.TakeItem("key")
	want := "You take the key."

	if response != want {
		t.Errorf("Wanted %q, got %q", want, response)
	}
	if !g.Player.Contains("key") {
		t.Error("Key was not transferred to inventory.")
	}
	if g.Entities["Room"].Contains("key") {
		t.Errorf("Key is still in the room.")
	}
}

func TestTakeItemFailsIfItemIsNotPresent(t *testing.T) {
	entities := rage.GameData{
		"Room": {
			Name:     "Room",
			Contents: []string{"key", "player"},
		},
		"key": {
			Name:     "key",
			Location: "Room",
		},
		"phone": {
			Name: "phone",
		},
		"player": {
			Name:     "player",
			Location: "Room",
		},
	}
	g, err := rage.NewGame(entities, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	response := g.TakeItem("phone")
	want := "I can't see that here."

	if response != want {
		t.Errorf("Wanted %q, got %q", want, response)
	}

	if g.Player.Contains("phone") {
		t.Errorf("Player took the phone.")
	}
}

func TestPlayerLocationReturnsNameOfRoomWherePlayerIs(t *testing.T) {
	t.Parallel()
	data := map[string]*rage.Entity{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination: "Bedroom",
				},
			},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Kind:        "Room",
			Contents:    []string{"player"},
		},
		"player": {
			Name:     "player",
			Location: "Bedroom",
		},
	}
	game, err := rage.NewGame(data, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	want := "Bedroom"
	got := game.PlayerLocation()
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPlayerCanUseExitToMoveBetweenRooms(t *testing.T) {
	t.Parallel()
	data := map[string]*rage.Entity{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination: "Bedroom",
				},
			},
			Contents: []string{"player"},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Kind:        "Room",
		},
		"player": {
			Name:     "player",
			Kind:     "Character",
			Location: "Living Room",
		},
	}
	game, err := rage.NewGame(data, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	cmd := rage.Command{Action: "north"}
	start := game.PlayerLocation()
	err = game.Do(cmd)
	finish := game.PlayerLocation()
	if err != nil {
		t.Fatal(err)
	}

	if start == finish {
		t.Error("player did not move rooms")
	}
}

func TestPlayerCannotPassExitWithoutKey(t *testing.T) {
	t.Parallel()
	data := map[string]*rage.Entity{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination:    "Bedroom",
					Requires:       "key",
					FailureMessage: "The door appears to be locked.",
				},
			},
			Contents: []string{"key", "player"},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Kind:        "Room",
		},
		"key": {
			Name:     "key",
			Location: "Living Room",
			Kind:     "Thing",
		},
		"player": {
			Name:     "player",
			Location: "Living Room",
			Kind:     "Character",
		},
	}

	writer := bytes.Buffer{}
	game, err := rage.NewGame(data, "player", &writer)
	if err != nil {
		t.Fatal(err)
	}

	cmd := rage.Command{Action: "north"}
	start := game.PlayerLocation()
	_ = game.Do(cmd)
	finish := game.PlayerLocation()

	want := "The door appears to be locked."
	got := strings.TrimSpace(writer.String())
	if got != want {
		t.Error(cmp.Diff(want, got))
	}

	if start != finish {
		t.Error("player was able to move through exit without key")
	}
}

func TestPlayerCannotPassExitWithoutKeyAndSeesDefaultFailureMessage(t *testing.T) {
	t.Parallel()
	data := map[string]*rage.Entity{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination: "Bedroom",
					Requires:    "key",
				},
			},
			Contents: []string{"key", "player"},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Kind:        "Room",
		},
		"key": {
			Name:     "key",
			Location: "Living Room",
			Kind:     "Thing",
		},
		"player": {
			Name:     "player",
			Location: "Living Room",
			Kind:     "Character",
		},
	}
	writer := bytes.Buffer{}
	game, err := rage.NewGame(data, "player", &writer)
	if err != nil {
		t.Fatal(err)
	}

	cmd := rage.Command{Action: "north"}
	start := game.PlayerLocation()
	_ = game.Do(cmd)
	finish := game.PlayerLocation()

	want := rage.DefaultExitFailureMessage
	got := strings.TrimSpace(writer.String())
	if got != want {
		t.Error(cmp.Diff(want, got))
	}

	if start != finish {
		t.Error("player was able to move through exit without key")
	}
}
