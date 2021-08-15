package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tagsTable struct {
	*tview.Table
	tags         *tags
	handleSelect func(int, int)
}

func newTagsTable() *tagsTable {
	t := &tagsTable{
		Table: tview.NewTable(),
	}

	t.SetBorder(true)
	t.SetSelectable(true, false)

	return t
}

func (t *tagsTable) Render() {
	t.Clear()
	for i, tag := range t.tags.all {
		cell := tview.NewTableCell(tag.Label)

		if t.tags.isSelected(tag) {
			cell.SetTextColor(tcell.ColorYellow)
			cell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
		}
		t.Table.SetCell(i, 0, cell)

	}
	t.SetSelectedFunc(t.handleSelect)
	t.SetTitle(t.title())
}

func (t *tagsTable) title() string {
	l := len(t.tags.all)
	if l > 0 {
		return fmt.Sprintf("Tags (%d)", len(t.tags.all))
	}

	return "Tags"
}
