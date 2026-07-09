# File Tools

A simple, offline desktop app for everyday file conversion and manipulation —
PDF, image and document tools — with a clean, beginner-friendly interface.

Open the app, pick a tool, choose your files, click one button, get new files.
No terminal, no setup, nothing to configure.

- **Never overwrites your originals.** Every result is saved as a *new* file in a
  folder you choose. If a name is taken, it adds ` (1)`, ` (2)`, …
- **Works offline.** All core PDF and image tools need no internet and no extra
  software.
- **Clear progress and messages.** Long jobs show a progress bar and can be
  cancelled; results come with an **Open output folder** button.
- **Handles mistakes gracefully.** Bad files, wrong types, locked or unreadable
  files, missing output folders and very large files all produce a plain-language
  message instead of a crash.

## Install

Grab the latest build from the
[Releases page](https://github.com/galawaydude/filetools-desktop/releases/latest).

**Windows**
1. Download **`FileToolsSetup.exe`**.
2. Run it and follow the short wizard.
3. Launch **File Tools** from the Start Menu or Desktop.

**macOS** (universal — Intel and Apple Silicon)
1. Download **`FileTools.dmg`** and open it.
2. Drag **File Tools** into the **Applications** folder.
3. The app isn't signed with an Apple Developer account, so the first time you
   open it, **right-click the app → Open → Open** (a normal double-click will be
   blocked by Gatekeeper). You only need to do this once.
   - If macOS still refuses, run once in Terminal:
     `xattr -cr /Applications/FileTools.app`

No dependencies are required for the core tools. (Word ↔ PDF is optional — see
[Limitations](#limitations).)

> **Note:** the app renders with OpenGL, so it may not open inside a Remote
> Desktop session or a virtual machine that lacks GPU acceleration. Use a normal
> desktop session.

## What it does

### PDF tools
Merge, split, compress, rotate, delete pages, reorder pages, resize pages,
images → PDF, extract embedded images, extract text.

### Image tools
Convert between JPG, PNG, WEBP, BMP, TIFF and GIF; resize; compress; crop.
Every image tool can process many files at once.

### Document tools
Text file → PDF (built-in). Word → PDF and PDF → Word (need LibreOffice).

### Batch Convert
Convert a whole batch of images to one format, optionally shrinking large ones.

See [docs/USAGE.md](docs/USAGE.md) for a step-by-step guide to each tool.

## Limitations

Being honest about what is and isn't reliable offline in pure Go:

- **Word ↔ PDF** uses [LibreOffice](https://www.libreoffice.org/) if it is
  installed. If it isn't, those two tools show a friendly "please install
  LibreOffice" message — everything else still works. PDF → Word is best-effort;
  complex layouts may need tidying.
- **Extract images from PDF** saves the images *embedded* inside the PDF. It does
  not render each page to a picture (that needs an external renderer).
- **Extract text from PDF** works for PDFs that contain real text. Scanned/image
  PDFs have no selectable text, so nothing can be extracted.
- **WEBP** output is lossless, so WEBP files can be larger than lossy WEBP made by
  other tools. WEBP input (reading) is fully supported.

## Build from source

Requires [Go](https://go.dev/dl/) 1.24+ and a C compiler (needed by the GUI
toolkit). On macOS install the Xcode command line tools; on Windows install
MinGW-w64 (`choco install mingw`); on Linux install `gcc` and the usual X11/GL
dev headers.

```sh
git clone https://github.com/galawaydude/filetools-desktop
cd filetools-desktop
go run ./cmd/filetools     # run the app
go test ./...              # run the tests
```

### Build the installers locally

**Windows** (needs Go, MinGW and [NSIS](https://nsis.sourceforge.io/)):
```powershell
powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1
# -> build\FileToolsSetup.exe
```

**macOS** (needs Go and the Xcode command line tools):
```sh
./scripts/build-macos.sh
# -> FileTools.app and FileTools.dmg (universal)
```

### How releases are built

Pushing a version tag (e.g. `v0.3.0`) triggers
[.github/workflows/release.yml](.github/workflows/release.yml). A Windows runner
builds `FileToolsSetup.exe` (via `fyne package` + NSIS) and a macOS runner builds
a universal `FileTools.dmg`; both are attached to the GitHub Release. You can also
run the workflow manually from the Actions tab to get the installers as build
artifacts.

## How it's built

- **Language/UI:** Go with [Fyne](https://fyne.io) (single self-contained
  executable, works offline).
- **PDF:** [pdfcpu](https://github.com/pdfcpu/pdfcpu) (pure Go).
- **Images:** [imaging](https://github.com/disintegration/imaging) +
  [nativewebp](https://github.com/HugoSmits86/nativewebp) (pure Go).
- **Text/PDF:** [fpdf](https://github.com/go-pdf/fpdf) and
  [ledongthuc/pdf](https://github.com/ledongthuc/pdf).

The code is organised so new tools are easy to add: each tool implements one
`Tool` interface and registers itself, and the whole UI is generated from that
registry. See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Project layout

```
cmd/filetools        app entry point (wires the engines + UI)
internal/tool        Tool interface + registry (the core abstraction)
internal/engine/pdf  PDF tools
internal/engine/image  image tools
internal/engine/doc  document tools (TXT->PDF, Word<->PDF)
internal/validate    plain-language input/output checks
internal/platform    OS glue (unique names, open folder, find LibreOffice)
internal/ui          Fyne screens (home, category, tool flow)
build/               icon, NSIS installer, packaging metadata
```

## License

[MIT](LICENSE)
