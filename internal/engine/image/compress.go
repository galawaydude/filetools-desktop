package image

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/galawaydude/filetools-desktop/internal/platform"
	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "image.compress",
		Name:        "Compress images",
		Description: "Make images smaller by lowering their quality a little. Works best for JPG images.",
		Category:    tool.CategoryImage,
		Input:       tool.InputMultiFile,
		Extensions:  inputExtensions,
		Options: []tool.Option{{
			Key: "quality", Label: "Quality (1–100)", Type: tool.OptInt, Default: "70", Min: 1, Max: 100,
			Help: "Lower means smaller files. 70 is a good balance.",
		}},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if len(req.Inputs) == 0 {
				return tool.Result{}, errors.New("please choose at least one image")
			}
			quality := req.Options.Int("quality", 70)
			if quality < 1 || quality > 100 {
				quality = 70
			}
			var outputs []string
			var sawLossless bool
			for i, in := range req.Inputs {
				if err := cancelled(ctx); err != nil {
					return tool.Result{Outputs: outputs}, err
				}
				img, err := decodeImage(in)
				if err != nil {
					return tool.Result{}, fmt.Errorf("%s: %w", filepath.Base(in), err)
				}
				format := formatFromExt(in)
				if format != "JPG" {
					sawLossless = true
				}
				out := platform.UniqueName(req.OutDir, platform.OutputName(in, " (compressed)", extFor(format)))
				if err := encodeImage(img, out, format, quality); err != nil {
					return tool.Result{}, fmt.Errorf("%s: could not save: %w", filepath.Base(in), err)
				}
				outputs = append(outputs, out)
				p.Update(float64(i+1)/float64(len(req.Inputs)), "Compressing images…")
			}
			msg := ""
			if sawLossless {
				msg = "Some images use a lossless format (like PNG), so the quality setting has little effect on them."
			}
			return tool.Result{Outputs: outputs, Message: msg}, nil
		},
	}))
}
