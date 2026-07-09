// Package doc implements document tools. TXT->PDF is pure Go; Word<->PDF
// (added later) rely on an optional external converter.
package doc

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "doc.txttopdf",
		Name:        "Text file to PDF",
		Description: "Turn a plain .txt file into a simple, readable PDF.",
		Category:    tool.CategoryDoc,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".txt"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			data, err := os.ReadFile(in)
			if err != nil {
				return tool.Result{}, fmt.Errorf("could not read this text file: %w", err)
			}
			if err := cancelledCtx(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.3, "Creating PDF…")

			pdf := fpdf.New("P", "mm", "A4", "")
			pdf.SetMargins(15, 15, 15)
			pdf.SetAutoPageBreak(true, 15)
			pdf.AddPage()
			pdf.SetFont("Courier", "", 10)
			tr := pdf.UnicodeTranslatorFromDescriptor("") // map text into the core font's charset

			text := strings.ReplaceAll(string(data), "\t", "    ")
			pdf.MultiCell(0, 5, tr(text), "", "", false)
			if pdf.Err() {
				return tool.Result{}, fmt.Errorf("could not build the PDF: %w", pdf.Error())
			}

			out := platform.UniqueName(req.OutDir, platform.OutputName(in, "", "pdf"))
			if err := pdf.OutputFileAndClose(out); err != nil {
				return tool.Result{}, fmt.Errorf("could not save the PDF: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}

func cancelledCtx(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
