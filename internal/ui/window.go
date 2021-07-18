package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
	"github.com/tarrsalah/pkt"

	"sort"
)

// The Window model
type Window struct {
	*tview.Application

	items         pkt.Items
	selectedItems pkt.Items
	tags          pkt.Tags
	selectedTags  map[int]struct{}

	itemsTable *itemsTable
	tagsTable  *tagsTable

	widgets       []tview.Primitive
	currentWidget int
}

func (w *Window) nextWidget() {
	next := w.currentWidget + 1
	if next >= len(w.widgets) {
		next = 0
	}

	w.Application.SetFocus(w.widgets[next])
	w.currentWidget = next

}

func (w *Window) handleSelectItem(i, j int) {
	browser.OpenURL(w.selectedItems[i].URL())
}

func (w *Window) handleSelectTag(i, j int) {
	if _, ok := w.selectedTags[i]; ok {
		delete(w.selectedTags, i)
	} else {
		w.selectedTags[i] = struct{}{}
	}

	// TODO: move this logic to the business models
	// Show all pocket items
	if len(w.selectedTags) == 0 {
		w.selectedItems = w.items
		w.Render()
		return
	}

	// Show only pocket items with tag
	w.selectedItems = make([]pkt.Item, 0)

	for _, item := range w.items {
		isTagged := false

		for _, tag := range item.Tags() {
			for i := range w.selectedTags {
				if w.tags[i].Label == tag.Label {
					isTagged = true
					break
				}
			}
			if isTagged {
				break
			}
		}

		if isTagged {
			w.selectedItems = append(w.selectedItems, item)
		}
	}
	w.Render()
}

// NewWindow returns a new UI window
func NewWindow(items pkt.Items) *Window {
	w := &Window{
		Application:  tview.NewApplication(),
		selectedTags: make(map[int]struct{}),
	}

	w.items = items
	w.selectedItems = w.items
	w.tags = w.items.Tags()
	sort.Sort(w.tags)

	w.itemsTable = newItemsTable()
	w.tagsTable = newTagsTable()

	w.itemsTable.handleSelect = w.handleSelectItem
	w.tagsTable.handleSelect = w.handleSelectTag

	w.widgets = append(w.widgets, w.itemsTable)
	w.widgets = append(w.widgets, w.tagsTable)

	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(w.tagsTable, 0, 1, true).
		AddItem(w.itemsTable, 0, 5, true)

	pages := tview.NewPages().
		AddAndSwitchToPage("main", flex, true)

	w.SetRoot(pages, true).EnableMouse(true)

	w.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			w.nextWidget()
			return event
		default:
			return event
		}
	})

	w.Render()
	w.Application.SetFocus(w.widgets[w.currentWidget])

	return w
}

// Render the children
func (w *Window) Render() {
	w.itemsTable.items = w.selectedItems

	w.tagsTable.tags = w.tags
	w.tagsTable.selectedTags = w.selectedTags

	w.itemsTable.Render()
	w.tagsTable.Render()
}
