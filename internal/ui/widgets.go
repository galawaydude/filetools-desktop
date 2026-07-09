package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// tappableCard is a large, obvious, clickable card used on the home screen and
// tool lists. Kept deliberately simple: an icon, a title, a subtitle.
type tappableCard struct {
	widget.BaseWidget
	icon     fyne.Resource
	title    string
	subtitle string
	onTap    func()
}

func newTappableCard(icon fyne.Resource, title, subtitle string, onTap func()) *tappableCard {
	c := &tappableCard{icon: icon, title: title, subtitle: subtitle, onTap: onTap}
	c.ExtendBaseWidget(c)
	return c
}

func (c *tappableCard) Tapped(_ *fyne.PointEvent) {
	if c.onTap != nil {
		c.onTap()
	}
}

func (c *tappableCard) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	bg.CornerRadius = 10
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	icon := widget.NewIcon(c.icon)
	title := widget.NewLabelWithStyle(c.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel(c.subtitle)
	subtitle.Wrapping = fyne.TextWrapWord

	content := container.NewBorder(
		nil, nil,
		container.NewPadded(icon), nil,
		container.NewVBox(title, subtitle),
	)
	obj := container.NewStack(bg, container.NewPadded(content))
	return widget.NewSimpleRenderer(obj)
}
