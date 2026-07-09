package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// categoryIcon maps a category to a built-in Fyne icon.
func categoryIcon(c tool.Category) fyne.Resource {
	switch c {
	case tool.CategoryPDF:
		return theme.DocumentIcon()
	case tool.CategoryImage:
		return theme.MediaPhotoIcon()
	case tool.CategoryDoc:
		return theme.FileTextIcon()
	case tool.CategoryBatch:
		return theme.ContentCopyIcon()
	default:
		return theme.FileIcon()
	}
}

// showHome renders the four big category cards.
func (u *UI) showHome() {
	title := widget.NewLabelWithStyle(AppName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.SizeName = theme.SizeNameHeadingText
	subtitle := widget.NewLabel("Simple, offline file conversion. Pick a tool to begin.")

	counts := u.registry.Counts()
	var cards []fyne.CanvasObject
	for _, c := range tool.Categories {
		c := c
		n := counts[c]
		sub := describeCategory(c, n)
		cards = append(cards, newTappableCard(categoryIcon(c), string(c), sub, func() {
			u.showCategory(c)
		}))
	}

	grid := container.NewGridWithColumns(2, cards...)
	header := container.NewVBox(title, subtitle, widget.NewSeparator())
	body := container.NewBorder(header, nil, nil, nil, container.NewPadded(grid))
	u.setContent(container.NewPadded(body))
}

func describeCategory(c tool.Category, n int) string {
	tools := "tools"
	if n == 1 {
		tools = "tool"
	}
	switch c {
	case tool.CategoryPDF:
		return fmt.Sprintf("Merge, split, compress, rotate and more — %d %s", n, tools)
	case tool.CategoryImage:
		return fmt.Sprintf("Convert, resize, compress and crop images — %d %s", n, tools)
	case tool.CategoryDoc:
		return fmt.Sprintf("Word, PDF and text documents — %d %s", n, tools)
	case tool.CategoryBatch:
		return fmt.Sprintf("Process many files at once — %d %s", n, tools)
	default:
		return fmt.Sprintf("%d %s", n, tools)
	}
}
