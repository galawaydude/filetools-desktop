package doc_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"

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
