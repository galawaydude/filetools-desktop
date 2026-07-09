package platform

import (
	"os/exec"
	"runtime"
)

// OpenFolder opens the given directory in the system file manager.
func OpenFolder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// LookPath reports whether an executable is available on PATH.
func LookPath(name string) (string, bool) {
	p, err := exec.LookPath(name)
	return p, err == nil
}
