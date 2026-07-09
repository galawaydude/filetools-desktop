package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// appHeader is the branded title block shown at the top of the home screen.
func appHeader() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(appIcon)
	logo.FillMode = canvas.ImageFillContain
	logoBox := container.NewGridWrap(fyne.NewSize(52, 52), logo)

	title := canvas.NewText(AppName, themeColor(theme.ColorNameForeground))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 28
	subtitle := mutedText("Simple, offline file conversion — pick a tool to begin.")

	return container.NewBorder(nil, nil, container.NewPadded(logoBox), nil,
		vcenter(container.NewVBox(title, subtitle)),
	)
}

// showHome renders the four category cards in a responsive two-column grid.
func (u *UI) showHome() {
	counts := u.registry.Counts()
	var cards []fyne.CanvasObject
	for _, c := range tool.Categories {
		c := c
		cards = append(cards, newHoverCard(
			categoryIcon(c), categoryColor(c),
			string(c), describeCategory(c, counts[c]),
			func() { u.showCategory(c) },
		))
	}

	grid := container.NewGridWithColumns(2, cards...)
	body := container.NewBorder(
		container.NewVBox(container.NewPadded(appHeader()), widget.NewSeparator()),
		nil, nil, nil,
		container.NewVScroll(container.NewPadded(container.NewVBox(grid))),
	)
	u.setContent(container.NewPadded(body))
}

func describeCategory(c tool.Category, n int) string {
	tools := "tools"
	if n == 1 {
		tools = "tool"
	}
	switch c {
	case tool.CategoryPDF:
		return fmt.Sprintf("Merge, split, compress, rotate and more · %d %s", n, tools)
	case tool.CategoryImage:
		return fmt.Sprintf("Convert, resize, compress and crop · %d %s", n, tools)
	case tool.CategoryDoc:
		return fmt.Sprintf("Word, PDF and text documents · %d %s", n, tools)
	case tool.CategoryBatch:
		return fmt.Sprintf("Process many files at once · %d %s", n, tools)
	default:
		return fmt.Sprintf("%d %s", n, tools)
	}
}
