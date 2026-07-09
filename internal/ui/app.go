// Package ui builds the Fyne desktop interface. It renders entirely from the
// tool.Registry, so it never needs editing when new tools are added.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// AppName is shown in the title bar and about text.
const AppName = "File Tools"

// UI owns the Fyne app + main window and drives navigation between screens.
type UI struct {
	fyne     fyne.App
	win      fyne.Window
	registry *tool.Registry
}

// New creates the application shell backed by the given registry.
func New(registry *tool.Registry) *UI {
	a := app.NewWithID("ai.filetools.desktop")
	a.Settings().SetTheme(theme.DefaultTheme())

	w := a.NewWindow(AppName)
	w.Resize(fyne.NewSize(880, 620))
	w.CenterOnScreen()

	return &UI{fyne: a, win: w, registry: registry}
}

// Run shows the home screen and starts the event loop (blocks until close).
func (u *UI) Run() {
	u.showHome()
	u.win.ShowAndRun()
}

// setContent swaps the whole window body.
func (u *UI) setContent(c fyne.CanvasObject) {
	u.win.SetContent(c)
}
