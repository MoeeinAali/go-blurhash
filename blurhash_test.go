package blurhash

import (
	"errors"
	"image"
	"image/color"
	"testing"
)

func mkImage(w, h int, rgb []uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			o := 3 * (x + y*w)
			img.SetRGBA(x, y, color.RGBA{R: rgb[o], G: rgb[o+1], B: rgb[o+2], A: 255})
		}
	}
	return img
}

func TestEncodeGoldenBlack1x1(t *testing.T) {
	img := mkImage(1, 1, []uint8{0, 0, 0})
	hash, err := Encode(img, WithComponents(1, 1), WithMaxSize(0))
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	if hash != "000000" {
		t.Fatalf("hash mismatch: got=%q want=%q", hash, "000000")
	}
}

func TestEncodeCrossLanguageVectors(t *testing.T) {
	tests := []struct {
		name string
		w    int
		h    int
		rgb  []uint8
		cx   int
		cy   int
		want string
	}{
		{
			name: "mix_3x2_c32",
			w:    3,
			h:    2,
			cx:   3,
			cy:   2,
			rgb: []uint8{
				255, 0, 0, 0, 255, 0, 0, 0, 255,
				255, 255, 0, 255, 255, 255, 0, 255, 255,
			},
			want: "B~LrYI~c{H?b=::k",
		},
		{
			name: "mix_2x2_c22",
			w:    2,
			h:    2,
			cx:   2,
			cy:   2,
			rgb: []uint8{
				120, 30, 200, 90, 180, 10,
				60, 80, 220, 240, 130, 40,
			},
			want: "A~HB1Axg%hx2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			img := mkImage(tc.w, tc.h, tc.rgb)
			hash, err := Encode(img, WithComponents(tc.cx, tc.cy), WithMaxSize(0))
			if err != nil {
				t.Fatalf("encode error: %v", err)
			}
			if hash != tc.want {
				t.Fatalf("hash mismatch: got=%q want=%q", hash, tc.want)
			}
		})
	}
}

func TestDecodeValidation(t *testing.T) {
	_, err := Decode("abc", 32, 32)
	if !errors.Is(err, ErrHashTooShort) {
		t.Fatalf("expected ErrHashTooShort, got %v", err)
	}

	_, err = Decode("000000", 0, 32)
	if !errors.Is(err, ErrInvalidDimensions) {
		t.Fatalf("expected ErrInvalidDimensions, got %v", err)
	}
}

func TestIsValid(t *testing.T) {
	ok, reason := IsValid("000000")
	if !ok || reason != "" {
		t.Fatalf("expected valid hash, got ok=%v reason=%q", ok, reason)
	}
	ok, reason = IsValid("bad")
	if ok || reason == "" {
		t.Fatalf("expected invalid hash with reason, got ok=%v reason=%q", ok, reason)
	}
}
