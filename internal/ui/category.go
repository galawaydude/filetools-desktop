package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// screenHeader builds a back button + accent badge + title/subtitle row reused
// by the category and tool screens.
func screenHeader(icon fyne.Resource, accent fyne.ThemeColorName, title, subtitle string, onBack func()) fyne.CanvasObject {
	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), onBack)
	back.Importance = widget.LowImportance

	titleText := canvas.NewText(title, themeColor(theme.ColorNameForeground))
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 22

	var titleBlock fyne.CanvasObject
	if subtitle != "" {
		titleBlock = container.NewVBox(titleText, mutedText(subtitle))
	} else {
		titleBlock = container.NewVBox(titleText)
	}

	row := container.NewBorder(nil, nil, vcenter(accentBadge(icon, accent, 44)), nil, vcenter(titleBlock))
	return container.NewVBox(container.NewHBox(back), container.NewPadded(row), widget.NewSeparator())
}

// showCategory lists every tool within a category as a clickable card.
func (u *UI) showCategory(c tool.Category) {
	accent := categoryColor(c)
	header := screenHeader(categoryIcon(c), accent, string(c), "Choose a tool", func() { u.showHome() })

	tools := u.registry.ByCategory(c)
	var body fyne.CanvasObject
	if len(tools) == 0 {
		body = container.NewCenter(mutedText("Tools coming soon."))
	} else {
		var cards []fyne.CanvasObject
		for _, t := range tools {
			t := t
			cards = append(cards, newHoverCard(toolIcon(t), accent, t.Name(), t.Description(), func() {
				u.showTool(t)
			}))
		}
		body = container.NewVScroll(container.NewPadded(container.NewVBox(cards...)))
	}

	u.setContent(container.NewPadded(container.NewBorder(header, nil, nil, nil, body)))
}
