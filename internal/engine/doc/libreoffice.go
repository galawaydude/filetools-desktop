package doc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/galawaydude/filetools-desktop/internal/platform"
)

// installHint tells a non-technical user how to enable Word conversions.
const installHint = "This tool needs the free program LibreOffice, which was not found on your computer. " +
	"Install it from libreoffice.org, then reopen File Tools and try again."

// findLibreOffice locates the soffice executable on PATH or in common install
// locations. The second return value is false when it is not installed.
func findLibreOffice() (string, bool) {
	for _, name := range []string{"soffice", "libreoffice"} {
		if p, ok := platform.LookPath(name); ok {
			return p, true
		}
	}
	var candidates []string
	switch runtime.GOOS {
	case "windows":
		candidates = []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		}
	case "darwin":
		candidates = []string{"/Applications/LibreOffice.app/Contents/MacOS/soffice"}
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c, true
		}
	}
	return "", false
}

// libreConvert converts inPath to the given LibreOffice output filter (e.g.
// "pdf" or "docx"), writing a new file into outDir. It runs headless with a
// private profile so it never clashes with a LibreOffice window the user may
// already have open, and honours ctx cancellation.
func libreConvert(ctx context.Context, inPath, outDir, filter, outExt string) (string, error) {
	soffice, ok := findLibreOffice()
	if !ok {
		return "", fmt.Errorf("%s", installHint)
	}

	profile, err := os.MkdirTemp("", "filetools-lo-*")
	if err != nil {
		return "", fmt.Errorf("could not prepare the converter: %w", err)
	}
	defer os.RemoveAll(profile)

	work, err := os.MkdirTemp(outDir, ".filetools-*")
	if err != nil {
		return "", fmt.Errorf("could not create a working folder in the output location: %w", err)
	}
	defer os.RemoveAll(work)

	cmd := exec.CommandContext(ctx, soffice,
		"--headless", "--norestore", "--nolockcheck",
		"-env:UserInstallation="+fileURI(profile),
		"--convert-to", filter,
		"--outdir", work,
		inPath,
	)
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", fmt.Errorf("the converter could not process this file: %s", firstLine(string(out)))
	}

	produced := filepath.Join(work, platform.Stem(inPath)+"."+strings.TrimPrefix(outExt, "."))
	if _, err := os.Stat(produced); err != nil {
		return "", fmt.Errorf("the converter did not produce an output file for %q", filepath.Base(inPath))
	}
	final := platform.UniqueName(outDir, platform.OutputName(inPath, "", outExt))
	if err := os.Rename(produced, final); err != nil {
		return "", fmt.Errorf("could not save the converted file: %w", err)
	}
	return final, nil
}

func fileURI(path string) string {
	p := filepath.ToSlash(path)
	if runtime.GOOS == "windows" {
		return "file:///" + p
	}
	return "file://" + p
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if s == "" {
		return "unknown error"
	}
	return s
}

// LibreOfficeAvailable reports whether Word conversions can run right now.
func LibreOfficeAvailable() bool {
	_, ok := findLibreOffice()
	return ok
}
