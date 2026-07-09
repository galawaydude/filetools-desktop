package pdf_test

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-pdf/fpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// makePNG writes a small solid-colour PNG and returns its path.
func makePNG(t *testing.T, dir, name string, c color.Color) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 64, 48))
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, c)
		}
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("makePNG: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("makePNG encode: %v", err)
	}
	return path
}

// makePDF writes a PDF with the given number of pages and returns its path.
func makePDF(t *testing.T, dir, name string, pages int) string {
	t.Helper()
	doc := fpdf.New("P", "mm", "A4", "")
	for i := 0; i < pages; i++ {
		doc.AddPage()
		doc.SetFont("Arial", "", 16)
		doc.Cell(40, 10, fmt.Sprintf("Page %d", i+1))
	}
	path := filepath.Join(dir, name)
	if err := doc.OutputFileAndClose(path); err != nil {
		t.Fatalf("makePDF: %v", err)
	}
	return path
}

func run(t *testing.T, id string, req tool.Request) tool.Result {
	t.Helper()
	tl, ok := tool.Default.Get(id)
	if !ok {
		t.Fatalf("tool %q not registered", id)
	}
	res, err := tl.Run(context.Background(), req, tool.NopProgress)
	if err != nil {
		t.Fatalf("%s: %v", id, err)
	}
	return res
}

func pages(t *testing.T, path string) int {
	t.Helper()
	n, err := api.PageCountFile(path)
	if err != nil {
		t.Fatalf("PageCountFile(%s): %v", path, err)
	}
	return n
}

func TestMerge(t *testing.T) {
	dir := t.TempDir()
	a := makePDF(t, dir, "a.pdf", 2)
	b := makePDF(t, dir, "b.pdf", 3)
	out := t.TempDir()
	res := run(t, "pdf.merge", tool.Request{Inputs: []string{a, b}, OutDir: out})
	if len(res.Outputs) != 1 {
		t.Fatalf("want 1 output, got %d", len(res.Outputs))
	}
	if got := pages(t, res.Outputs[0]); got != 5 {
		t.Fatalf("merged pages = %d, want 5", got)
	}
}

func TestMergeNeverOverwrites(t *testing.T) {
	dir := t.TempDir()
	a := makePDF(t, dir, "a.pdf", 1)
	b := makePDF(t, dir, "b.pdf", 1)
	out := t.TempDir()
	first := run(t, "pdf.merge", tool.Request{Inputs: []string{a, b}, OutDir: out})
	second := run(t, "pdf.merge", tool.Request{Inputs: []string{a, b}, OutDir: out})
	if first.Outputs[0] == second.Outputs[0] {
		t.Fatalf("second run reused the same path %q (would overwrite)", first.Outputs[0])
	}
	if _, err := os.Stat(first.Outputs[0]); err != nil {
		t.Fatalf("first output missing: %v", err)
	}
}

func TestSplit(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 4)
	out := t.TempDir()
	res := run(t, "pdf.split", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"span": "2"}})
	if len(res.Outputs) != 2 {
		t.Fatalf("split into %d files, want 2", len(res.Outputs))
	}
	for _, o := range res.Outputs {
		if got := pages(t, o); got != 2 {
			t.Fatalf("part %s has %d pages, want 2", o, got)
		}
	}
}

func TestCompress(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 3)
	out := t.TempDir()
	res := run(t, "pdf.compress", tool.Request{Inputs: []string{src}, OutDir: out})
	if got := pages(t, res.Outputs[0]); got != 3 {
		t.Fatalf("compressed pages = %d, want 3", got)
	}
}

func TestRotate(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 2)
	out := t.TempDir()
	res := run(t, "pdf.rotate", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"angle": "90"}})
	if got := pages(t, res.Outputs[0]); got != 2 {
		t.Fatalf("rotated pages = %d, want 2", got)
	}
}

func TestDeletePages(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 4)
	out := t.TempDir()
	res := run(t, "pdf.delete", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"pages": "2,3"}})
	if got := pages(t, res.Outputs[0]); got != 2 {
		t.Fatalf("after delete pages = %d, want 2", got)
	}
}

func TestReorder(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 3)
	out := t.TempDir()
	res := run(t, "pdf.reorder", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"order": "3,2,1"}})
	if got := pages(t, res.Outputs[0]); got != 3 {
		t.Fatalf("reordered pages = %d, want 3", got)
	}
}

func TestResize(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 2)
	out := t.TempDir()
	res := run(t, "pdf.resize", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"target": "A5"}})
	if _, err := os.Stat(res.Outputs[0]); err != nil {
		t.Fatalf("resize output missing: %v", err)
	}
}

func TestImagesToPDF(t *testing.T) {
	dir := t.TempDir()
	a := makePNG(t, dir, "a.png", color.RGBA{200, 30, 30, 255})
	b := makePNG(t, dir, "b.png", color.RGBA{30, 30, 200, 255})
	out := t.TempDir()
	res := run(t, "pdf.imagestopdf", tool.Request{Inputs: []string{a, b}, OutDir: out})
	if got := pages(t, res.Outputs[0]); got != 2 {
		t.Fatalf("images->pdf pages = %d, want 2", got)
	}
}

func TestExtractImages(t *testing.T) {
	dir := t.TempDir()
	img := makePNG(t, dir, "a.png", color.RGBA{10, 200, 10, 255})
	built := t.TempDir()
	src := run(t, "pdf.imagestopdf", tool.Request{Inputs: []string{img}, OutDir: built}).Outputs[0]

	out := t.TempDir()
	res := run(t, "pdf.extractimages", tool.Request{Inputs: []string{src}, OutDir: out})
	if len(res.Outputs) == 0 {
		t.Fatal("expected at least one extracted image")
	}
}

func TestExtractText(t *testing.T) {
	dir := t.TempDir()
	src := makePDF(t, dir, "src.pdf", 1)
	out := t.TempDir()
	res := run(t, "pdf.extracttext", tool.Request{Inputs: []string{src}, OutDir: out})
	if len(res.Outputs) != 1 {
		t.Fatalf("want 1 text file, got %d", len(res.Outputs))
	}
	if _, err := os.Stat(res.Outputs[0]); err != nil {
		t.Fatalf("text output missing: %v", err)
	}
}

func TestBadFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "not.pdf")
	os.WriteFile(bad, []byte("this is not a pdf"), 0o644)
	tl, _ := tool.Default.Get("pdf.split")
	_, err := tl.Run(context.Background(), tool.Request{Inputs: []string{bad}, OutDir: t.TempDir(), Options: tool.Options{"span": "1"}}, tool.NopProgress)
	if err == nil {
		t.Fatal("expected an error for a non-PDF file, got nil")
	}
}
