package ui

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
	"github.com/galawaydude/filetools-desktop/internal/validate"
)

// toolView holds the state of a single tool's screen.
type toolView struct {
	ui   *UI
	t    tool.Tool
	win  fyne.Window
	multi bool

	inputs  []string
	outDir  string
	getOpts []func() (string, string) // returns (key, value)

	filesBox  *fyne.Container
	outLabel  *widget.Label
	bottom    *fyne.Container
	cancel    context.CancelFunc
}

// showTool renders the select -> options -> output -> run flow for a tool.
func (u *UI) showTool(t tool.Tool) {
	v := &toolView{ui: u, t: t, win: u.win, multi: t.InputKind() == tool.InputMultiFile}

	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() { u.showCategory(t.Category()) })
	back.Importance = widget.LowImportance
	title := widget.NewLabelWithStyle(t.Name(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.SizeName = theme.SizeNameHeadingText
	desc := widget.NewLabel(t.Description())
	desc.Wrapping = fyne.TextWrapWord
	header := container.NewVBox(container.NewHBox(back), title, desc, widget.NewSeparator())

	form := container.NewVBox(
		v.buildFilesSection(),
		v.buildOptionsSection(),
		v.buildOutputSection(),
	)

	v.bottom = container.NewStack(v.runButton())

	content := container.NewBorder(header, container.NewVBox(widget.NewSeparator(), v.bottom), nil, nil,
		container.NewVScroll(container.NewPadded(form)))
	u.setContent(container.NewPadded(content))
}

// ---- Files ----

func (v *toolView) buildFilesSection() fyne.CanvasObject {
	label := "Choose a file"
	if v.multi {
		label = "Choose files (add as many as you like)"
	}
	v.filesBox = container.NewVBox()
	v.refreshFiles()

	add := widget.NewButtonWithIcon("Add a file", theme.ContentAddIcon(), v.addFile)
	add.Importance = widget.HighImportance
	return sectionCard("1. "+label, container.NewVBox(v.filesBox, container.NewHBox(add)))
}

func (v *toolView) addFile() {
	d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil || r == nil {
			return
		}
		defer r.Close()
		p := r.URI().Path()
		if !v.multi {
			v.inputs = []string{p}
		} else {
			v.inputs = append(v.inputs, p)
		}
		if v.outDir == "" {
			v.setOutDir(filepath.Dir(p))
		}
		v.refreshFiles()
	}, v.win)
	if exts := v.t.Extensions(); len(exts) > 0 {
		d.SetFilter(storage.NewExtensionFileFilter(exts))
	}
	d.Show()
}

func (v *toolView) refreshFiles() {
	v.filesBox.Objects = v.filesBox.Objects[:0]
	if len(v.inputs) == 0 {
		hint := widget.NewLabel("No files chosen yet.")
		hint.TextStyle = fyne.TextStyle{Italic: true}
		v.filesBox.Add(hint)
	}
	for i := range v.inputs {
		v.filesBox.Add(v.fileRow(i))
	}
	v.filesBox.Refresh()
}

func (v *toolView) fileRow(i int) fyne.CanvasObject {
	name := widget.NewLabel(filepath.Base(v.inputs[i]))
	name.Truncation = fyne.TextTruncateEllipsis

	remove := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		v.inputs = append(v.inputs[:i], v.inputs[i+1:]...)
		v.refreshFiles()
	})
	remove.Importance = widget.LowImportance

	var left fyne.CanvasObject = widget.NewIcon(theme.FileIcon())
	if v.multi {
		up := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() { v.move(i, -1) })
		down := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() { v.move(i, 1) })
		up.Importance, down.Importance = widget.LowImportance, widget.LowImportance
		if i == 0 {
			up.Disable()
		}
		if i == len(v.inputs)-1 {
			down.Disable()
		}
		left = container.NewHBox(widget.NewLabel(fmt.Sprintf("%d.", i+1)), up, down)
	}
	return container.NewBorder(nil, nil, left, remove, name)
}

func (v *toolView) move(i, delta int) {
	j := i + delta
	if j < 0 || j >= len(v.inputs) {
		return
	}
	v.inputs[i], v.inputs[j] = v.inputs[j], v.inputs[i]
	v.refreshFiles()
}

// ---- Options ----

