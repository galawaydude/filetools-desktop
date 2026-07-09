package doc_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/engine/doc"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func TestTxtToPDF(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(src, []byte("Hello world\nLine two\n\tindented"), 0o644); err != nil {
		t.Fatal(err)
	}
	tl, ok := tool.Default.Get("doc.txttopdf")
	if !ok {
		t.Fatal("doc.txttopdf not registered")
	}
	out := t.TempDir()
	res, err := tl.Run(context.Background(), tool.Request{Inputs: []string{src}, OutDir: out}, tool.NopProgress)
	if err != nil {
		t.Fatalf("txt->pdf: %v", err)
	}
	if n, err := api.PageCountFile(res.Outputs[0]); err != nil || n < 1 {
		t.Fatalf("output PDF invalid: pages=%d err=%v", n, err)
	}
}

// When LibreOffice is not installed, Word conversions must fail gracefully with
// a plain-language hint rather than crashing.
func TestWordToPDFGracefulWithoutLibreOffice(t *testing.T) {
	if doc.LibreOfficeAvailable() {
		t.Skip("LibreOffice is installed; graceful-fallback path not exercised here")
	}
	dir := t.TempDir()
	src := filepath.Join(dir, "letter.docx")
	if err := os.WriteFile(src, []byte("not a real docx, just a placeholder"), 0o644); err != nil {
		t.Fatal(err)
	}
	tl, _ := tool.Default.Get("doc.wordtopdf")
	_, err := tl.Run(context.Background(), tool.Request{Inputs: []string{src}, OutDir: t.TempDir()}, tool.NopProgress)
	if err == nil {
		t.Fatal("expected an error when LibreOffice is missing")
	}
	if !strings.Contains(err.Error(), "LibreOffice") {
		t.Fatalf("error should mention LibreOffice, got: %v", err)
	}
}
