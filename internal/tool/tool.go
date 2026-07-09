// Package tool defines the core abstractions every feature of the app plugs
// into: a Tool (Strategy/Command) and a Registry that the UI iterates over.
//
// Adding a new feature = implement Tool + Register() it. The UI, batch runner
// and everything else pick it up automatically. Nothing else needs to change.
package tool

import "context"

// Category groups tools into the four home-screen cards.
type Category string

const (
	CategoryPDF   Category = "PDF Tools"
	CategoryImage Category = "Image Tools"
	CategoryDoc   Category = "Document Tools"
	CategoryBatch Category = "Batch Convert"
)

// Categories is the display order used by the home screen.
var Categories = []Category{CategoryPDF, CategoryImage, CategoryDoc, CategoryBatch}

// InputKind tells the UI how many files the tool consumes.
type InputKind int

const (
	// InputSingleFile: exactly one file (e.g. rotate a PDF).
	InputSingleFile InputKind = iota
	// InputMultiFile: one or more files (e.g. merge PDFs, batch convert).
	InputMultiFile
)

// OptionType tells the UI which control to render for an option.
type OptionType int

const (
	OptText   OptionType = iota // free text
	OptInt                      // whole number
	OptFloat                    // decimal number
	OptBool                     // on/off checkbox
	OptChoice                   // pick-one dropdown
)

// Option describes a single, plain-language setting the user can adjust.
type Option struct {
	Key     string     // machine key used to read the value back
	Label   string     // plain-language label shown to the user
	Help    string     // one-line hint under the control
	Type    OptionType // control to render
	Default string     // default value (as string; parsed by helpers)
	Choices []string   // for OptChoice
	Min     float64    // for OptInt/OptFloat (inclusive); 0 = unset
	Max     float64    // for OptInt/OptFloat (inclusive); 0 = unset
}

// Request is the validated work order handed to Tool.Run.
type Request struct {
	Inputs  []string // absolute paths, already validated
	OutDir  string   // existing, writable output directory chosen by the user
	Options Options  // user-selected option values
}

// Result is what a successful run reports back to the UI.
type Result struct {
	Outputs []string // absolute paths of newly created files (never originals)
	Message string   // optional plain-language note (e.g. limitations)
}

// Progress is the Observer the UI implements to reflect long-running work.
// Fraction is 0..1; message is a short plain-language status line.
type Progress interface {
	Update(fraction float64, message string)
}

// ProgressFunc adapts a plain function to the Progress interface.
type ProgressFunc func(fraction float64, message string)

// Update implements Progress.
func (f ProgressFunc) Update(fraction float64, message string) {
	if f != nil {
		f(fraction, message)
	}
}

// NopProgress is a Progress that discards updates (handy for tests).
var NopProgress Progress = ProgressFunc(func(float64, string) {})

// Tool is one file operation. Implementations must be safe to call from a
// background goroutine and must honour ctx cancellation promptly.
//
// Contract:
//   - Never modify or overwrite any input file.
//   - Write every output into Request.OutDir as a NEW file.
//   - Return promptly with ctx.Err() when the context is cancelled.
type Tool interface {
	ID() string           // stable unique id, e.g. "pdf.merge"
	Name() string         // plain-language name, e.g. "Merge PDFs"
	Category() Category   // home card it belongs to
	Description() string  // one/two sentence plain-language description
	InputKind() InputKind // single vs multiple files
	Extensions() []string // accepted input extensions, lower-case with dot
	Options() []Option    // user-adjustable settings (may be empty)
	Run(ctx context.Context, req Request, p Progress) (Result, error)
}
