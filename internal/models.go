package internal

import (
	"fmt"
	"net/url"
	"sort"
)

// PageCount is the default response page count
const PageCount = 100

// Tag is a packet tag
type Tag struct {
	ID    string `json:"item_id"`
	Label string `json:"tag"`
}

func (t Tag) String() string {
	return t.Label
}

// Tags is a set of Tags
type Tags []Tag

// Len return the length of a list of tags
func (t Tags) Len() int {
	return len(t)
}

// Swap swaps two tags
func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less compares two tags by label
func (t Tags) Less(i, j int) bool {
	return t[i].Label < t[j].Label
}

func (t Tags) String() string {
	s := ""
	for _, tag := range t {
		s += fmt.Sprintf("(%s)", tag)
	}

	return s
}

// Item is the poeckt item
type Item struct {
	ID            string         `json:"item_id"`
	GivenURL      string         `json:"given_url"`
	GivenTitle    string         `json:"given_title"`
	ResolvedTitle string         `json:"resolved_title"`
	AddedAt       string         `json:"time_added"`
	TagsMap       map[string]Tag `json:"tags"`
}

// Title gets the item's title
func (item Item) Title() string {
	if len(item.ResolvedTitle) > 3 {
		return item.ResolvedTitle
	}

	return item.GivenTitle
}

// URL returns the item's given URL
func (item Item) URL() string {
	return item.GivenURL
}

// Host returns the item's given URL host
func (item Item) Host() string {
	itemURL, _ := url.Parse(item.GivenURL)
	return itemURL.Host
}

// Tags returns the list of the item's tag
func (item Item) Tags() Tags {
	tags := []Tag{}
	for _, tag := range item.TagsMap {
		tags = append(tags, tag)
	}

	return tags
}

// Items is a list of items
type Items []Item

// Tags return a list of tags from a list of items
func (items Items) Tags() Tags {
	tags := Tags([]Tag{})
	tagsMap := map[string]Tag{}

	_ = tagsMap

	for _, item := range items {
		for _, tag := range item.Tags() {
			tagsMap[tag.Label] = tag
		}
	}

	for _, tag := range tagsMap {
		tags = append(tags, tag)
	}

	sort.Sort(tags)
	return tags
}

// GetTagged filter items by tags
func (items Items) GetTagged(tags Tags) Items {
	if len(tags) == 0 {
		taggedItems := make([]Item, len(items))
		copy(taggedItems, items)
		return taggedItems
	}

	taggedItemes := []Item{}
	for _, item := range items {
		isTagged := false
		for _, tag := range item.Tags() {
			for _, filter := range tags {
				if filter.Label == tag.Label {
					isTagged = true
					break
				}
			}

			if isTagged {
				break
			}
		}

		if isTagged {
			taggedItemes = append(taggedItemes, item)
		}
	}

	return taggedItemes

}

func (item Item) isTagged(tags Tags) bool {
	for _, tag := range item.Tags() {
		for _, filter := range tags {
			if filter.Label == tag.Label {
				return true
			}
		}
	}
	return false
}
