package pkt

import (
	"testing"
)

func TestGetTaggedItems(t *testing.T) {
	var items Items
	items = []Item{
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
		tags   Tags
		result int
	}{
		{[]Tag{{Label: "elixir"}}, 3},
		{[]Tag{{Label: "python"}}, 2},
		{[]Tag{{Label: "ruby"}}, 1},
		{[]Tag{{Label: "golang"}}, 1},
	}

	for _, tc := range tests {
		filteredItems := Items(items).GetTagged(tc.tags)
		got := len(filteredItems)

		if len(filteredItems) != tc.result {
			t.Errorf("Expected %d item(s) with %v filters, got %d", tc.result, tc.tags, got)
		}
	}
}
