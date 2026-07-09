package pdf

import (
	"context"
	"errors"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpulib "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.imagestopdf",
		Name:        "Images to PDF",
		Description: "Turn one or more images into a single PDF, one image per page, in the order you add them.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputMultiFile,
		Extensions:  []string{".jpg", ".jpeg", ".png"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if len(req.Inputs) == 0 {
				return tool.Result{}, errors.New("please choose at least one image")
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Building PDF from images…")
			out := platform.UniqueName(req.OutDir, "Images.pdf")
			if err := api.ImportImagesFile(req.Inputs, out, pdfcpulib.DefaultImportConfig(), nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not create a PDF from these images: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}
