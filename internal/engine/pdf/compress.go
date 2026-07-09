package pdf

import (
	"context"
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.compress",
		Name:        "Compress PDF",
		Description: "Make a PDF smaller by removing unused data and optimising its structure.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			before := fileSize(in)
			p.Update(0.2, "Compressing PDF…")
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (compressed)", "pdf"))
			if err := api.OptimizeFile(in, out, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not compress this PDF: %w", err)
			}
			p.Update(1, "Done")
			after := fileSize(out)
			msg := "Compression finished."
			if before > 0 && after > 0 && after < before {
				saved := 100 * (1 - float64(after)/float64(before))
				msg = fmt.Sprintf("Reduced from %s to %s (%.0f%% smaller).", humanSize(before), humanSize(after), saved)
			} else if after >= before {
				msg = "This PDF was already well optimised, so the size barely changed."
			}
			return tool.Result{Outputs: []string{out}, Message: msg}, nil
		},
	}))
}

func fileSize(path string) int64 {
	if fi, err := os.Stat(path); err == nil {
		return fi.Size()
	}
	return 0
}

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGT"[exp])
}
