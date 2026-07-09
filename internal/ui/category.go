package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// showCategory lists every tool within a category as a clickable card.
func (u *UI) showCategory(c tool.Category) {
	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() { u.showHome() })
	back.Importance = widget.LowImportance

	title := widget.NewLabelWithStyle(string(c), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.SizeName = theme.SizeNameHeadingText

	header := container.NewVBox(
		container.NewHBox(back),
		title,
		widget.NewSeparator(),
	)

	tools := u.registry.ByCategory(c)
	var body fyne.CanvasObject
	if len(tools) == 0 {
		body = container.NewCenter(widget.NewLabel("Tools coming soon."))
	} else {
		var cards []fyne.CanvasObject
		for _, t := range tools {
			t := t
			cards = append(cards, newTappableCard(categoryIcon(c), t.Name(), t.Description(), func() {
				u.showTool(t)
			}))
		}
		body = container.NewVScroll(container.NewPadded(container.NewVBox(cards...)))
	}

	u.setContent(container.NewPadded(container.NewBorder(header, nil, nil, nil, body)))
}
