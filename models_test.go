package pkt

import (
	"testing"
)

func TestGetTaggedItems(t *testing.T) {
	items := []Item{
		{
			TagsMap: map[string]Tag{
				"elixir": {Label: "elixir"},
			},
		},

		{
			TagsMap: map[string]Tag{
				"elixir": {Label: "elixir"},
			},
		},

		{
			TagsMap: map[string]Tag{
				"ruby":   {Label: "ruby"},
				"python": {Label: "python"},
			},
		},

		{
			TagsMap: map[string]Tag{
				"elixir": {Label: "elixir"},
				"golang": {Label: "golang"},
				"python": {Label: "python"},
			},
		},
	}

	tests := []struct {
		tags   []string
		result int
	}{
		{[]string{"elixir"}, 3},
		{[]string{"python"}, 2},
		{[]string{"ruby"}, 1},
		{[]string{"golang"}, 1},
	}

	for _, tc := range tests {
		filteredItems := getTaggedItems(items, tc.tags)
		got := len(filteredItems)

		if len(filteredItems) != tc.result {
			t.Errorf("Expected %d item(s) with %v filters, got %d", tc.result, tc.tags, got)
		}
	}
}
