package ui

import (
	"image/png"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	_ "github.com/galawaydude/filetools-desktop/internal/engine/doc"
	_ "github.com/galawaydude/filetools-desktop/internal/engine/image"
	_ "github.com/galawaydude/filetools-desktop/internal/engine/pdf"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// TestPreviewCapture renders the main screens to PNGs so the design can be
// eyeballed without a display. Set FILETOOLS_PREVIEW_DIR to capture.
func TestPreviewCapture(t *testing.T) {
	dir := os.Getenv("FILETOOLS_PREVIEW_DIR")
	if dir == "" {
		t.Skip("set FILETOOLS_PREVIEW_DIR to capture UI previews")
	}
	app := test.NewApp()
	app.Settings().SetTheme(newAppTheme())
	w := test.NewWindow(nil)
	w.Resize(fyne.NewSize(940, 660))
	u := &UI{fyne: app, win: w, registry: tool.Default}

	capture := func(name string) {
		w.Content().Refresh()
		f, err := os.Create(dir + "/" + name + ".png")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := png.Encode(f, w.Canvas().Capture()); err != nil {
			t.Fatal(err)
		}
	}

	u.showHome()
	capture("home")

	u.showCategory(tool.CategoryPDF)
	capture("category")

	if tl, ok := tool.Default.Get("pdf.split"); ok {
		u.showTool(tl)
		capture("tool")
	}
}
