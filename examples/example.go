//go:build ignore
// +build ignore

// This file demonstrates how to use the blurhash package.
// To run this example: go run examples/example.go

package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/moeeinaali/go-blurhash"
)

func main() {
	// Example 1: Basic encode/decode
	basicExample()

	// Example 2: Using options
	optionsExample()

	// Example 3: Validation
	validationExample()

	// Example 4: Round-trip test
	roundTripExample()
}

// Example 1: Basic encode/decode with an in-memory image
func basicExample() {
	fmt.Println("=== Example 1: Basic Encode/Decode ===")

	// Create a simple test image with a gradient
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			r := uint8(x * 255 / 100)
			g := uint8(y * 255 / 100)
			b := uint8(128)
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Encode the image to a blurhash
	hash, err := blurhash.Encode(img)
	if err != nil {
		fmt.Printf("Error encoding: %v\n", err)
		return
	}

	fmt.Printf("✓ Encoded hash: %s\n", hash)

	// Decode the hash back to an image (at a smaller size)
	decoded, err := blurhash.Decode(hash, 50, 50)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return
	}

	fmt.Printf("✓ Decoded image size: %dx%d\n", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	fmt.Println()
}

// Example 2: Using encoding options
func optionsExample() {
	fmt.Println("=== Example 2: Using Options ===")

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))
	for x := 0; x < 200; x++ {
		for y := 0; y < 150; y++ {
			img.SetRGBA(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 100, 255})
		}
	}

	// Encode with different component counts
	hash1, _ := blurhash.Encode(img, blurhash.WithComponents(4, 3))
	fmt.Printf("✓ 4x3 components: %s (length: %d)\n", hash1, len(hash1))

	hash2, _ := blurhash.Encode(img, blurhash.WithComponents(6, 5))
	fmt.Printf("✓ 6x5 components: %s (length: %d)\n", hash2, len(hash2))

	// Auto-detect components based on aspect ratio
	hash3, _ := blurhash.Encode(img, blurhash.WithAutoComponents())
	fmt.Printf("✓ Auto components: %s (length: %d)\n", hash3, len(hash3))

	// Decode with different punch factors (contrast control)
	fmt.Println("\nDecoding with punch factors:")
	decoded1, _ := blurhash.Decode(hash1, 200, 150, blurhash.WithPunch(1.0))
	fmt.Printf("✓ Punch 1.0: %dx%d\n", decoded1.Bounds().Dx(), decoded1.Bounds().Dy())

	decoded2, _ := blurhash.Decode(hash1, 200, 150, blurhash.WithPunch(2.0))
	fmt.Printf("✓ Punch 2.0 (more contrast): %dx%d\n", decoded2.Bounds().Dx(), decoded2.Bounds().Dy())

	fmt.Println()
}

// Example 3: Validation
func validationExample() {
	fmt.Println("=== Example 3: Hash Validation ===")

	testHashes := []string{
		"B~LrYI~c{H?b=::k", // Valid hash
		"000000",           // Valid (solid black)
		"invalid!",         // Invalid - contains '!'
		"",                 // Invalid - empty
	}

	for _, hash := range testHashes {
		isValid, reason := blurhash.IsValid(hash)
		status := "✓"
		if !isValid {
			status = "✗"
		}
		fmt.Printf("%s %q - %s\n", status, hash, reason)
	}
	fmt.Println()
}

// Example 4: Round-trip encode/decode test
func roundTripExample() {
	fmt.Println("=== Example 4: Round-trip Test ===")

	// Create a simple colorful image
	original := image.NewRGBA(image.Rect(0, 0, 120, 80))
	colors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{255, 255, 0, 255}, // Yellow
	}

	// Fill with color blocks
	for x := 0; x < 120; x++ {
		for y := 0; y < 80; y++ {
			colorIdx := ((x / 30) + (y/20)*4) % len(colors)
			original.SetRGBA(x, y, colors[colorIdx])
		}
	}

	// Encode
	hash, _ := blurhash.Encode(original)
	fmt.Printf("✓ Original image: 120x80 -> hash: %s\n", hash)

	// Decode back to original size
	decoded, _ := blurhash.Decode(hash, 120, 80)
	fmt.Printf("✓ Decoded image: %dx%d\n", decoded.Bounds().Dx(), decoded.Bounds().Dy())

	// Create a preview
	preview, _ := blurhash.Decode(hash, 30, 20)
	fmt.Printf("✓ Created preview: %dx%d\n", preview.Bounds().Dx(), preview.Bounds().Dy())

	fmt.Println()
}

// Note: If you have actual image files, you can load them like this:
// EXAMPLE: Loading from file
//
//   f, _ := os.Open("photo.jpg")
//   defer f.Close()
//   img, format, _ := image.Decode(f)
//   hash, _ := blurhash.Encode(img)
//   fmt.Printf("Photo loaded from %s, encoded as: %s\n", format, hash)
