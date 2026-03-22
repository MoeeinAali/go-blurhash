package blurhash

import (
	"image"
	"image/color"
	"testing"
)

func benchmarkImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := uint8((x*5 + y*3) % 256)
			g := uint8((x*7 + y*11) % 256)
			b := uint8((x*13 + y*17) % 256)
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func BenchmarkEncode32(b *testing.B) {
	img := benchmarkImage(32, 32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Encode(img, WithComponents(4, 3), WithMaxSize(0))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode32(b *testing.B) {
	img := benchmarkImage(32, 32)
	hash, err := Encode(img, WithComponents(4, 3), WithMaxSize(0))
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Decode(hash, 32, 32)
		if err != nil {
			b.Fatal(err)
		}
	}
}
