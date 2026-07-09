package pdf

import (
	"context"
	"fmt"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func parseSelection(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("please type which pages, for example: 1,3,5-7")
	}
	sel, err := api.ParsePageSelection(s)
	if err != nil {
		return nil, fmt.Errorf("could not understand the page numbers %q — use a format like 1,3,5-7", s)
	}
	return sel, nil
}

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.delete",
		Name:        "Delete pages",
		Description: "Remove selected pages from a PDF. The rest are saved as a new PDF.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Options: []tool.Option{{
			Key: "pages", Label: "Pages to delete", Type: tool.OptText,
			Default: "", Help: "For example: 1,3,5-7",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			sel, err := parseSelection(req.Options.String("pages"))
			if err != nil {
				return tool.Result{}, err
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Removing pages…")
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (pages removed)", "pdf"))
			if err := api.RemovePagesFile(in, out, sel, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not delete those pages: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))

	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.reorder",
		Name:        "Reorder pages",
		Description: "Rebuild a PDF with its pages in an order you choose.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Options: []tool.Option{{
			Key: "order", Label: "New page order", Type: tool.OptText,
			Default: "", Help: "List every page in the order you want, e.g. 3,1,2",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			sel, err := parseSelection(req.Options.String("order"))
			if err != nil {
				return tool.Result{}, err
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Reordering pages…")
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (reordered)", "pdf"))
			if err := api.CollectFile(in, out, sel, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not reorder the pages: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}
