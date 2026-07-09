package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"

	"github.com/galawaydude/filetools-desktop/internal/tool"
	// Populate the registry with the real tools.
	_ "github.com/galawaydude/filetools-desktop/internal/engine/doc"
	_ "github.com/galawaydude/filetools-desktop/internal/engine/image"
	_ "github.com/galawaydude/filetools-desktop/internal/engine/pdf"
)

// TestBuildEveryToolView makes sure each registered tool's screen constructs
// without panicking (missing icons, nil option controls, bad layouts, …).
func TestBuildEveryToolView(t *testing.T) {
	test.NewApp()
	w := test.NewWindow(nil)
	defer w.Close()

	all := tool.Default.All()
	if len(all) == 0 {
		t.Fatal("no tools registered")
	}
	for _, tl := range all {
		v := &toolView{t: tl, win: w, multi: tl.InputKind() == tool.InputMultiFile}
		v.buildFilesSection()
		v.buildOptionsSection()
		v.buildOutputSection()
		if v.runButton() == nil {
			t.Fatalf("%s: nil run button", tl.ID())
		}
	}
}

func TestOptionControlsForAllTypes(t *testing.T) {
	test.NewApp()
	w := test.NewWindow(nil)
	defer w.Close()
	v := &toolView{win: w}
	cases := []tool.Option{
		{Key: "a", Type: tool.OptText, Default: "hi"},
		{Key: "b", Type: tool.OptInt, Default: "5"},
		{Key: "c", Type: tool.OptFloat, Default: "1.5"},
		{Key: "d", Type: tool.OptBool, Default: "true"},
		{Key: "e", Type: tool.OptChoice, Default: "x", Choices: []string{"x", "y"}},
	}
	for _, o := range cases {
		_, get := v.optionControl(o)
		if got := get(); got == "" && o.Type != tool.OptText {
			t.Fatalf("option %q returned empty value", o.Key)
		}
	}
}
