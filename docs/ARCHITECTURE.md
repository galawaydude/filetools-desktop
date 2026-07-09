# Architecture

The guiding idea: **adding a feature should touch one small file and nothing
else.** Everything hangs off a single `Tool` abstraction and a registry the UI
renders from.

## The core abstraction

`internal/tool` defines:

- **`Tool`** — one file operation (id, name, category, accepted extensions,
  options, and a `Run` method). This is the Strategy/Command pattern: the runner
  and UI treat every operation identically.
- **`Registry`** — a thread-safe set of tools, grouped by category. `tool.Default`
  is the process-wide registry.
- **`Config` + `Define`** — a declarative helper so a tool is written as one data
  block instead of an eight-method type.

Engine packages register themselves in `init()`:

```go
func init() {
    tool.Register(tool.Define(tool.Config{
        ID: "pdf.merge", Name: "Merge PDFs", Category: tool.CategoryPDF,
        Input: tool.InputMultiFile, Extensions: []string{".pdf"},
        Run: func(ctx, req, p) (tool.Result, error) { /* ... */ },
    }))
}
```

`cmd/filetools/main.go` blank-imports the engine packages so their `init()`
functions run (the `database/sql` driver pattern). The UI iterates the registry —
it never names a specific tool, so **new tools appear automatically**.

## Layers

```
UI (Fyne)                internal/ui
  │ builds screens from the registry; marshals background updates with fyne.Do
  ▼
validation               internal/validate   (plain-language pre-flight checks)
  ▼
Tool.Run                 internal/engine/{pdf,image,doc}
  │ Facade over pdfcpu / imaging; Adapter over LibreOffice
  ▼
platform + libraries     internal/platform + third-party libs
```

- **Facade** — `engine/pdf` and `engine/image` hide the details of pdfcpu and
  imaging behind small functions.
- **Adapter** — `engine/doc` wraps LibreOffice so its *absence* is a capability
  flag and a friendly message, never a crash.
- **Observer + Context** — `Run` receives a `Progress` callback and a
  `context.Context`; the UI shows progress and cancels through them, keeping the
  window responsive because work runs on a goroutine.
- **Never-overwrite pipeline** — `platform.UniqueName` guarantees new files;
  multi-file tools write into a same-filesystem temp dir and move results in.

## Why these choices

- **Fyne** builds to one self-contained, offline executable and cross-packages to
  Windows — the shareable-installer requirement.
- **pdfcpu / imaging / nativewebp / fpdf** are pure Go, so the core tools have no
  runtime dependencies and the build stays simple.
- The only optional external program is **LibreOffice**, and only for Word ↔ PDF,
  which genuinely cannot be done reliably in pure Go.

## Adding a new tool

1. Create a file in the relevant `internal/engine/*` package.
2. `tool.Register(tool.Define(tool.Config{ ... }))` in its `init()`.
3. If it's a brand-new engine package, add one blank import to `main.go`.

That's it — the home card count, the tool list, the file picker, the options
form, progress, validation and result dialog all come for free.
