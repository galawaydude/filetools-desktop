package image_test

import (
	"context"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/disintegration/imaging"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func makeImg(t *testing.T, dir, name string, w, h int) string {
	t.Helper()
	img := imaging.New(w, h, color.NRGBA{120, 180, 90, 255})
	path := filepath.Join(dir, name)
	if err := imaging.Save(img, path); err != nil {
		t.Fatalf("makeImg: %v", err)
	}
	return path
}

func run(t *testing.T, id string, req tool.Request) tool.Result {
	t.Helper()
	tl, ok := tool.Default.Get(id)
	if !ok {
		t.Fatalf("tool %q not registered", id)
	}
	res, err := tl.Run(context.Background(), req, tool.NopProgress)
	if err != nil {
		t.Fatalf("%s: %v", id, err)
	}
	return res
}

func decoded(t *testing.T, path string) image.Image {
	t.Helper()
	img, err := imaging.Open(path)
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return img
}

func TestConvertToEachFormat(t *testing.T) {
	dir := t.TempDir()
	src := makeImg(t, dir, "src.png", 40, 30)
	for _, format := range []string{"JPG", "PNG", "WEBP", "BMP", "TIFF", "GIF"} {
		out := t.TempDir()
		res := run(t, "image.convert", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"format": format}})
		if len(res.Outputs) != 1 {
			t.Fatalf("%s: want 1 output, got %d", format, len(res.Outputs))
		}
		// Every produced file must decode back to a valid image.
		img := decoded(t, res.Outputs[0])
		if img.Bounds().Dx() != 40 {
			t.Fatalf("%s: width = %d, want 40", format, img.Bounds().Dx())
		}
	}
}

func TestConvertNeverOverwrites(t *testing.T) {
	dir := t.TempDir()
	src := makeImg(t, dir, "src.png", 20, 20)
	out := t.TempDir()
	a := run(t, "image.convert", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"format": "JPG"}})
	b := run(t, "image.convert", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"format": "JPG"}})
	if a.Outputs[0] == b.Outputs[0] {
		t.Fatalf("second convert reused %q (would overwrite)", a.Outputs[0])
	}
}

func TestResizeKeepsAspect(t *testing.T) {
	dir := t.TempDir()
	src := makeImg(t, dir, "src.png", 200, 100)
	out := t.TempDir()
	res := run(t, "image.resize", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"width": "100", "height": "0"}})
	img := decoded(t, res.Outputs[0])
	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 50 {
		t.Fatalf("resized to %dx%d, want 100x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestCompressJPGShrinks(t *testing.T) {
	dir := t.TempDir()
	// A detailed image so JPEG quality matters.
	big := imaging.New(300, 300, color.NRGBA{0, 0, 0, 255})
	for y := 0; y < 300; y++ {
		for x := 0; x < 300; x++ {
			big.Set(x, y, color.NRGBA{uint8(x), uint8(y), uint8(x * y), 255})
		}
	}
	src := filepath.Join(dir, "src.jpg")
	if err := imaging.Save(big, src, imaging.JPEGQuality(100)); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	res := run(t, "image.compress", tool.Request{Inputs: []string{src}, OutDir: out, Options: tool.Options{"quality": "30"}})
	before, _ := os.Stat(src)
	after, _ := os.Stat(res.Outputs[0])
	if after.Size() >= before.Size() {
		t.Fatalf("compressed size %d not smaller than %d", after.Size(), before.Size())
	}
}

func TestCrop(t *testing.T) {
	dir := t.TempDir()
	src := makeImg(t, dir, "src.png", 100, 80)
	out := t.TempDir()
	res := run(t, "image.crop", tool.Request{Inputs: []string{src}, OutDir: out,
		Options: tool.Options{"x": "10", "y": "10", "width": "50", "height": "40"}})
	img := decoded(t, res.Outputs[0])
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 40 {
		t.Fatalf("cropped to %dx%d, want 50x40", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestBatchConvertAndShrink(t *testing.T) {
	dir := t.TempDir()
	a := makeImg(t, dir, "a.png", 500, 400)
	b := makeImg(t, dir, "b.png", 120, 100)
	out := t.TempDir()
	res := run(t, "image.batch", tool.Request{Inputs: []string{a, b}, OutDir: out,
		Options: tool.Options{"format": "JPG", "maxwidth": "200"}})
	if len(res.Outputs) != 2 {
		t.Fatalf("batch produced %d files, want 2", len(res.Outputs))
	}
	for _, o := range res.Outputs {
		if !strings.HasSuffix(o, ".jpg") {
			t.Fatalf("output %q is not .jpg", o)
		}
		if w := decoded(t, o).Bounds().Dx(); w > 200 {
			t.Fatalf("output width %d exceeds max 200", w)
		}
	}
}
