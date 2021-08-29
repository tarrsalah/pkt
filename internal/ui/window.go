package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
	"github.com/tarrsalah/pkt"
)

// The UI model
type model struct {
	tags         pkt.Tags
	selectedTags map[string]struct{}
	items        pkt.Items
	selecedItems pkt.Items
}

func newModel(items pkt.Items) *model {
	return &model{
		items:        items,
		selecedItems: items,
		tags:         items.Tags(),
		selectedTags: make(map[string]struct{}),
	}
}

func (m *model) getItem(i int) pkt.Item {
	return m.selecedItems[i]
}

func (m *model) getTag(i int) pkt.Tag {
	return m.tags[i]
}

func (m *model) toggleTag(tag pkt.Tag) {
	if _, ok := m.selectedTags[tag.Label]; ok {
		delete(m.selectedTags, tag.Label)
	} else {
		m.selectedTags[tag.Label] = struct{}{}
	}
	m.selecedItems = m.items.GetTagged(m.getSelectedTags())
}

func (m *model) isTagSelected(tag pkt.Tag) bool {
	_, isSelected := m.selectedTags[tag.Label]
	return isSelected
}

func (m *model) getSelectedTags() pkt.Tags {
	selected := []pkt.Tag{}
	for _, tag := range m.tags {
		if m.isTagSelected((tag)) {
			selected = append(selected, tag)
		}
	}

	return selected
}

// Window is the global ui component
type Window struct {
	*tview.Application
	model *model

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
	browser.OpenURL(w.model.getItem(i).URL())
}

func (w *Window) handleSelectTag(i, _ int) {
	w.model.toggleTag(w.model.getTag(i))

	w.itemsTable.refresh()
	w.tagsTable.refresh()
}

// NewWindow returns a new UI window
func NewWindow(items pkt.Items) *Window {
	w := &Window{
		Application: tview.NewApplication(),
		model:       newModel(items),
	}

	w.itemsTable = newItemsTable(w.model)
	w.tagsTable = newTagsTable(w.model)

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
