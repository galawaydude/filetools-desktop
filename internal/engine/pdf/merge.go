package pdf

import (
	"context"
	"errors"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.merge",
		Name:        "Merge PDFs",
		Description: "Combine several PDF files into one, in the order you add them.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputMultiFile,
		Extensions:  []string{".pdf"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if len(req.Inputs) < 2 {
				return tool.Result{}, errors.New("please choose at least two PDF files to merge")
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Merging PDFs…")
			out := platform.UniqueName(req.OutDir, "Merged.pdf")
			if err := api.MergeCreateFile(req.Inputs, out, false, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not merge these PDFs: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}
