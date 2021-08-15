package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
	"github.com/tarrsalah/pkt"
)

// Window is the global ui component
type Window struct {
	*tview.Application

	items *items
	tags  *tags

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
	browser.OpenURL(w.items.get(i).URL())
}

func (w *Window) handleSelectTag(i, _ int) {
	w.tags.toggle(w.tags.get(i))
	w.items.filter(w.tags)

	w.itemsTable.refresh()
	w.tagsTable.refresh()
}

// NewWindow returns a new UI window
func NewWindow(items pkt.Items) *Window {
	w := &Window{
		Application: tview.NewApplication(),
		items:       newItems(items),
		tags:        newTags(items.Tags()),
	}

	w.itemsTable = newItemsTable(w.items)
	w.tagsTable = newTagsTable(w.tags)

	w.itemsTable.SetSelectedFunc(w.handleSelectItem)
	w.tagsTable.SetSelectedFunc(w.handleSelectTag)

	w.widgets = append(w.widgets, w.itemsTable)
	w.widgets = append(w.widgets, w.tagsTable)

	page := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(w.tagsTable, 0, 1, true).
		AddItem(w.itemsTable, 0, 5, true)

	pages := tview.NewPages().
		AddAndSwitchToPage("main", page, true)

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
	w.Application.SetFocus(w.widgets[w.currentWidget])
	return w
}
