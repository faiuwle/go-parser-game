package rage_test

import (
	"testing"

	"github.com/faiuwle/go-parser-game/rage"
	"github.com/google/go-cmp/cmp"
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
	want := "You can go north and east from here."

	if exitString != want {
		t.Errorf("Wanted %q, got %q", want, exitString)
	}
}

func TestListItems(t *testing.T) {
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
			Contents: []string{"key", "phone"},
		},
		"key": {
			Name:     "key",
			Location: "Room",
		},
		"phone": {
			Name:     "phone",
			Location: "Room",
		},
	}
	g, err := rage.NewGame(entities, "Room")
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
			Contents: []string{"key"},
		},
		"key": {
			Name:     "key",
			Location: "Room",
		},
		"phone": {
			Name: "phone",
		},
	}
	g, err := rage.NewGame(entities, "Room")
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
					Destination: 1,
				},
			},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Kind:        "Room",
		},
	}
	game, err := rage.NewGame(data, "Bedroom")
	if err != nil {
		t.Fatal(err)
	}
	want := "Bedroom"
	got := game.PlayerLocation()
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

// func TestPlayerCanUseExitToMoveBetweenRooms(t *testing.T) {
// 	t.Parallel()
// 	data := []*rage.Entity{
// 		{
// 			Id:          0,
// 			Name:        "Living Room",
// 			Description: "The living room",
// 			Kind:        "Room",
// 			Exits: map[string]rage.Exit{
// 				"north": {
// 					Destination: 1,
// 				},
// 			},
// 		},
// 		{
// 			Id:          1,
// 			Name:        "Bedroom",
// 			Description: "The bedroom",
// 			Kind:        "Room",
// 		},
// 		{
// 			Id:       3,
// 			Location: 0,
// 			Kind:     "Character",
// 		},
// 	}
// 	game := rage.NewGame(data)
// 	start := game.PlayerLocation()
// 	_ = game.Do("north")
// 	finish := game.PlayerLocation()
// 	if start == finish {
// 		t.Error("player did not move rooms")
// 	}
// }

// func TestPlayerCannotPassExitWithoutKey(t *testing.T) {
// 	t.Parallel()
// 	data := []*rage.Entity{
// 		{
// 			Id:          0,
// 			Name:        "Living Room",
// 			Description: "The living room",
// 			Kind:        "Room",
// 			Exits: map[string]rage.Exit{
// 				"north": {
// 					Destination: 1,
// 					Requires:    2,
// 				},
// 			},
// 		},
// 		{
// 			Id:          1,
// 			Name:        "Bedroom",
// 			Description: "The bedroom",
// 			Kind:        "Room",
// 		},
// 		{
// 			Id:       2,
// 			Name:     "key",
// 			Location: 1,
// 			Kind:     "Thing",
// 		},
// 		{
// 			Id:       3,
// 			Location: 0,
// 			Kind:     "Character",
// 		},
// 	}
// 	game := rage.NewGame(data)
// 	_ = game.Do("north")
// 	if game.PlayerLocation() != 0 {
// 		t.Error("player was able to move through exit without key")
// 	}
// }
