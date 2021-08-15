package ui

import (
	"github.com/tarrsalah/pkt"
)

type tags struct {
	all      pkt.Tags
	selected map[string]struct{}
}

func newTags(l pkt.Tags) *tags {
	return &tags{
		all:      l,
		selected: make(map[string]struct{}),
	}
}

func (t *tags) toggle(tag pkt.Tag) {
	if _, ok := t.selected[tag.Label]; ok {
		delete(t.selected, tag.Label)
	} else {
		t.selected[tag.Label] = struct{}{}
	}
}

func (t *tags) isSelected(tag pkt.Tag) bool {
	if _, ok := t.selected[tag.Label]; ok {
		return true
	}

	return false
}

func (t *tags) getSelected() pkt.Tags {
	selected := []pkt.Tag{}
	for _, tag := range t.all {
		if t.isSelected((tag)) {
			selected = append(selected, tag)
		}
	}

	return selected
}

func (t *tags) get(index int) pkt.Tag {
	return t.all[index]
}

type items struct {
	all      pkt.Items
	selected pkt.Items
}

func newItems(i pkt.Items) *items {
	return &items{
		all:      i,
		selected: i,
	}
}

func (i *items) filter(tags *tags) {
	i.selected = i.all.GetTagged(tags.getSelected())
}

func (i *items) getSelected() pkt.Items {
	return i.selected
}

func (i *items) get(index int) pkt.Item {
	return i.selected[index]
}
