package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// showTool renders a tool's screen. The full select -> options -> output -> run
// flow is built in a later milestone; this placeholder keeps navigation working.
func (u *UI) showTool(t tool.Tool) {
	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() { u.showCategory(t.Category()) })
	back.Importance = widget.LowImportance

	title := widget.NewLabelWithStyle(t.Name(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.SizeName = theme.SizeNameHeadingText
	desc := widget.NewLabel(t.Description())
	desc.Wrapping = fyne.TextWrapWord

	header := container.NewVBox(container.NewHBox(back), title, widget.NewSeparator())
	body := container.NewCenter(widget.NewLabel("This tool's screen is coming soon."))
	u.setContent(container.NewPadded(container.NewBorder(header, nil, nil, nil, container.NewVBox(desc, body))))
}
