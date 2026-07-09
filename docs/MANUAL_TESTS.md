# Manual test checklist

Run through these after installing a new build. Use throwaway copies of real
files. ✅ = works, and importantly, **the original file is unchanged** every time.

## Setup
- [ ] Installer runs and finishes; Start Menu and Desktop shortcuts appear.
- [ ] App opens to the home screen with four cards.

## PDF
- [ ] **Merge**: add 2–3 PDFs, reorder with ↑/↓, Process → one combined PDF in the
      right order.
- [ ] **Split**: split a multi-page PDF, "pages per file" = 2 → several PDFs.
- [ ] **Compress**: produces a new PDF; success message shows size change.
- [ ] **Rotate**: 90° → pages are rotated in the output.
- [ ] **Delete pages**: `1,3` → those pages are gone.
- [ ] **Reorder**: `3,2,1` → pages reversed.
- [ ] **Resize**: A4 → new PDF saved.
- [ ] **Images to PDF**: add 2 images → 2-page PDF.
- [ ] **Extract images**: on a PDF with images → images saved; on one without →
      friendly "no images found" message.
- [ ] **Extract text**: text PDF → `.txt` with text; scanned PDF → message that no
      text was found.

## Image
- [ ] **Convert** a PNG to JPG, WEBP, BMP, TIFF, GIF → each opens correctly.
- [ ] **Resize** width 100, height 0 → width becomes 100, aspect kept.
- [ ] **Compress** a JPG at quality 30 → smaller file.
- [ ] **Crop** with sensible box → cropped image of the right size.
- [ ] Convert several images at once → all converted.

## Document
- [ ] **Text to PDF**: a `.txt` → readable PDF.
- [ ] **Word to PDF** *with* LibreOffice installed → PDF produced.
- [ ] **Word to PDF** *without* LibreOffice → clear "install LibreOffice" message,
      no crash, no files changed.

## Batch
- [ ] **Batch convert images** to JPG, max width 200 → all converted and shrunk.

## Cross-cutting behaviour
- [ ] **No overwrite**: run the same tool twice into the same folder → second
      output is named `… (1)`.
- [ ] **Open output folder** button opens the correct folder.
- [ ] **Cancel** a long job (e.g. many images) → "cancelled, no files changed".
- [ ] **Bad file**: pick a non-matching/corrupt file → plain-language error.
- [ ] **Locked file**: keep a PDF open in a viewer, try to use it → readable error
      (behaviour may vary by program).
- [ ] **Read-only output folder** → clear "cannot write to this folder" message.
- [ ] **Large file** (>300 MB) → asked to confirm before running.
