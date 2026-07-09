// Command filetools is a simple, offline desktop app for everyday file
// conversion and manipulation (PDF, image and document tools).
//
// The blank imports below pull in each engine package purely for the side
// effect of self-registering its tools into tool.Default (the database/sql
// driver pattern). Adding a feature never requires touching this file beyond
// adding one blank import for a brand-new engine package.
package main

import (
	"github.com/galawaydude/filetools-desktop/internal/tool"
	"github.com/galawaydude/filetools-desktop/internal/ui"
	// Engine packages register their tools via init():
	_ "github.com/galawaydude/filetools-desktop/internal/engine/doc"
	_ "github.com/galawaydude/filetools-desktop/internal/engine/image"
	_ "github.com/galawaydude/filetools-desktop/internal/engine/pdf"
)

func main() {
	ui.New(tool.Default).Run()
}
