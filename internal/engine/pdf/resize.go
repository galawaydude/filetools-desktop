package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpulib "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// resizePresets maps a plain-language choice to a pdfcpu resize config string.
// pdfcpu's parser fills in the internal page dimensions, so we let it build the
// Resize value rather than constructing the struct by hand.
var resizePresets = map[string]string{
	"A4":          "formsize:A4",
	"A3":          "formsize:A3",
	"A5":          "formsize:A5",
	"Letter":      "formsize:Letter",
	"Legal":       "formsize:Legal",
	"50% smaller": "scalefactor:0.5",
	"75% smaller": "scalefactor:0.75",
}

func init() {
	choices := []string{"A4", "A3", "A5", "Letter", "Legal", "50% smaller", "75% smaller"}
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.resize",
		Name:        "Resize PDF pages",
		Description: "Change the page size of a PDF (for example to A4 or Letter) or scale every page down.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Options: []tool.Option{{
			Key: "target", Label: "Resize to", Type: tool.OptChoice,
			Default: "A4", Choices: choices, Help: "Pick a paper size or shrink every page.",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			target := req.Options.StringOr("target", "A4")
			cfg, ok := resizePresets[target]
			if !ok {
				return tool.Result{}, fmt.Errorf("unknown resize option %q", target)
			}
			resize, err := pdfcpulib.ParseResizeConfig(cfg, types.POINTS)
			if err != nil {
				return tool.Result{}, fmt.Errorf("could not prepare the resize: %w", err)
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Resizing pages…")
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (resized)", "pdf"))
			if err := api.ResizeFile(in, out, nil, resize, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not resize this PDF: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}
