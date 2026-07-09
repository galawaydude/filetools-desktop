//go:build ignore

// mkicon converts build/appicon.png into build/appicon.ico (a modern PNG-backed
// .ico used by the Windows installer). Run with: go run build/mkicon.go
package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	img, err := imaging.Open("build/appicon.png")
	if err != nil {
		panic(err)
	}
	resized := imaging.Resize(img, 256, 256, imaging.Lanczos)

	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, resized); err != nil {
		panic(err)
	}
	data := pngBuf.Bytes()

	var out bytes.Buffer
	le := binary.LittleEndian
	binary.Write(&out, le, uint16(0))         // reserved
	binary.Write(&out, le, uint16(1))         // type: icon
	binary.Write(&out, le, uint16(1))         // image count
	out.WriteByte(0)                          // width 0 => 256
	out.WriteByte(0)                          // height 0 => 256
	out.WriteByte(0)                          // palette
	out.WriteByte(0)                          // reserved
	binary.Write(&out, le, uint16(1))         // colour planes
	binary.Write(&out, le, uint16(32))        // bits per pixel
	binary.Write(&out, le, uint32(len(data))) // image size
	binary.Write(&out, le, uint32(22))        // offset (6 + 16)
	out.Write(data)

	if err := os.WriteFile("build/appicon.ico", out.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
