package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tagsTable struct {
	*tview.Table
	model *model
}

func newTagsTable(model *model) *tagsTable {
	t := &tagsTable{
		Table: tview.NewTable(),
		model: model,
	}

	t.SetBorder(true)
	t.SetSelectable(true, false)

	t.refresh()
	return t
}

func (t *tagsTable) refresh() {
	t.Clear()
	for i, tag := range t.model.tags {
		cell := tview.NewTableCell(tag.Label)

		if t.model.isTagSelected(tag) {
			cell.SetTextColor(tcell.ColorYellow)
			cell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
		}
		t.Table.SetCell(i, 0, cell)

	}
	t.SetTitle(t.title())
}

func (t *tagsTable) title() string {
	l := len(t.model.tags)
	if l > 0 {
		return fmt.Sprintf("Tags (%d)", len(t.model.tags))
	}

	return "Tags"
}
