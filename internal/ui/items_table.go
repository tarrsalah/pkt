package ui

import (
	"fmt"
	"github.com/rivo/tview"
)

type itemsTable struct {
	*tview.Table
	model *model
}

func newItemsTable(model *model) *itemsTable {
	t := &itemsTable{
		Table: tview.NewTable(),
		model: model,
	}

	t.SetBorder(true)
	t.SetSelectable(true, false)
	t.refresh()

	return t
}

func (t *itemsTable) refresh() {
	t.Clear()

	for i, item := range t.model.selecedItems {
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

	t.SetTitle(t.title())
}

func (t *itemsTable) title() string {
	count := len(t.model.selecedItems)
	if count > 0 {
		return fmt.Sprintf("Pocket items (%d)", count)
	}
	return "Pocket items"
}
