package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/tarrsalah/pkt"
)

type itemsTable struct {
	*tview.Table
	items        pkt.Items
	handleSelect func(int, int)
}

func newItemsTable() *itemsTable {
	t := &itemsTable{
		Table: tview.NewTable(),
	}

	t.SetBorder(true)
	t.SetSelectable(true, false)

	return t
}

func (t *itemsTable) Render() {
	t.Clear()
	for i, item := range t.items {
		tags := fmt.Sprintf("[yellow] %s", item.Tags())
		title := fmt.Sprintf("%d. %s [green](%s)",
			i+1, item.Title(),
			item.Host())

		t.SetCell(i, 0, tview.NewTableCell(title).
			SetMaxWidth(1).
			SetExpansion(3))

		t.SetCell(i, 1, tview.NewTableCell(tags).
			SetMaxWidth(1).
			SetExpansion(2))
	}

	t.SetSelectedFunc(t.handleSelect)
	t.SetTitle(t.title())
}

func (t *itemsTable) title() string {
	count := len(t.items)
	if count > 0 {
		return fmt.Sprintf("Pocket items (%d)", count)
	}
	return "Pocket items"
}
