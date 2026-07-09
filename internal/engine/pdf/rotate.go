package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.rotate",
		Name:        "Rotate PDF",
		Description: "Turn every page of a PDF clockwise by 90, 180 or 270 degrees.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Options: []tool.Option{{
			Key: "angle", Label: "Rotate clockwise by", Type: tool.OptChoice,
			Default: "90", Choices: []string{"90", "180", "270"},
			Help: "Degrees to turn each page.",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			angle := req.Options.Int("angle", 90)
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Rotating pages…")
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (rotated)", "pdf"))
			if err := api.RotateFile(in, out, angle, nil, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not rotate this PDF: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}
