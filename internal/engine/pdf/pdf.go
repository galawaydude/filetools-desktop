// Package pdf implements the PDF tools on top of pdfcpu (pure Go). Each tool
// self-registers into tool.Default via init().
package pdf

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/galawaydude/filetools-desktop/internal/platform"
)

func init() {
	// Keep pdfcpu from writing a config directory on the user's machine.
	api.DisableConfigDir()
}

// pageCount returns the number of pages in a PDF, wrapping pdfcpu errors as a
// plain-language "not a readable PDF" message.
func pageCount(path string) (int, error) {
	n, err := api.PageCountFile(path)
	if err != nil {
		return 0, fmt.Errorf("this file could not be read as a PDF: %w", err)
	}
	return n, nil
}

// tempDirIn creates a scratch directory on the same filesystem as outDir so
// finished files can be moved with a cheap rename.
func tempDirIn(outDir string) (string, func(), error) {
	d, err := os.MkdirTemp(outDir, ".filetools-*")
	if err != nil {
		return "", nil, fmt.Errorf("could not create a working folder in the output location: %w", err)
	}
	return d, func() { os.RemoveAll(d) }, nil
}

// moveInto moves src into dir under a non-colliding name and returns the path.
func moveInto(src, dir string) (string, error) {
	dst := platform.UniqueName(dir, filepath.Base(src))
	if err := os.Rename(src, dst); err == nil {
		return dst, nil
	}
	// Cross-device fallback: copy then remove.
	if err := copyFile(src, dst); err != nil {
		return "", err
	}
	os.Remove(src)
	return dst, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// cancelled reports whether ctx is done, mapping to its error.
func cancelled(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
