package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "image.convert",
		Name:        "Convert image format",
		Description: "Change one or more images to another format (JPG, PNG, WEBP, BMP, TIFF or GIF).",
		Category:    tool.CategoryImage,
		Input:       tool.InputMultiFile,
		Extensions:  inputExtensions,
		Options: []tool.Option{{
			Key: "format", Label: "Convert to", Type: tool.OptChoice,
			Default: "PNG", Choices: formatChoices, Help: "The format your images will be saved as.",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if len(req.Inputs) == 0 {
				return tool.Result{}, errors.New("please choose at least one image")
			}
			format := req.Options.StringOr("format", "PNG")
			outputs, err := convertImages(ctx, req.Inputs, req.OutDir, format, 0, 90, p)
			if err != nil {
				return tool.Result{}, err
			}
			return tool.Result{Outputs: outputs, Message: webpNote(format)}, nil
		},
	}))

	tool.Register(tool.Define(tool.Config{
		ID:          "image.batch",
		Name:        "Batch convert images",
		Description: "Convert a whole batch of images to one format at once, and optionally shrink very large ones.",
		Category:    tool.CategoryBatch,
		Input:       tool.InputMultiFile,
		Extensions:  inputExtensions,
		Options: []tool.Option{
			{Key: "format", Label: "Convert to", Type: tool.OptChoice, Default: "JPG", Choices: formatChoices},
			{Key: "maxwidth", Label: "Shrink images wider than (pixels)", Type: tool.OptInt, Default: "0", Min: 0,
				Help: "Leave 0 to keep the original size."},
		},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if len(req.Inputs) == 0 {
				return tool.Result{}, errors.New("please choose at least one image")
			}
			format := req.Options.StringOr("format", "JPG")
			maxWidth := req.Options.Int("maxwidth", 0)
			outputs, err := convertImages(ctx, req.Inputs, req.OutDir, format, maxWidth, 85, p)
			if err != nil {
				return tool.Result{}, err
			}
			return tool.Result{Outputs: outputs, Message: fmt.Sprintf("Converted %d image(s). %s", len(outputs), webpNote(format))}, nil
		},
	}))
}

// webpNote explains WEBP's lossless-only encoding when relevant.
func webpNote(format string) string {
	if format == "WEBP" {
		return "WEBP files are saved lossless, so they can be larger than lossy WEBP from other tools."
	}
	return ""
}
