package pkt

import "net/url"

const PAGE_COUNT = 100

type Tag struct {
	Id    string `json:"item_id"`
	Label string `json:"tag"`
}

type Store interface {
	Get() []Item
	Put([]Item)
	Close()
}

type Item struct {
	Id            string         `json:"item_id"`
	GivenUrl      string         `json:"given_url"`
	GivenTitle    string         `json:"given_title"`
	ResolvedTitle string         `json:"resolved_title"`
	AddedAt       string         `json:"time_added"`
	TagsMap       map[string]Tag `json:"tags"`
}

func (item Item) Title() string {
	if len(item.ResolvedTitle) > 3 {
		return item.ResolvedTitle
	}

	return item.GivenTitle
}

func (item Item) Url() string {
	return item.GivenUrl
}

func (item Item) Host() string {
	itemUrl, _ := url.Parse(item.GivenUrl)
	return itemUrl.Host
}

func (item Item) Tags() []string {
	tags := []string{}
	for _, tag := range item.TagsMap {
		tags = append(tags, tag.Label)
	}

	return tags
}
