# File Tools

A simple, offline desktop app for everyday file conversion and manipulation —
PDF, image and document tools — with a clean, beginner-friendly interface.

> **Status:** in active development. See [milestones](#milestones).

## What it does

Open the app, pick a tool, choose your files, click one button, get new files.
No terminal, no setup, nothing to configure.

- **Never overwrites your originals.** Every result is saved as a new file in a
  folder you choose.
- **Works offline.** Core PDF and image tools need no internet and no extra
  software.
- **Clear progress and messages.** Long jobs show progress and can be cancelled;
  results come with an "Open output folder" button.

## Feature overview

| Area | Tools | Notes |
|------|-------|-------|
| PDF | merge, split, compress, rotate, delete pages, reorder, resize, images→PDF, extract images, extract text, TXT→PDF | Pure Go, always available |
| Image | convert (JPG/PNG/BMP/TIFF/GIF/WEBP), resize, compress, crop, batch convert | Pure Go, always available |
| Document | Word→PDF, PDF→Word | Optional — uses LibreOffice if installed (see [Limitations](#limitations)) |

## Install (Windows)

Download the latest **`FileToolsSetup.exe`** from the
[Releases page](https://github.com/galawaydude/filetools-desktop/releases) and
run it. That's it — no dependencies to install for the core tools.

## Build from source

Requires [Go](https://go.dev/dl/) 1.24+.

```sh
git clone https://github.com/galawaydude/filetools-desktop
cd filetools-desktop
go run ./cmd/filetools
```

Full build/packaging instructions are added in a later milestone.

## Limitations

- **Word ↔ PDF** relies on [LibreOffice](https://www.libreoffice.org/) being
  installed. If it isn't, those two tools show a plain-language message instead
  of failing; every other tool still works.
- **PDF → images** extracts images embedded in the PDF. Rendering each *page* to
  an image needs an external renderer and is not included.

## Milestones

1. Scaffold + toolchain + repo ✅
2. PDF core tools
3. PDF ↔ assets (images/text)
4. Image tools
5. Documents (optional LibreOffice)
6. UI polish + UX
7. Packaging + CI installer
8. Docs + tests + release

## License

[MIT](LICENSE)
