package rage_test

import (
	"github.com/faiuwle/go-parser-game/rage"
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
