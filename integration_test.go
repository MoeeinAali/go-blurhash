package blurhash

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestRealImageEncode tests encoding various real image files
func TestRealImageEncode(t *testing.T) {
	testdata := "testdata"
	files := []string{
		"solid_red.png",
		"checkerboard.png",
		"gradient.png",
		"noisy.png",
		"small.png",
		"large.png",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(testdata, file)
			f, err := os.Open(path)
			if err != nil {
				t.Skipf("Test image not found: %v", err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("Failed to decode image: %v", err)
			}

			hash, err := Encode(img)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			if len(hash) < 6 {
				t.Errorf("Hash too short: %s", hash)
			}

			// Validate hash structure
			isValid, reason := IsValid(hash)
			if !isValid {
				t.Errorf("Invalid hash produced: %v", reason)
			}

			t.Logf("✓ %s -> %s (%d chars)", file, hash, len(hash))
		})
	}
}

// TestRealImageDecodeRoundTrip tests encode->decode round-trip with real images
func TestRealImageDecodeRoundTrip(t *testing.T) {
	testdata := "testdata"
	files := []string{
		"gradient.png",
		"noisy.png",
		"small.png",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(testdata, file)
			f, err := os.Open(path)
			if err != nil {
				t.Skipf("Test image not found: %v", err)
			}
			defer f.Close()

			original, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			// Encode the original image
			hash, err := Encode(original, WithComponents(4, 3))
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// Decode back using original dimensions
			bounds := original.Bounds()
			decoded, err := Decode(hash, bounds.Dx(), bounds.Dy(), WithPunch(1.2))
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Verify decoded image has same dimensions and is RGBA
			decodedBounds := decoded.Bounds()
			if decodedBounds.Dx() != bounds.Dx() || decodedBounds.Dy() != bounds.Dy() {
				t.Errorf("Decoded image size mismatch: expected %dx%d, got %dx%d",
					bounds.Dx(), bounds.Dy(), decodedBounds.Dx(), decodedBounds.Dy())
			}

			// Sample some pixels to ensure they're reasonable
			_, _, _, a := decoded.At(0, 0).RGBA()
			if a == 0 {
				t.Error("Decoded image has zero alpha")
			}

			t.Logf("✓ %s: decoded %dx%d image with hash %s", file,
				decodedBounds.Dx(), decodedBounds.Dy(), hash)
		})
	}
}

// TestVariousComponentSettings tests encoding with different component counts
func TestVariousComponentSettings(t *testing.T) {
	tests := []struct {
		name  string
		compX int
		compY int
	}{
		{"1x1", 1, 1},
		{"2x2", 2, 2},
		{"4x3", 4, 3},
		{"6x5", 6, 5},
		{"9x9", 9, 9},
	}

	path := filepath.Join("testdata", "gradient.png")
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("Test image not found: %v", err)
	}
	defer f.Close()

	img, _, _ := image.Decode(f)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := Encode(img, WithComponents(tc.compX, tc.compY))
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			isValid, _ := IsValid(hash)
			if !isValid {
				t.Errorf("Invalid hash: %s", hash)
			}

			// Decode and verify
			decoded, err := Decode(hash, 256, 256)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if decoded.Bounds().Dx() != 256 || decoded.Bounds().Dy() != 256 {
				t.Errorf("Wrong decoded size")
			}

			t.Logf("✓ %dx%d components -> hash length %d", tc.compX, tc.compY, len(hash))
		})
	}
}

// TestMaxSizeDownscaling tests maxSize parameter
func TestMaxSizeDownscaling(t *testing.T) {
	path := filepath.Join("testdata", "large.png")
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("Test image not found: %v", err)
	}
	defer f.Close()

	img, _, _ := image.Decode(f)

	sizes := []int{32, 64, 128, 256}
	for _, sz := range sizes {
		hash1, _ := Encode(img, WithMaxSize(sz))
		hash2, _ := Encode(img, WithMaxSize(sz*2))

		// Hashes might differ but both should be valid
		valid1, _ := IsValid(hash1)
		valid2, _ := IsValid(hash2)
		if !valid1 || !valid2 {
			t.Errorf("Invalid hash for maxSize %d", sz)
		}

		t.Logf("✓ MaxSize %d -> %s", sz, hash1)
	}
}

// TestPunchFactor tests the punch (contrast) multiplier in decode
func TestPunchFactor(t *testing.T) {
	path := filepath.Join("testdata", "gradient.png")
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("Test image not found: %v", err)
	}
	defer f.Close()

	img, _, _ := image.Decode(f)
	hash, _ := Encode(img)

	punches := []float64{0.5, 1.0, 1.5, 2.0}
	for _, punch := range punches {
		decoded, err := Decode(hash, 256, 256, WithPunch(punch))
		if err != nil {
			t.Fatalf("Decode with punch %f failed: %v", punch, err)
		}

		if decoded.Bounds().Dx() != 256 {
			t.Errorf("Wrong size with punch %f", punch)
		}

		t.Logf("✓ Punch %.1f -> %dx%d", punch, decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}
}

