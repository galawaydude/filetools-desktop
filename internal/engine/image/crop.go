package image

import (
	"context"
	"fmt"
	"image"

	"github.com/disintegration/imaging"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "image.crop",
		Name:        "Crop image",
		Description: "Cut out a rectangle from an image. Give the top-left corner and the size to keep.",
		Category:    tool.CategoryImage,
		Input:       tool.InputSingleFile,
		Extensions:  inputExtensions,
		Options: []tool.Option{
			{Key: "x", Label: "Left (pixels from left)", Type: tool.OptInt, Default: "0", Min: 0},
			{Key: "y", Label: "Top (pixels from top)", Type: tool.OptInt, Default: "0", Min: 0},
			{Key: "width", Label: "Width to keep (pixels)", Type: tool.OptInt, Default: "0", Min: 1},
			{Key: "height", Label: "Height to keep (pixels)", Type: tool.OptInt, Default: "0", Min: 1},
		},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			in := req.Inputs[0]
			x := req.Options.Int("x", 0)
			y := req.Options.Int("y", 0)
			w := req.Options.Int("width", 0)
			h := req.Options.Int("height", 0)
			if w <= 0 || h <= 0 {
				return tool.Result{}, fmt.Errorf("please enter a width and height greater than 0")
			}
			if err := cancelled(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.3, "Cropping…")
			img, err := decodeImage(in)
			if err != nil {
				return tool.Result{}, err
			}
			b := img.Bounds()
			if x >= b.Dx() || y >= b.Dy() {
				return tool.Result{}, fmt.Errorf("the top-left corner is outside the image (image is %d×%d pixels)", b.Dx(), b.Dy())
			}
			rect := image.Rect(x, y, min(x+w, b.Dx()), min(y+h, b.Dy()))
			cropped := imaging.Crop(img, rect)
			format := formatFromExt(in)
			out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (cropped)", extFor(format)))
			if err := encodeImage(cropped, out, format, 90); err != nil {
				return tool.Result{}, fmt.Errorf("could not save the cropped image: %w", err)
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))
}
