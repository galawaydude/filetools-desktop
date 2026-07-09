// Package validate performs plain-language pre-flight checks so the engines can
// assume their inputs are sane and users get friendly messages up front.
package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// largeFileBytes is the point past which we mention a file is large (info only).
const largeFileBytes = 300 << 20 // 300 MB

// Inputs checks the selected files exist, are readable, and match the tool's
// accepted extensions. It enforces the single-vs-multiple file rule.
func Inputs(inputs []string, allowed []string, kind tool.InputKind) error {
	if len(inputs) == 0 {
		return fmt.Errorf("please choose at least one file")
	}
	if kind == tool.InputSingleFile && len(inputs) != 1 {
		return fmt.Errorf("this tool works on a single file — please choose just one")
	}
	for _, p := range inputs {
		if err := oneFile(p, allowed); err != nil {
			return err
		}
	}
	return nil
}

func oneFile(path string, allowed []string) error {
	name := filepath.Base(path)
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the file %q no longer exists", name)
		}
		return fmt.Errorf("the file %q could not be read: %v", name, err)
	}
	if fi.IsDir() {
		return fmt.Errorf("%q is a folder, not a file", name)
	}
	if !extAllowed(path, allowed) {
		return fmt.Errorf("%q is not a supported file type for this tool (allowed: %s)", name, strings.Join(allowed, ", "))
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied for %q — it may be open in another program", name)
		}
		return fmt.Errorf("could not open %q: %v", name, err)
	}
	f.Close()
	return nil
}

func extAllowed(path string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	e := strings.ToLower(filepath.Ext(path))
	for _, a := range allowed {
		if e == strings.ToLower(a) {
			return true
		}
	}
	return false
}

// OutputDir checks the chosen output folder exists and is writable.
func OutputDir(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return fmt.Errorf("please choose a folder to save the results in")
	}
	fi, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("the output folder could not be found — please choose another")
	}
	if !fi.IsDir() {
		return fmt.Errorf("the output location is not a folder — please choose a folder")
	}
	probe, err := os.CreateTemp(dir, ".filetools-write-*")
	if err != nil {
		return fmt.Errorf("this folder cannot be written to (permission issue) — please choose another")
	}
	probe.Close()
	os.Remove(probe.Name())
	return nil
}

// HasLargeFile reports whether any input is large enough to be worth warning
// about, along with that file's name.
func HasLargeFile(inputs []string) (string, bool) {
	for _, p := range inputs {
		if fi, err := os.Stat(p); err == nil && fi.Size() > largeFileBytes {
			return filepath.Base(p), true
		}
	}
	return "", false
}
