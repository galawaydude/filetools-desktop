package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// vcenter vertically centres an object within a taller region.
func vcenter(o fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(layout.NewSpacer(), o, layout.NewSpacer())
}

// accentBadge is a rounded, filled square in an accent colour with a white icon.
func accentBadge(icon fyne.Resource, accent fyne.ThemeColorName, size float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(themeColor(accent))
	bg.CornerRadius = size * 0.28

	img := canvas.NewImageFromResource(theme.NewColoredResource(icon, colorOnAccent))
	img.FillMode = canvas.ImageFillContain

	inset := size * 0.24
	stack := container.NewStack(bg, container.New(&paddingLayout{inset}, img))
	return container.NewGridWrap(fyne.NewSize(size, size), stack)
}

// paddingLayout adds a uniform inset around a single object.
type paddingLayout struct{ pad float32 }

func (p *paddingLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	m := objs[0].MinSize()
	return fyne.NewSize(m.Width+2*p.pad, m.Height+2*p.pad)
}

func (p *paddingLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	objs[0].Move(fyne.NewPos(p.pad, p.pad))
	objs[0].Resize(fyne.NewSize(size.Width-2*p.pad, size.Height-2*p.pad))
}

// mutedText renders soft, wrapping secondary text.
func mutedText(s string) *widget.RichText {
	seg := &widget.TextSegment{Text: s, Style: widget.RichTextStyle{ColorName: colorMuted, SizeName: theme.SizeNameText}}
	rt := widget.NewRichText(seg)
	rt.Wrapping = fyne.TextWrapWord
	return rt
}

// boldText renders a bold title in the foreground colour.
func boldText(s string) *canvas.Text {
	t := canvas.NewText(s, themeColor(theme.ColorNameForeground))
	t.TextStyle = fyne.TextStyle{Bold: true}
	t.TextSize = 15
	return t
}

// hoverCard is a rounded, clickable card with an accent badge, a title, a
// subtitle and a chevron. It highlights on hover and shows a pointer cursor.
type hoverCard struct {
	widget.BaseWidget
	icon     fyne.Resource
	accent   fyne.ThemeColorName
	title    string
	subtitle string
	onTap    func()

	bg *canvas.Rectangle
}

func newHoverCard(icon fyne.Resource, accent fyne.ThemeColorName, title, subtitle string, onTap func()) *hoverCard {
	c := &hoverCard{icon: icon, accent: accent, title: title, subtitle: subtitle, onTap: onTap}
	c.ExtendBaseWidget(c)
	return c
}

func (c *hoverCard) Tapped(_ *fyne.PointEvent) {
	if c.onTap != nil {
		c.onTap()
	}
}

func (c *hoverCard) Cursor() desktop.Cursor { return desktop.PointerCursor }

func (c *hoverCard) MouseIn(_ *desktop.MouseEvent) {
	c.bg.FillColor = themeColor(colorCardHi)
	c.bg.StrokeColor = themeColor(c.accent)
	c.bg.StrokeWidth = 1.5
	canvas.Refresh(c.bg)
}

func (c *hoverCard) MouseMoved(_ *desktop.MouseEvent) {}

func (c *hoverCard) MouseOut() {
	c.bg.FillColor = themeColor(colorCard)
	c.bg.StrokeColor = themeColor(theme.ColorNameInputBorder)
	c.bg.StrokeWidth = 1
	canvas.Refresh(c.bg)
}

func (c *hoverCard) CreateRenderer() fyne.WidgetRenderer {
	c.bg = canvas.NewRectangle(themeColor(colorCard))
	c.bg.CornerRadius = 14
	c.bg.StrokeColor = themeColor(theme.ColorNameInputBorder)
	c.bg.StrokeWidth = 1

	badge := accentBadge(c.icon, c.accent, 46)
	text := container.NewVBox(boldText(c.title), mutedText(c.subtitle))
	chevron := widget.NewIcon(theme.NewColoredResource(theme.NavigateNextIcon(), colorMuted))

	row := container.NewBorder(nil, nil,
		vcenter(badge), vcenter(chevron),
		vcenter(text),
	)
	content := container.NewStack(c.bg, container.NewPadded(container.NewPadded(row)))
	return widget.NewSimpleRenderer(content)
}
