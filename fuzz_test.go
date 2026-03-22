package blurhash

import (
	"image"
	"image/color"
	"testing"
)

func FuzzEncodeDecodeNoPanic(f *testing.F) {
	f.Add(uint8(32), uint8(32), uint8(3), uint8(3))
	f.Add(uint8(7), uint8(5), uint8(2), uint8(2))

	f.Fuzz(func(t *testing.T, w, h, cx, cy uint8) {
		width := int(w%32) + 1
		height := int(h%32) + 1
		compX := int(cx%9) + 1
		compY := int(cy%9) + 1

		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				v := uint8((x*37 + y*17 + width + height) % 256)
				img.SetRGBA(x, y, color.RGBA{R: v, G: v / 2, B: 255 - v, A: 255})
			}
		}

		hash, err := Encode(img, WithComponents(compX, compY), WithMaxSize(0))
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		out, err := Decode(hash, width, height)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if out.Bounds().Dx() != width || out.Bounds().Dy() != height {
			t.Fatalf("size mismatch after decode")
		}
	})
}
