# BlurHash Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/moeeinaali/go-blurhash.svg)](https://pkg.go.dev/github.com/moeeinaali/go-blurhash)
[![Go Report Card](https://goreportcard.com/badge/github.com/moeeinaali/go-blurhash)](https://goreportcard.com/report/github.com/moeeinaali/go-blurhash)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A fast, efficient, and pure Go implementation of the [BlurHash](https://blurHash.com) image encoding algorithm. BlurHash is a compact representation of a placeholder for an image, allowing you to show a smooth placeholder while loading images.

## Features

- 🚀 **Fast & Efficient** - Pure Go with zero external dependencies
- 📦 **Spec Compliant** - Matches official implementations (C, TypeScript, Swift)
- 🎨 **Flexible** - Configurable component counts, custom punch factors, and auto-sizing
- 🧪 **Well Tested** - Comprehensive test suite with real images and cross-language validation
- 📊 **Streaming Support** - Read/write images directly from files or buffers
- ⚡ **Optimized** - Pre-computed lookup tables, minimal allocations
- 🛠️ **CLI Tool** - Command-line utility for encoding/decoding images
- 📖 **Full Documentation** - Clear examples and API reference

## Installation

```bash
go get github.com/moeeinaali/go-blurhash
```

## Quick Start

### Encoding an Image

```go
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

	"github.com/moeeinaali/go-blurhash"
)

func main() {
	// Open and decode an image
	f, _ := os.Open("photo.jpg")
	img, _, _ := image.Decode(f)
	f.Close()

	// Encode to BlurHash
	hash, _ := blurhash.Encode(img)
	fmt.Println(hash) // e.g., "B~LrYI~c{H?b=::k"
}
```

### Decoding a Hash to an Image

```go
hash := "B~LrYI~c{H?b=::k"

// Decode to a specific size
img, _ := blurhash.Decode(hash, 256, 256)

// Save as PNG
f, _ := os.Create("preview.png")
defer f.Close()
png.Encode(f, img)
```

### Validating a Hash

```go
isValid, reason := blurhash.IsValid("B~LrYI~c{H?b=::k")
if !isValid {
	fmt.Printf("Invalid hash: %s\n", reason)
}
```

## API Reference

### Encoding

```go
// Encode image with default options (4x3 components)
hash, err := blurhash.Encode(img)

// Encode with custom components
hash, err := blurhash.Encode(img, blurhash.WithComponents(6, 5))

// Auto-detect components based on image aspect ratio
hash, err := blurhash.Encode(img, blurhash.WithAutoComponents())

// Limit image size before encoding (useful for large images)
hash, err := blurhash.Encode(img, blurhash.WithMaxSize(100))

// Combine options
hash, err := blurhash.Encode(img,
	blurhash.WithComponents(4, 3),
	blurhash.WithMaxSize(200),
)
```

### Decoding

```go
// Decode to specific dimensions
img, err := blurhash.Decode(hash, width, height)

// Decode with custom punch (contrast multiplier)
img, err := blurhash.Decode(hash, width, height, blurhash.WithPunch(1.5))

// Combine options
img, err := blurhash.Decode(hash, width, height,
	blurhash.WithPunch(1.2),
)
```

### Validation

```go
// Check if a hash is valid
isValid, reason := blurhash.IsValid(hash)

// Examples of reasons:
// - "valid"
// - "hash too short (minimum 6 characters)"
// - "invalid character: !"
// - "error decoding components"
```

### Streaming

```go
// Encode from a file reader
f, _ := os.Open("image.jpg")
defer f.Close()
hash, err := blurhash.EncodeReader(f)

// Decode to a file writer
f, _ := os.Create("preview.png")
defer f.Close()
err := blurhash.DecodeToWriter(f, hash, 256, 256)
```

### Reusable Encoders/Decoders

For better performance when encoding/decoding multiple images, use structured instances:

```go
// Create an encoder with default options
encoder := blurhash.NewEncoder(
	blurhash.WithAutoComponents(),
	blurhash.WithMaxSize(100),
)

// Reuse for multiple images
for _, imgPath := range imageList {
	f, _ := os.Open(imgPath)
	img, _, _ := image.Decode(f)
	hash, _ := encoder.Encode(img)
	fmt.Println(hash)
	f.Close()
}

// Create a decoder
decoder := blurhash.NewDecoder(blurhash.WithPunch(1.2))

// Reuse for multiple hashes
for _, hash := range hashList {
	img, _ := decoder.Decode(hash, 256, 256)
	// Process img
}
```

## Options

| Option | Type | Default | Purpose |
|--------|------|---------|---------|
| `WithComponents(x, y)` | `int, int` | `4, 3` | Set component counts (each 1-9) |
| `WithAutoComponents()` | - | - | Auto-detect components based on aspect ratio |
| `WithMaxSize(pixels)` | `int` | `32` | Maximum image dimension before downscaling |
| `WithPunch(factor)` | `float64` | `1.0` | Contrast multiplier during decoding (>0) |

## CLI Tool

Build the command-line tool:

```bash
go build -o blurhash ./cmd/blurhash
```

### Encoding

```bash
# Encode with default settings (4x3 components)
./blurhash encode photo.jpg

# Encode with custom components
./blurhash encode -x 6 -y 5 photo.jpg

# Auto-detect components
./blurhash encode -auto photo.jpg

# Save hash to file
./blurhash encode photo.jpg -out hash.txt

# With max size
./blurhash encode -maxsize 256 large-photo.jpg
```

### Decoding

```bash
# Decode to 256x256 PNG
./blurhash decode "B~LrYI~c{H?b=::k"

# Custom output size
./blurhash decode -w 512 -h 512 "B~LrYI~c{H?b=::k" -out large.png

# With punch factor
./blurhash decode -punch 1.5 "hash" -out punchy.png
```

### Validation

```bash
# Validate one hash
./blurhash validate "B~LrYI~c{H?b=::k"

# Validate multiple hashes
./blurhash validate "hash1" "hash2" "hash3"
```

## Examples

See `examples/example.go` for more detailed usage examples:

```bash
go run examples/example.go
```

Topics covered:
- Basic encode/decode operations
- Using encoding options
- Working with file I/O
- Hash validation and error handling
- Round-trip testing

## Performance

Benchmarks on Apple M4 (256×256 image):

```
BenchmarkEncode256     1036 ns/op    7.6 KB/op    1043 allocs/op
BenchmarkDecode256    84.8 µs/op    6.6 KB/op      24 allocs/op
```

Performance characteristics:
- **Encoding**: O(width × height × components) - typically 1-2ms for 256×256 images
- **Decoding**: O(width × height × components) - typically 50-200µs for 256×256 output
- **Memory**: Minimal allocations; reusable buffers when using Encoder/Decoder
- **Downscaling**: Images larger than `maxSize` are downscaled before encoding, reducing computation

## Testing

Comprehensive test suite includes:

- **Unit tests** - Base83 encoding, color conversions, DCT, component validation
- **Golden tests** - Output verification against official C implementation
- **Cross-language vectors** - Validation against C and TypeScript implementations
- **Integration tests** - Real image files (solid colors, gradients, photographs)
- **Component variations** - 1×1 through 9×9 configurations
- **Fuzz tests** - Random inputs to ensure robustness
- **Stress tests** - 100+ encode/decode cycles
- **Benchmarks** - Performance measurement suite
- **Streaming** - File I/O operations

Run tests:

```bash
go test ./...              # All tests
go test -run Integration   # Integration tests only
go test -bench .           # Benchmarks
go test -fuzz Fuzz         # Fuzz tests
```

## Technical Details

### Supported Image Formats

Encoding accepts any `image.Image`:
- PNG, JPEG, GIF (via standard library)
- Custom image formats (if they implement `image.Image`)

Decoding produces:
- `image.RGBA` - 32-bit RGBA image

### Component Count

Component counts determine hash length and quality:
- More components = longer hash, higher quality
- X and Y are independent: 4×3, 6×5, 9×9 are all valid
- Minimum 1×1, maximum 9×9 for each dimension

Hash length formula: `4 + 2×(compX×compY)`

Examples:
- 1×1: 6 characters
- 4×3: 28 characters
- 6×5: 64 characters
- 9×9: 166 characters

### Color Space

The algorithm uses linear RGB internally:
- sRGB input (0-255 uint8) → linear RGB (0-1 float)
- DCT applied in linear RGB
- linear RGB → sRGB output (0-255 uint8)

This ensures perceptually accurate results.

## Comparison with Alternatives

| Feature | blurhash-go | Official C | TypeScript | Swift |
|---------|-------------|-----------|-----------|-------|
| Language | Go | C | TypeScript/JS | Swift |
| Dependencies | None | None | None | None |
| Performance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Spec Compliant | ✓ | ✓ | ✓ | ✓ |
| Test Suite | ✓ | - | ✓ | ✓ |

## Error Handling

The library provides descriptive errors:

```go
hash, err := blurhash.Encode(img, blurhash.WithComponents(10, 3))
// Error: invalid component count (outside 1-9 range)

img, err := blurhash.Decode("ab", 256, 256)
// Error: hash too short (minimum 6 characters)

img, err := blurhash.Decode("!!!!!!!", 256, 256)
// Error: invalid character in hash
```

## Streaming Examples

### Encode from HTTP response

```go
resp, _ := http.Get("https://example.com/photo.jpg")
defer resp.Body.Close()

hash, err := blurhash.EncodeReader(resp.Body)
```

### Decode to HTTP response

```go
w.Header().Set("Content-Type", "image/png")
blurhash.DecodeToWriter(w, hash, 256, 256)
```

## Contributing

Contributions welcome! Areas of interest:

- Performance optimizations
- Additional image format support
- More comprehensive examples
- Documentation improvements
- Bug reports and fixes

Please ensure:
- Tests pass: `go test ./...`
- Code follows conventions: `gofmt -w .`
- No new external dependencies without discussion
- Commit messages follow Conventional Commits

## Algorithm Reference

- **BlurHash Specification**: https://blurHash.com
- **Algorithm Paper**: DCT-based image placeholder encoding
- **Official Implementations**: https://github.com/woltapp/blurhash

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [BlurHash](https://woltapp.github.io/blurhash/) - Original algorithm and specification by Wolt
- Reference implementations in C, TypeScript, and Swift
- The Go community for standard library excellence

## Changelog

### v1.0.0 (March 2026)

- ✅ Full BlurHash implementation
- ✅ Encode/Decode API with options pattern
- ✅ CLI tool for command-line usage
- ✅ Streaming support for file I/O
- ✅ Comprehensive test suite
- ✅ Validation utilities
- ✅ Performance benchmarks
- ✅ Complete documentation

## Future Roadmap

- [ ] WASM bindings for web use
- [ ] WebP and other format support
- [ ] Batch processing optimizations
- [ ] Hardware acceleration (SIMD)
- [ ] Integration examples (web frameworks)
- [ ] Metrics and profiling utilities

## FAQ

**Q: How does BlurHash differ from simple thumbnails?**
A: BlurHash creates a compact (6-200 character) hash that decodes to any size smoothly, while thumbnails are full images. BlurHash is ideal for placeholders during image loading.

**Q: Can I use a BlurHash from another language?**
A: Yes! BlurHash is standardized. Hashes created by Python/JavaScript/Swift libraries work with this Go implementation and vice versa.

**Q: What's the recommended component count?**
A: Start with 4×3 (28 characters) for balanced quality/size. Use 6×5 (64 chars) for higher quality, 1×1 (6 chars) for minimal size.

**Q: Does it work with animated images (GIF)?**
A: Yes, but only the first frame is encoded. The library accepts any `image.Image`, including GIF frames.

**Q: Can I use this in production?**
A: Absolutely! The library is fully tested, spec-compliant, and handles all inputs safely (no panics).

---

**Made with ❤️ for the Go community**

Questions? [Open an issue](https://github.com/moeeinaali/go-blurhash/issues) or [start a discussion](https://github.com/moeeinaali/go-blurhash/discussions).

## Generated Output Matrix

This section is auto-generated by `go run ./scripts/generate_readme_output_matrix.go`.

<!-- GENERATED_OUTPUT_MATRIX_START -->
Inputs:
- files: testdata/1.jpg, testdata/2.jpg
- sizes (width=height): 8, 128
- components: x=1 y=1, x=9 y=1, x=1 y=9, x=9 y=9

Rendering notes:
- each image is decoded at its real size (8x8 or 128x128)
- each image is displayed in README at 300x300 pixels
- image src is embedded as base64 PNG data URL

### testdata/1.jpg

Original image (real size: 1170x658, display: 300x300)

<img src="testdata/1.jpg" alt="original testdata/1.jpg" width="300" height="300" style="object-fit:contain;" />

#### size=8 (real output: 8x8, display: 300x300)

<table>
  <thead>
    <tr><th>Case</th><th>x=1 y=1</th><th>x=9 y=1</th><th>x=1 y=9</th><th>x=9 y=9</th></tr>
  </thead>
  <tbody>
    <tr><th>row-1</th><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAGUlEQVR4nGJZOncKAzbAhFV00EoAAgAA//91BAHpJHRsSAAAAABJRU5ErkJggg==" alt="testdata/1.jpg size 8 x=1 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAALklEQVR4nGL58unLyTPHrpzdLyIu8/39p/+Mfxh/fnv/4TcTAw4wOCUAAQAA///XERGNCh9qHAAAAABJRU5ErkJggg==" alt="testdata/1.jpg size 8 x=9 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAOUlEQVR4nGK5fnIfAzbAcubcRewSXPx82CW4OTiwS7BzcmOXOHzwBHYJTmFO7BJsLD+wSgACAAD//3jUCLqkGhNIAAAAAElFTkSuQmCC" alt="testdata/1.jpg size 8 x=1 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAA10lEQVR4nATAv06DQBgA8OP47NHKn0oLBROIm4OLUQcH39XZxEdwramJibo4ONRFWwMKnvYOe8fd5w8Qcfl0z6v3OEuY4gqZlm3TuPDVtC+r6vnhdlaUsuUeWPG9qVeaMgJ5vh8yPxnPrPyNgnE2ScuDlPpx0JnexA4OsDw625DtYMJ+VA2oTBIMR9qb7kZ2x3pO0X2u9/wE/ox7M18s7uaPr8vqDSFsCiLXH4RuFc/yaOhrStU0rAJ3xEl/cnEO1tFS6NPD41oo3VkXHNGnl1fX/wEAAP//n6NkWSsPJFYAAAAASUVORK5CYII=" alt="testdata/1.jpg size 8 x=9 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td></tr>
  </tbody>
</table>

#### size=128 (real output: 128x128, display: 300x300)

<table>
  <thead>
    <tr><th>Case</th><th>x=1 y=1</th><th>x=9 y=1</th><th>x=1 y=9</th><th>x=9 y=9</th></tr>
  </thead>
  <tbody>
    <tr><th>row-1</th><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAABNUlEQVR4nOzRMQ0AIADAMELwf/Fy4hMZPVgVLNm6Z484Uwf8rgFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AHsBAAD//8PhAtl1ZPcpAAAAAElFTkSuQmCC" alt="testdata/1.jpg size 128 x=1 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAAB0klEQVR4nOzRW1LFIBAA0SG6/8Va90deY0Fe7qB/+lydQhMEq78/P5+IyG3mnGN9xhh99K211luvbamt1vrM9ttq7a32Xntro7cx2hx9zpFzZM6ILLGUUo5jj+PY42sv9izH+WwtStnv7k0Z58zM63KZMc9bzvvHmc+j/b1nxLs/4rrB84n/s9yvrUVe562vyPfcuRbzOWc+vztvdO2495X9Z665L3L+P+uEct+mxOsIoQwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDMADADwAwAMwDsLwAA//94xcGB+dRoUQAAAABJRU5ErkJggg==" alt="testdata/1.jpg size 128 x=9 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAAB3klEQVR4nOzd23WFIBQGYZDpv6S8p4IUI0rKmAfnq4B1/rVFLmfL3+/PiIczjj2GT6Pf38X7vvYYPq0AZDwFoOJ5HnsMn8azC8DE3tsew6exqwAVdxWg6hEk475vewyfVgXImoRlVYCM3UpYxfO0F2RqM07G+3YiY+pIUnbZA/g6xpj2GD6NOQvAxLwKwMR1NQ2YuFYBmLiuZY/h06oAGas5QFUFyHoLkhWArABkXK2EVe0FyQpARr+/i9l5gKoKkPUOKisAWQHImN1KUVUBsgKQFYCM0UJA1UpY1iNI1m6orABkzNlTyFQFyApA1iNIVgXICkDWVoSsCpD1L0lZr0CyApDVqkDGOQVg4tQtRVX3dFkByGpZJqt9vay2lbICkBWAjH0XgIm7AFS1r5cVgKxHkKwKkPURHxn7bivCxH6qAFO7obICkPF2JqzqVoSsm3Gy+gXJmLWtVNW6WFYAsgKQsSgAE6sKULHoEyam5gAZa7USM1UBMlbfD1D1KUMZbQW52g2VFYCsAGSMjiRVBSBjdCtCxRndCzJVATJO/5RX1S1FVgXIqgBZ/YJkBSDjLQBVFSCjrSBXFSCjt1AXp3WYqjlA1iNIxjgdCpu6FSRrDpDR7++qAmT/AQAA///jza5jKOChgwAAAABJRU5ErkJggg==" alt="testdata/1.jpg size 128 x=1 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,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" alt="testdata/1.jpg size 128 x=9 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td></tr>
  </tbody>
</table>


### testdata/2.jpg

Original image (real size: 1920x1200, display: 300x300)

<img src="testdata/2.jpg" alt="original testdata/2.jpg" width="300" height="300" style="object-fit:contain;" />

#### size=8 (real output: 8x8, display: 300x300)

<table>
  <thead>
    <tr><th>Case</th><th>x=1 y=1</th><th>x=9 y=1</th><th>x=1 y=9</th><th>x=9 y=9</th></tr>
  </thead>
  <tbody>
    <tr><th>row-1</th><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAGUlEQVR4nGLJr61nwAaYsIoOWglAAAAA//8iIQF+uAIVhgAAAABJRU5ErkJggg==" alt="testdata/2.jpg size 8 x=1 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAALUlEQVR4nGKJK8h88YmV+ccPfs7vqmoqUtJq3398v3fvERMDDjA4JQABAAD//w82DHmQdTfhAAAAAElFTkSuQmCC" alt="testdata/2.jpg size 8 x=9 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAOUlEQVR4nGJxy53AgA2wiHKxY5dQU5HBLiHD9gG7xPdPT7BLXH3wHLvEk6e7sEv85PiMVQIQAAD//1bPDOJ6e4/HAAAAAElFTkSuQmCC" alt="testdata/2.jpg size 8 x=1 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAA2ElEQVR4nADIADf/BC9Fh9ECDgAL4lIbEQkOBQL5+w0SIsXg5wQDFg4Z7vLn7gviFQQ2CQ8jIBL5A/04/QkESDQZ4QQm8+rqEAL4JCIUGhET/voB8/cIBP/l0hIMFhv52goHCAj1Bh4nIfz38R4SCgL2/vL15NL58vkE/vMGAL7/+O0REgmzr7IDHx8Xutbp/wP02MruLTxF/QYNo6mO7vL5Ayko2PoB8O/e4ggf1PDr2/Lw3P79BuLnyAQCCQDoB+7pCAAA/wAdD/n9/QH39fH9BwkBAAD//+h0XkHx9QKOAAAAAElFTkSuQmCC" alt="testdata/2.jpg size 8 x=9 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td></tr>
  </tbody>
</table>

#### size=128 (real output: 128x128, display: 300x300)

<table>
  <thead>
    <tr><th>Case</th><th>x=1 y=1</th><th>x=9 y=1</th><th>x=1 y=9</th><th>x=9 y=9</th></tr>
  </thead>
  <tbody>
    <tr><th>row-1</th><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAABNElEQVR4nOzRQQkAIADAQBH7mkCsbox7uEsw2NrnjjhTB/yuAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAawDWAKwBWAOwBmANwBqANQBrANYArAFYA7AGYA3AGoA1AGsA1gCsAVgDsAZgDcAagDUAewEAAP//SmYCbstm7VYAAAAASUVORK5CYII=" alt="testdata/2.jpg size 128 x=1 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAAB1ElEQVR4nOzRQU7DMBBA0Rm752TPjvtL8aCxWwtO8Df/KZ2YQGyV//r6+Y6I6k9WRFWsyJ59Zc+1532y4vy2+q+j+tav5Wdmz3tFRo3e/s/DPifOot/Kyox7jRGZOca92pxnzr2Y872ac756jDnma89+OMbMsV/I2VvsPbJ375nRx/S5fXTur3G+ftXnU7Xq3NYeq1bPvXir59yfdRf/f171POeFWk+dxd5mH7D/ifvUEUIZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYAmAFgBoAZAGYA2G8AAAD//wZfqyim+s12AAAAAElFTkSuQmCC" alt="testdata/2.jpg size 128 x=9 y=1" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAAB3ElEQVR4nOzd25HkIBQEUUBp23qzvoy/I2DMyA/lsYDoirqo9eTf/58Rz7IX8HUFIGPaK/i4ApAxSkBFm4CrESRjloCqESRjVQFVI0hWA2QFICsAGWsVgIk1OxA18dQAFWvVAFMByFhPAZh4aoCqESRrBMkaQbIaIGsPkDWCZI0gWQ2QdTpa1tlQWQ2QdUlSRgdBrhogoy3AxZrXXsOn8dQAVQ2QMUcBmFgFoGKNY6/h0wpAxiwAFfMWgKkRJKsBMkYNUDHvttfwaY0gGaMGqBg1QFUAsgKQFYCM2/8AFbcGqBi3CzKmGiDj1AAVtwBU3C7KqxpBsgKQcfr9VTVAVgAydv/DVDVAxm4XVtUAGfu0CZgaQbIaIGN3HKqqAbIaIGPvbk00FYCsAGS8BaDi7Ndew6exC0BVALICkLUHyDj7117Dp3FODTBxTg0wcW8BmDgFoOLe9gATY9QAE7cAVIzRCDIxZg0wMWYNMDFWF2RMjFkApgKQMVb3BZmYvb5exVgFYOL2BQ1VH7GSMfqMlYo7CsBEW7CrBsi4HQap6DFhFz0i5up9QTJ6PsPVCJKxdzPIVACyApCx3wIwsd92YRNvI0hVA2RtwrJeVyOrATJOZ+NUfUFD9hcAAP//gQHHWuDRVfUAAAAASUVORK5CYII=" alt="testdata/2.jpg size 128 x=1 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td><td><img src="data:image/png;base64,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" alt="testdata/2.jpg size 128 x=9 y=9" width="300" height="300" style="object-fit:contain; image-rendering: pixelated;" /></td></tr>
  </tbody>
</table>
<!-- GENERATED_OUTPUT_MATRIX_END -->
