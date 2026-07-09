package image

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "image.resize",
		Name:        "Resize images",
		Description: "Change the size of one or more images. Set only width or only height to keep the shape.",
		Category:    tool.CategoryImage,
		Input:       tool.InputMultiFile,
		Extensions:  inputExtensions,
		Options: []tool.Option{
			{Key: "width", Label: "Width (pixels)", Type: tool.OptInt, Default: "800", Min: 0,
				Help: "Set to 0 to size automatically from the height."},
			{Key: "height", Label: "Height (pixels)", Type: tool.OptInt, Default: "0", Min: 0,
				Help: "Set to 0 to size automatically from the width."},
		},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if len(req.Inputs) == 0 {
				return tool.Result{}, errors.New("please choose at least one image")
			}
			w := req.Options.Int("width", 0)
			h := req.Options.Int("height", 0)
			if w <= 0 && h <= 0 {
				return tool.Result{}, errors.New("please enter a width or a height (or both)")
			}
			var outputs []string
			for i, in := range req.Inputs {
				if err := cancelled(ctx); err != nil {
					return tool.Result{Outputs: outputs}, err
				}
				img, err := decodeImage(in)
				if err != nil {
					return tool.Result{}, fmt.Errorf("%s: %w", filepath.Base(in), err)
				}
				resized := imaging.Resize(img, w, h, imaging.Lanczos)
				format := formatFromExt(in)
				out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (resized)", extFor(format)))
				if err := encodeImage(resized, out, format, 90); err != nil {
					return tool.Result{}, fmt.Errorf("%s: could not save: %w", filepath.Base(in), err)
				}
				outputs = append(outputs, out)
				p.Update(float64(i+1)/float64(len(req.Inputs)), "Resizing images…")
			}
			return tool.Result{Outputs: outputs}, nil
		},
	}))
}
