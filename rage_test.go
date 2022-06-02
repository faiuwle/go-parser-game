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

var commonGameData = rage.GameData{
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

func TestTakeItemSucceedsIfItemIsPresent(t *testing.T) {
	entities := commonGameData
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
	if g.GetEntity("Room").Contains("key") {
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
	bedroom := &rage.Entity{
		Name:        "Bedroom",
		Description: "The bedroom",
		Kind:        "Room",
		Contents:    []string{"player"},
	}
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
		"Bedroom": bedroom,
		"player": {
			Name:     "player",
			Location: "Bedroom",
		},
	}
	game, err := rage.NewGame(data, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	want := bedroom
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

func TestGetEntityPanicsWhenGivenNonExistentEntity(t *testing.T) {
	data := map[string]*rage.Entity{}
	game, err := rage.NewGame(data, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("GetEntity did not panic")
		}
	}()
	game.GetEntity("whatever")
}

func TestGetEntityReturnsEntityWhenItExists(t *testing.T) {
	want := &rage.Entity{Name: "entity"}
	data := map[string]*rage.Entity{
		want.Name: want,
	}
	game, err := rage.NewGame(data, "player", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	entity := game.GetEntity(want.Name)
	if entity != want {
		t.Fatalf("Found %#v instead of %#v", entity, want)
	}
}

func TestNewGameReturnsErrorWithInconsistentData(t *testing.T) {
	testCases := map[string]rage.GameData{
		"non-existent location": {
			"object": &rage.Entity{
				Name:     "object",
				Location: "non-existent",
				Kind:     "Thing",
			},
			"player": {
				Name: "player",
			},
		},
		"non-existant object in room": {
			"room": &rage.Entity{
				Name: "room",
				Contents: []string{
					"non-existent",
				},
				Kind: "Room",
			},
			"player": {
				Name: "player",
			},
		},
		"empty location on thing": {
			"object": &rage.Entity{
				Name:     "object",
				Location: "",
				Kind:     "Thing",
			},
			"player": {
				Name: "player",
			},
		},
		"non-empty location on room": {
			"room": &rage.Entity{
				Name:     "room",
				Location: "somewhere",
				Kind:     "Room",
			},
			"somewhere": &rage.Entity{
				Name:     "somewhere",
				Location: "room",
				Kind:     "Room",
			},
			"player": {
				Name: "player",
			},
		},
		"no player in game data": {},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := rage.NewGame(tc, "player", io.Discard)
			if err == nil {
				t.Fatal("Did not fail with invalid game data")
			}
		})
	}
}

func TestNewGameCreatesGameFromConsistentData(t *testing.T) {
	data := map[string]*rage.Entity{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination:    "Bedroom",
					Requires:       "key",
					FailureMessage: "You rattle the handle, but the door seems locked.",
				},
			},
			Contents: []string{"key", "Shera"},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Exits: map[string]rage.Exit{
				"south": {
					Destination: "Living Room",
				},
			},
			Kind: "Room",
		},
		"key": {
			Name:     "key",
			Location: "Living Room",
			Kind:     "Thing",
		},
		"Shera": {
			Name:     "Shera",
			Location: "Living Room",
		},
	}

	_, err := rage.NewGame(data, "Shera", io.Discard)
	if err != nil {
		t.Fatalf("Errored on consistent game data: %s", err)
	}
}
