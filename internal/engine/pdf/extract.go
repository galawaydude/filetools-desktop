package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	extractpdf "github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.extractimages",
		Name:        "Extract images from PDF",
		Description: "Save the images embedded inside a PDF as separate image files. (This does not turn each page into a picture.)",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			if _, err := pageCount(in); err != nil {
				return tool.Result{}, err
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.2, "Looking for images…")
			tmp, cleanup, err := tempDirIn(req.OutDir)
			if err != nil {
				return tool.Result{}, err
			}
			defer cleanup()
			if err := api.ExtractImagesFile(in, tmp, nil, nil); err != nil {
				return tool.Result{}, fmt.Errorf("could not read images from this PDF: %w", err)
			}
			outputs, err := moveAll(tmp, req.OutDir, nil, p, 0.2, 0.9)
			if err != nil {
				return tool.Result{}, err
			}
			p.Update(1, "Done")
			if len(outputs) == 0 {
				return tool.Result{Message: "No embedded images were found in this PDF."}, nil
			}
			return tool.Result{Outputs: outputs, Message: fmt.Sprintf("Extracted %d image(s).", len(outputs))}, nil
		},
	}))

	tool.Register(tool.Define(tool.Config{
		ID:          "pdf.extracttext",
		Name:        "Extract text from PDF",
		Description: "Pull the selectable text out of a PDF into a plain text file. Works for text PDFs, not scanned pages.",
		Category:    tool.CategoryPDF,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.3, "Reading text…")
			text, err := extractText(in)
			if err != nil {
				return tool.Result{}, err
			}
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, "", "txt"))
			if err := os.WriteFile(out, []byte(text), 0o644); err != nil {
				return tool.Result{}, fmt.Errorf("could not save the text file: %w", err)
			}
			p.Update(1, "Done")
			if strings.TrimSpace(text) == "" {
				return tool.Result{Outputs: []string{out}, Message: "No selectable text was found — this PDF is likely made of scanned images."}, nil
			}
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}

// extractText returns the plain text of a PDF. It guards against panics in the
// third-party reader on malformed files.
func extractText(path string) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("this PDF could not be read for text")
		}
	}()
	f, r, err := extractpdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("this file could not be read as a PDF: %w", err)
	}
	defer f.Close()
	rc, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("could not extract text from this PDF: %w", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, rc); err != nil {
		return "", fmt.Errorf("could not extract text from this PDF: %w", err)
	}
	return buf.String(), nil
}
