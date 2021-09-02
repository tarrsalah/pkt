package internal

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
)

type App struct {
	ui            *tview.Application
	widgets       []tview.Primitive
	pages         *tview.Pages
	currentWidget int

	list struct {
		table    *tview.Table
		all      Items
		selected Items
	}

	tags struct {
		table    *tview.Table
		all      Tags
		selected map[string]struct{}
	}
}

// Get selected Tags
func (app *App) getSelectedTags() Tags {
	selected := []Tag{}
	for _, tag := range app.tags.all {
		if _, ok := app.tags.selected[tag.Label]; ok {
			selected = append(selected, tag)
		}
	}
	return selected
}

func (app *App) nextWidget() {
	next := app.currentWidget + 1
	if next >= len(app.widgets) {
		next = 0
	}

	app.ui.SetFocus(app.widgets[next])
	app.currentWidget = next
}

func (app *App) show() {
	// show list
	app.list.table.Clear()
	selectedTags := app.getSelectedTags()

	app.list.selected = []Item{}
	for _, l := range app.list.all {
		if len(selectedTags) == 0 || l.isTagged(selectedTags) {
			app.list.selected = append(app.list.selected, l)
		}
	}

	for i, l := range app.list.selected {
		tags := fmt.Sprintf("[yellow] %s", l.Tags())
		title := fmt.Sprintf("%d. %s [green](%s)",
			i+1, l.Title(),
			l.Host())

		app.list.table.SetCell(i, 0, tview.NewTableCell(title).
			SetMaxWidth(1).
			SetExpansion(3))

		app.list.table.SetCell(i, 1, tview.NewTableCell(tags).
			SetMaxWidth(1).
			SetExpansion(2))
	}

	app.list.table.SetTitle(fmt.Sprintf("list (%d)", len(app.list.selected)))
	app.list.table.Select(0, 0)

	// show tags
	app.tags.table.Clear()
	for i, tag := range app.tags.all {
		cell := tview.NewTableCell(tag.Label)

		if _, ok := app.tags.selected[tag.Label]; ok {
			cell.SetTextColor(tcell.ColorYellow)
			cell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
		}
		app.tags.table.SetCell(i, 0, cell)
	}
	app.tags.table.SetTitle(fmt.Sprintf("tags (%d)", len(app.tags.all)))
}

func (app *App) handleSelectItem(i, j int) {
	browser.OpenURL(app.list.selected[i].URL())
}

func (app *App) handleSelectTag(i, j int) {
	tag := app.tags.all[i]
	if _, ok := app.tags.selected[tag.Label]; ok {
		delete(app.tags.selected, tag.Label)
	} else {
		app.tags.selected[tag.Label] = struct{}{}
	}
	app.show()
}

func NewApp(list Items) *App {
	app := &App{}

	app.ui = tview.NewApplication()

	// list
	app.list.all = list
	app.list.table = tview.NewTable()
	app.list.table.SetBorder(true)
	app.list.table.SetSelectable(true, false)

	// tags
	app.tags.all = list.Tags()
	app.tags.selected = make(map[string]struct{})
	app.tags.table = tview.NewTable()
	app.tags.table.SetBorder(true)
	app.tags.table.SetSelectable(true, false)

	app.list.table.SetSelectedFunc(app.handleSelectItem)
	app.tags.table.SetSelectedFunc(app.handleSelectTag)

	app.widgets = append(app.widgets, app.list.table)
	app.widgets = append(app.widgets, app.tags.table)

	page := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(app.tags.table, 0, 1, true).
		AddItem(app.list.table, 0, 5, true)

	app.pages = tview.NewPages().AddAndSwitchToPage("main", page, true)
	app.ui.SetRoot(app.pages, true).EnableMouse(true)

	app.ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.nextWidget()
		}
		return event
	})

	app.ui.SetFocus(app.widgets[app.currentWidget])

	return app
}

func (app *App) Run() {
	app.show()
	app.ui.Run()
}