func (v *toolView) buildOptionsSection() fyne.CanvasObject {
	opts := v.t.Options()
	if len(opts) == 0 {
		return container.NewVBox()
	}
	rows := container.NewVBox()
	for _, o := range opts {
		o := o
		control, getter := v.optionControl(o)
		v.getOpts = append(v.getOpts, func() (string, string) { return o.Key, getter() })
		lbl := widget.NewLabelWithStyle(o.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		row := container.NewVBox(lbl, control)
		if o.Help != "" {
			help := widget.NewLabel(o.Help)
			help.TextStyle = fyne.TextStyle{Italic: true}
			row.Add(help)
		}
		rows.Add(row)
	}
	return sectionCard("2. Choose options", rows)
}

func (v *toolView) optionControl(o tool.Option) (fyne.CanvasObject, func() string) {
	switch o.Type {
	case tool.OptChoice:
		sel := widget.NewSelect(o.Choices, nil)
		sel.SetSelected(o.Default)
		return sel, func() string { return sel.Selected }
	case tool.OptBool:
		chk := widget.NewCheck("", nil)
		chk.SetChecked(tool.Options{o.Key: o.Default}.Bool(o.Key))
		return chk, func() string {
			if chk.Checked {
				return "true"
			}
			return "false"
		}
	default: // OptText, OptInt, OptFloat
		e := widget.NewEntry()
		e.SetText(o.Default)
		if o.Type == tool.OptInt || o.Type == tool.OptFloat {
			e.Validator = numericValidator(o.Type == tool.OptFloat)
		}
		return e, func() string { return e.Text }
	}
}

// ---- Output ----

func (v *toolView) buildOutputSection() fyne.CanvasObject {
	v.outLabel = widget.NewLabel("No folder chosen — results will be saved here.")
	v.outLabel.Wrapping = fyne.TextWrapWord
	choose := widget.NewButtonWithIcon("Choose folder", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			v.setOutDir(list.Path())
		}, v.win)
	})
	return sectionCard("3. Save results in", container.NewBorder(nil, nil, nil, choose, v.outLabel))
}

func (v *toolView) setOutDir(dir string) {
	v.outDir = dir
	if v.outLabel != nil {
		v.outLabel.SetText(dir)
	}
}

// ---- Run / progress / result ----

func (v *toolView) runButton() *widget.Button {
	b := widget.NewButtonWithIcon("Process", theme.MediaPlayIcon(), v.start)
	b.Importance = widget.HighImportance
	return b
}

func (v *toolView) start() {
	if err := validate.Inputs(v.inputs, v.t.Extensions(), v.t.InputKind()); err != nil {
		dialog.ShowError(err, v.win)
		return
	}
	if err := validate.OutputDir(v.outDir); err != nil {
		dialog.ShowError(err, v.win)
		return
	}

	opts := tool.Options{}
	for _, g := range v.getOpts {
		k, val := g()
		opts[k] = val
	}

	req := tool.Request{Inputs: append([]string(nil), v.inputs...), OutDir: v.outDir, Options: opts}

	if name, big := validate.HasLargeFile(v.inputs); big {
		dialog.ShowConfirm("Large file",
			fmt.Sprintf("%q is quite large, so this may take a while. Continue?", name),
			func(ok bool) {
				if ok {
					v.execute(req)
				}
			}, v.win)
		return
	}
	v.execute(req)
}

func (v *toolView) execute(req tool.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	bar := widget.NewProgressBar()
	status := widget.NewLabel("Starting…")
	cancelBtn := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() { cancel() })
	v.setBottom(container.NewVBox(status, bar, container.NewHBox(cancelBtn)))

	progress := tool.ProgressFunc(func(f float64, msg string) {
		fyne.Do(func() {
			bar.SetValue(f)
			if msg != "" {
				status.SetText(msg)
			}
		})
	})

	go func() {
		res, err := v.t.Run(ctx, req, progress)
		fyne.Do(func() { v.finish(res, err) })
	}()
}

func (v *toolView) finish(res tool.Result, err error) {
	v.setBottom(v.runButton())
	switch {
	case errors.Is(err, context.Canceled):
		dialog.ShowInformation("Cancelled", "The job was cancelled. No files were changed.", v.win)
	case err != nil:
		dialog.ShowError(err, v.win)
	default:
		v.showSuccess(res)
	}
}

func (v *toolView) showSuccess(res tool.Result) {
	var b strings.Builder
	switch n := len(res.Outputs); n {
	case 0:
		b.WriteString("Finished.")
	case 1:
		fmt.Fprintf(&b, "Done! Saved 1 file:\n%s", filepath.Base(res.Outputs[0]))
	default:
		fmt.Fprintf(&b, "Done! Saved %d files in the output folder.", n)
	}
	if res.Message != "" {
		b.WriteString("\n\n")
		b.WriteString(res.Message)
	}
	msg := widget.NewLabel(b.String())
	msg.Wrapping = fyne.TextWrapWord
	dialog.ShowCustomConfirm("Success", "Open output folder", "Close", msg, func(open bool) {
		if open {
			platform.OpenFolder(v.outDir)
		}
	}, v.win)
}

func (v *toolView) setBottom(o fyne.CanvasObject) {
	v.bottom.Objects = []fyne.CanvasObject{o}
	v.bottom.Refresh()
}

// ---- shared helpers ----

func sectionCard(title string, body fyne.CanvasObject) fyne.CanvasObject {
	head := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return container.NewVBox(head, body, widget.NewLabel(""))
}

func numericValidator(allowDecimal bool) fyne.StringValidator {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		for i, r := range s {
			switch {
			case r >= '0' && r <= '9':
			case r == '-' && i == 0:
			case r == '.' && allowDecimal:
			default:
				return errors.New("please enter a number")
			}
		}
		return nil
	}
}
