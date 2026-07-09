package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.split",
		Name:        "Split PDF",
		Description: "Break one PDF into smaller PDFs. Choose how many pages go into each part.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Options: []tool.Option{{
			Key: "span", Label: "Pages per file", Type: tool.OptInt,
			Default: "1", Min: 1, Help: "Each new PDF will contain this many pages.",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			span := req.Options.Int("span", 1)
			if span < 1 {
				span = 1
			}
			n, err := pageCount(in)
			if err != nil {
				return tool.Result{}, err
			}
			if span >= n {
				return tool.Result{}, fmt.Errorf("this PDF only has %d page(s); choose a smaller number of pages per file", n)
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Splitting PDF…")

			tmp, cleanup, err := tempDirIn(req.OutDir)
			if err != nil {
				return tool.Result{}, err
			}
			defer cleanup()

			if err := api.SplitFile(in, tmp, span, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not split this PDF: %w", err)
			}
			outputs, err := moveAll(tmp, req.OutDir, []string{".pdf"}, p, 0.2, 0.9)
			if err != nil {
				return tool.Result{}, err
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: outputs, Message: fmt.Sprintf("Created %d files.", len(outputs))}, nil
		},
	}))
}
