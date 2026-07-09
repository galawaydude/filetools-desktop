// Package image implements the image tools on top of the pure-Go imaging
// library, with WEBP encoding provided by nativewebp. Each tool self-registers.
package image

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp" // register WEBP decoder for imaging.Open

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

// formatChoices are the output formats the convert tool offers.
var formatChoices = []string{"JPG", "PNG", "WEBP", "BMP", "TIFF", "GIF"}

// inputExtensions are the image types the tools accept as input.
var inputExtensions = []string{".jpg", ".jpeg", ".png", ".webp", ".bmp", ".tif", ".tiff", ".gif"}

// extFor returns the file extension (with dot) for a format name.
func extFor(format string) string {
	switch strings.ToUpper(format) {
	case "JPG", "JPEG":
		return ".jpg"
	case "PNG":
		return ".png"
	case "WEBP":
		return ".webp"
	case "BMP":
		return ".bmp"
	case "TIFF":
		return ".tiff"
	case "GIF":
		return ".gif"
	default:
		return ".png"
	}
}

// formatFromExt maps a file's extension back to a format name.
func formatFromExt(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		return "JPG"
	case ".webp":
		return "WEBP"
	case ".bmp":
		return "BMP"
	case ".tif", ".tiff":
		return "TIFF"
	case ".gif":
		return "GIF"
	default:
		return "PNG"
	}
}

// decodeImage loads an image, honouring EXIF orientation.
func decodeImage(path string) (image.Image, error) {
	img, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("this image could not be opened (it may be corrupted or an unsupported type): %w", err)
	}
	return img, nil
}

// encodeImage writes img to outPath in the given format. quality applies to JPG.
func encodeImage(img image.Image, outPath, format string, quality int) error {
	if strings.ToUpper(format) == "WEBP" {
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		if err := nativewebp.Encode(f, img, &nativewebp.Options{}); err != nil {
			f.Close()
			return err
		}
		return f.Close()
	}
	var opts []imaging.EncodeOption
	if strings.ToUpper(format) == "JPG" {
		if quality <= 0 || quality > 100 {
			quality = 90
		}
		opts = append(opts, imaging.JPEGQuality(quality))
	}
	return imaging.Save(img, outPath, opts...)
}

func cancelled(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// convertImages decodes each input, optionally caps its width, and writes it in
// the target format to outDir. It reports progress and never overwrites.
func convertImages(ctx context.Context, inputs []string, outDir, format string, maxWidth, quality int, p tool.Progress) ([]string, error) {
	ext := extFor(format)
	var outputs []string
	for i, in := range inputs {
		if err := cancelled(ctx); err != nil {
			return outputs, err
		}
		img, err := decodeImage(in)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", filepath.Base(in), err)
		}
		if maxWidth > 0 && img.Bounds().Dx() > maxWidth {
			img = imaging.Resize(img, maxWidth, 0, imaging.Lanczos)
		}
		out := platform.UniqueName(outDir, platform.OutputName(in, "", ext))
		if err := encodeImage(img, out, format, quality); err != nil {
			return nil, fmt.Errorf("%s: could not save as %s: %w", filepath.Base(in), format, err)
		}
		outputs = append(outputs, out)
		p.Update(float64(i+1)/float64(len(inputs)), "Converting images…")
	}
	return outputs, nil
}
