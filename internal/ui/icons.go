package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// categoryColor returns the accent colour name for a category.
func categoryColor(c tool.Category) fyne.ThemeColorName {
	switch c {
	case tool.CategoryPDF:
		return colorPDF
	case tool.CategoryImage:
		return colorImage
	case tool.CategoryDoc:
		return colorDoc
	case tool.CategoryBatch:
		return colorBatch
	default:
		return colorImage
	}
}

// categoryIcon returns a representative icon for a category.
func categoryIcon(c tool.Category) fyne.Resource {
	switch c {
	case tool.CategoryPDF:
		return theme.DocumentIcon()
	case tool.CategoryImage:
		return theme.MediaPhotoIcon()
	case tool.CategoryDoc:
		return theme.FileTextIcon()
	case tool.CategoryBatch:
		return theme.ContentCopyIcon()
	default:
		return theme.FileIcon()
	}
}

// toolIcon returns a per-tool icon, falling back to the category icon.
func toolIcon(t tool.Tool) fyne.Resource {
	switch t.ID() {
	case "pdf.merge":
		return theme.ContentAddIcon()
	case "pdf.split":
		return theme.ContentCutIcon()
	case "pdf.compress":
		return theme.MoveDownIcon()
	case "pdf.rotate":
		return theme.ViewRefreshIcon()
	case "pdf.delete":
		return theme.DeleteIcon()
	case "pdf.reorder":
		return theme.ListIcon()
	case "pdf.resize":
		return theme.ViewFullScreenIcon()
	case "pdf.imagestopdf":
		return theme.MediaPhotoIcon()
	case "pdf.extractimages":
		return theme.MediaPhotoIcon()
	case "pdf.extracttext":
		return theme.FileTextIcon()
	case "image.convert":
		return theme.ViewRefreshIcon()
	case "image.resize":
		return theme.ViewFullScreenIcon()
	case "image.compress":
		return theme.MoveDownIcon()
	case "image.crop":
		return theme.ContentCutIcon()
	case "image.batch":
		return theme.ContentCopyIcon()
	case "doc.txttopdf":
		return theme.DocumentCreateIcon()
	case "doc.wordtopdf":
		return theme.DocumentIcon()
	case "doc.pdftoword":
		return theme.FileTextIcon()
	default:
		return categoryIcon(t.Category())
	}
}
