package rage_test

import (
	"github.com/faiuwle/go-parser-game/rage"
	"golang.org/x/exp/slices"
	"testing"
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
		Exits: map[string]*rage.Entity{
			"north": nil,
			"east":  nil,
		},
	}

	exitString := rage.ListExits(room)
	want := "You can go north and east from here."

	if exitString != want {
		t.Errorf("Wanted %q, got %q", want, exitString)
	}
}

func TestListItems(t *testing.T) {
	entities := []*rage.Entity{
		{
			Name: "key",
		},
		{
			Name: "phone",
		},
		{
			Name: "chocolate",
		},
	}

	room := rage.Entity{
		Contents: []int{0, 1, 2},
	}

	exitString := rage.ListItems(room, entities)
	want := "You see here key, phone, and chocolate."

	if exitString != want {
		t.Errorf("Wanted %q, got %q", want, exitString)
	}
}

func TestTakeItem(t *testing.T) {
	entities := []*rage.Entity{
		{
			Name: "Player",
		},
		{
			Name: "key",
		},
		{
			Name: "phone",
		},
		{
			Name:     "Room",
			Contents: []int{0, 1},
		},
	}

	response := rage.TakeItem(entities, entities[0], entities[3], "key")
	want := "You take the key."

	if response != want {
		t.Errorf("Wanted %q, got %q", want, response)
	}

	if !slices.Contains(entities[0].Contents, 1) {
		t.Errorf("Key was not transferred to inventory.")
	}

	if slices.Contains(entities[3].Contents, 1) {
		t.Errorf("Key is still in the room.")
	}

	response = rage.TakeItem(entities, entities[0], entities[3], "phone")
	want = "I can't see that here."

	if response != want {
		t.Errorf("Wanted %q, got %q", want, response)
	}

	if slices.Contains(entities[0].Contents, 2) {
		t.Errorf("Player took the phone.")
	}
}
