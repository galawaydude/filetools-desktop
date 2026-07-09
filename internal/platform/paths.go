package platform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UniqueName returns a path in dir for filename, appending " (1)", " (2)"…
// before the extension until it does not collide with an existing file.
func UniqueName(dir, filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	candidate := filepath.Join(dir, filename)
	for i := 1; ; i++ {
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate
		}
		candidate = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", base, i, ext))
	}
}

// Stem returns a file's base name without its extension.
func Stem(path string) string {
	b := filepath.Base(path)
	return strings.TrimSuffix(b, filepath.Ext(b))
}

// OutputName builds "<stem><suffix>.<ext>" from an input path.
// ext may be given with or without a leading dot.
func OutputName(inputPath, suffix, ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	return fmt.Sprintf("%s%s.%s", Stem(inputPath), suffix, ext)
}