// TestAutoComponentsSelection tests auto component selection
func TestAutoComponentsSelection(t *testing.T) {
	testdata := "testdata"
	files := []struct {
		name string
		path string
	}{
		{"small", "small.png"},
		{"square", "noisy.png"},
		{"large", "large.png"},
	}

	for _, tc := range files {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(testdata, tc.path)
			f, err := os.Open(path)
			if err != nil {
				t.Skipf("Test image not found: %v", err)
			}
			defer f.Close()

			img, _, _ := image.Decode(f)

			manual, _ := Encode(img, WithComponents(4, 3))
			auto, _ := Encode(img, WithAutoComponents())

			if len(manual) == 0 || len(auto) == 0 {
				t.Error("Failed to encode")
			}

			t.Logf("✓ %s: manual=%d chars, auto=%d chars", tc.name, len(manual), len(auto))
		})
	}
}

// TestStreamingEncode tests EncodeReader for file I/O
func TestStreamingEncode(t *testing.T) {
	path := filepath.Join("testdata", "gradient.png")
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("Test image not found: %v", err)
	}
	defer f.Close()

	hash, err := EncodeReader(f)
	if err != nil {
		t.Fatalf("EncodeReader failed: %v", err)
	}

	isValid, _ := IsValid(hash)
	if !isValid {
		t.Errorf("Invalid hash from EncodeReader: %s", hash)
	}

	t.Logf("✓ EncodeReader: %s", hash)
}

// TestDifferentImageConfigs tests various image configurations
func TestDifferentImageConfigs(t *testing.T) {
	configs := []struct {
		name   string
		width  int
		height int
		color  color.Color
	}{
		{"red_100x100", 100, 100, color.RGBA{255, 0, 0, 255}},
		{"blue_50x200", 50, 200, color.RGBA{0, 0, 255, 255}},
		{"green_200x50", 200, 50, color.RGBA{0, 255, 0, 255}},
		{"grayscale_128x128", 128, 128, color.RGBA{128, 128, 128, 255}},
	}

	for _, cfg := range configs {
		t.Run(cfg.name, func(t *testing.T) {
			// Create a solid color image
			img := image.NewRGBA(image.Rect(0, 0, cfg.width, cfg.height))
			for x := 0; x < cfg.width; x++ {
				for y := 0; y < cfg.height; y++ {
					img.Set(x, y, cfg.color)
				}
			}

			hash, err := Encode(img)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			decoded, err := Decode(hash, cfg.width, cfg.height)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if decoded.Bounds().Dx() != cfg.width || decoded.Bounds().Dy() != cfg.height {
				t.Errorf("Size mismatch: expected %dx%d, got %dx%d",
					cfg.width, cfg.height,
					decoded.Bounds().Dx(), decoded.Bounds().Dy())
			}

			t.Logf("✓ %s: %s", cfg.name, hash)
		})
	}
}

// BenchmarkRealImageEncode benchmarks encoding various real images
func BenchmarkRealImageEncode(b *testing.B) {
	path := filepath.Join("testdata", "gradient.png")
	f, _ := os.Open(path)
	img, _, _ := image.Decode(f)
	f.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Encode(img)
	}
}

// BenchmarkRealImageDecode benchmarks decoding to various sizes
func BenchmarkRealImageDecode(b *testing.B) {
	path := filepath.Join("testdata", "gradient.png")
	f, _ := os.Open(path)
	img, _, _ := image.Decode(f)
	f.Close()

	hash, _ := Encode(img)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Decode(hash, 256, 256)
	}
}

// TestEncodeDecode100Times stress tests encode/decode
func TestEncodeDecode100Times(t *testing.T) {
	path := filepath.Join("testdata", "noisy.png")
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("Test image not found: %v", err)
	}
	defer f.Close()

	img, _, _ := image.Decode(f)

	for i := 0; i < 100; i++ {
		hash, err := Encode(img)
		if err != nil {
			t.Fatalf("Iteration %d encode failed: %v", i, err)
		}

		_, err = Decode(hash, 128, 128)
		if err != nil {
			t.Fatalf("Iteration %d decode failed: %v", i, err)
		}

		if i%20 == 0 {
			fmt.Printf("✓ Completed %d iterations\n", i)
		}
	}

	t.Logf("✓ Completed 100 encode/decode cycles")
}
