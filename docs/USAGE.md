# Usage guide

Every tool follows the same five steps:

1. **Open the app** and click a category card (PDF, Image, Document or Batch).
2. **Pick a tool.**
3. **Add your file(s)** with **Add a file**. For tools that take several files you
   can add more one at a time and reorder them with the ↑ / ↓ buttons (order
   matters for Merge and Images → PDF).
4. **Choose options** if the tool has any, and **choose the folder** to save
   results in (it defaults to the folder your first file came from).
5. Click **Process**. Watch the progress bar (you can **Cancel**), then use
   **Open output folder** to see your new files.

Your original files are never changed.

## PDF tools

| Tool | What it does | Options |
|------|--------------|---------|
| Merge PDFs | Joins several PDFs into one, in the order listed | — |
| Split PDF | Breaks one PDF into smaller PDFs | Pages per file |
| Compress PDF | Makes a PDF smaller | — |
| Rotate PDF | Turns every page clockwise | 90 / 180 / 270 |
| Delete pages | Removes chosen pages | Pages, e.g. `1,3,5-7` |
| Reorder pages | Rebuilds a PDF in a new page order | New order, e.g. `3,1,2` |
| Resize PDF pages | Changes page size or scales down | A4 / A3 / A5 / Letter / Legal / 50% / 75% |
| Images to PDF | One PDF, one image per page | — |
| Extract images from PDF | Saves images embedded in the PDF | — |
| Extract text from PDF | Saves selectable text to a `.txt` | — |

**Page numbers** use commas for individual pages and hyphens for ranges, e.g.
`1,3,5-7`.

## Image tools

| Tool | What it does | Options |
|------|--------------|---------|
| Convert image format | Changes images to another format | Target format |
| Resize images | Changes size (set only one side to keep the shape) | Width, Height |
| Compress images | Shrinks file size (best for JPG) | Quality 1–100 |
| Crop image | Cuts out a rectangle | Left, Top, Width, Height |

All image tools accept multiple files at once.

## Document tools

| Tool | What it does | Needs |
|------|--------------|-------|
| Text file to PDF | Turns a `.txt` into a simple PDF | — |
| Word to PDF | Converts `.docx` / `.doc` to PDF | LibreOffice |
| PDF to Word | Converts a PDF to `.docx` (best-effort) | LibreOffice |

If LibreOffice isn't installed, the two Word tools tell you so and do nothing to
your files. Install it free from [libreoffice.org](https://www.libreoffice.org/),
then reopen File Tools.

## Batch Convert

**Batch convert images** turns a whole set of images into one chosen format in a
single run, and can shrink images wider than a size you set (leave `0` to keep the
original size).

## Tips

- If a file is open in another program (e.g. a PDF open in a viewer), close it
  first — locked files can't be read.
- Very large files still work but take longer; you'll be asked to confirm.
- Nothing you do here uploads your files anywhere. It all runs on your computer.
