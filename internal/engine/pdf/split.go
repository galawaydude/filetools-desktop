package pdf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
			outputs, err := moveAllPDFs(tmp, req.OutDir, p)
			if err != nil {
				return tool.Result{}, err
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: outputs, Message: fmt.Sprintf("Created %d files.", len(outputs))}, nil
		},
	}))
}

func moveAllPDFs(fromDir, toDir string, p tool.Progress) ([]string, error) {
	entries, err := os.ReadDir(fromDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".pdf" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	var outputs []string
	for i, name := range names {
		dst, err := moveInto(filepath.Join(fromDir, name), toDir)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, dst)
		p.Update(0.2+0.7*float64(i+1)/float64(len(names)), "Saving files…")
	}
	return outputs, nil
}
