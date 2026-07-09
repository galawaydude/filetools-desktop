package ui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed appicon.png
var iconBytes []byte

var appIcon = fyne.NewStaticResource("appicon.png", iconBytes)

// Custom theme color names. Standard widgets never ask for these; our own
// components do (via theme.NewColoredResource and direct lookups), so the theme
// stays the single source of truth for the palette.
const (
	colorPDF      fyne.ThemeColorName = "catPDF"
	colorImage    fyne.ThemeColorName = "catImage"
	colorDoc      fyne.ThemeColorName = "catDoc"
	colorBatch    fyne.ThemeColorName = "catBatch"
	colorOnAccent fyne.ThemeColorName = "onAccent"
	colorMuted    fyne.ThemeColorName = "textMuted"
	colorCard     fyne.ThemeColorName = "surfaceCard"
	colorCardHi   fyne.ThemeColorName = "surfaceCardHover"
)

const accentPrimary = 0x12A594 // teal used for buttons, focus and links

type appTheme struct{ fyne.Theme }

func newAppTheme() fyne.Theme { return &appTheme{theme.DefaultTheme()} }

func hex(v uint32) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
}

func alpha(c color.NRGBA, a uint8) color.NRGBA { c.A = a; return c }

func (t *appTheme) Color(name fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	dark := v == theme.VariantDark
	pick := func(light, darkC uint32) color.Color {
		if dark {
			return hex(darkC)
		}
		return hex(light)
	}
	switch name {
	// Category accents (same in both variants — they're saturated enough).
	case colorPDF:
		return hex(0xE5484D)
	case colorImage:
		return hex(0x12A594)
	case colorDoc:
		return hex(0x3E63DD)
	case colorBatch:
		return hex(0x8E4EC6)
	case colorOnAccent:
		return hex(0xFFFFFF)

	// Neutrals.
	case colorMuted:
		return pick(0x6B7280, 0x9AA1AC)
	case colorCard:
		return pick(0xFFFFFF, 0x1E2127)
	case colorCardHi:
		return pick(0xEDEFF3, 0x272C34)

	// Standard names, retuned for a cleaner look.
	case theme.ColorNamePrimary, theme.ColorNameHyperlink:
		return hex(accentPrimary)
	case theme.ColorNameFocus:
		return alpha(hex(accentPrimary), 0x80)
	case theme.ColorNameBackground:
		return pick(0xF5F6F8, 0x15171C)
	case theme.ColorNameForeground:
		return pick(0x1A1D21, 0xE7E9EC)
	case theme.ColorNameInputBackground:
		return pick(0xFFFFFF, 0x1E2127)
	case theme.ColorNameInputBorder, theme.ColorNameSeparator:
		return pick(0xE1E4E9, 0x2C313A)
	case theme.ColorNameButton:
		return pick(0xEDEFF3, 0x272C34)
	case theme.ColorNameHover:
		return pick(0xE6E9EE, 0x2E343D)
	case theme.ColorNamePlaceHolder:
		return hex(0x9AA1AC)
	}
	return t.Theme.Color(name, v)
}

func (t *appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius:
		return 9
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 26
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameInputBorder:
		return 1
	}
	return t.Theme.Size(name)
}

// themeColor resolves a theme color name for the current variant.
func themeColor(name fyne.ThemeColorName) color.Color {
	return fyne.CurrentApp().Settings().Theme().Color(name, fyne.CurrentApp().Settings().ThemeVariant())
}
