package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
	"github.com/tarrsalah/pkt"
)

func draw(items []pkt.Item) {
	title := fmt.Sprintf("Pocket items (%d)", len(items))
	itemsTable := tview.NewTable().SetSelectable(true, false).
		Select(0, 0).SetFixed(1, 1)

	itemsTable.SetTitle(title).SetTitleAlign(tview.AlignLeft)
	itemsTable.SetBorder(true)

	headers := []string{
		"Title",
		"tags",
	}

	for i, header := range headers {
		itemsTable.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorBlack,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, item := range items {
		tags := ""
		if len(item.Tags()) > 0 {
			tags = "[yellow]"
			for _, tag := range item.Tags() {
				tags += fmt.Sprintf("(%s)", tag)
			}
		}
		title := fmt.Sprintf("%d. %s [green](%s)",
			i+1, item.Title(),
			item.Host())

		itemsTable.SetCell(i+1, 0, tview.NewTableCell(title).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(3))

		itemsTable.SetCell(i+1, 1, tview.NewTableCell(tags).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(2))
	}

	itemsTable.SetSelectedFunc(func(row, column int) {
		current := items[row-1]
		browser.OpenURL(current.Url())
	})

	app := tview.NewApplication()
	grid := tview.NewGrid().SetRows(1).
		AddItem(itemsTable, 0, 0, 2, 1, 0, 0, true)

	pages := tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	app.SetRoot(pages, true)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
