package ui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed appicon.png
var iconBytes []byte

// appIcon is the embedded application icon used for the window and taskbar.
var appIcon = fyne.NewStaticResource("appicon.png", iconBytes)

// appTheme is the default Fyne theme with a teal accent to match the branding.
type appTheme struct{ fyne.Theme }

func newAppTheme() fyne.Theme { return &appTheme{theme.DefaultTheme()} }

func (t *appTheme) Color(name fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary, theme.ColorNameHyperlink, theme.ColorNameFocus:
		return color.NRGBA{R: 0x14, G: 0x8f, B: 0x8f, A: 0xff}
	}
	return t.Theme.Color(name, v)
}
