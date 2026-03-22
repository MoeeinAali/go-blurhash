package blurhash

import (
	"fmt"
	"image"
	"math"
	"strings"
)

// Encoder encodes images into blurhash strings.
type Encoder struct {
	opts options
}

// NewEncoder constructs an encoder with options.
func NewEncoder(opts ...Option) *Encoder {
	o := defaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return &Encoder{opts: o}
}

// Encode encodes an image into a BlurHash string.
func (e *Encoder) Encode(img image.Image) (string, error) {
	if img == nil {
		return "", fmt.Errorf("%w: nil image", ErrMismatchedImageSize)
	}
	return encodeImage(img, e.opts)
}

func encodeImage(img image.Image, opts options) (string, error) {
	rgb, width, height, err := prepareRGB(img, opts.maxSize)
	if err != nil {
		return "", err
	}

	compX, compY := opts.compX, opts.compY
	if opts.mode == modeAuto {
		compX, compY = autoComponents(width, height)
	}
	if compX < 1 || compX > 9 || compY < 1 || compY > 9 {
		return "", ErrInvalidComponents
	}

	factors := computeFactors(rgb, width, height, compX, compY)
	dc := factors[0]
	ac := factors[1:]

	var b strings.Builder
	b.Grow(2 + 4 + len(ac)*2)

	sizeFlag := (compX - 1) + (compY-1)*9
	b.WriteString(encode83(sizeFlag, 1))

	maximumValue := 1.0
	if len(ac) > 0 {
		actualMaximumValue := 0.0
		for _, c := range ac {
			actualMaximumValue = math.Max(actualMaximumValue, math.Abs(c.r))
			actualMaximumValue = math.Max(actualMaximumValue, math.Abs(c.g))
			actualMaximumValue = math.Max(actualMaximumValue, math.Abs(c.b))
		}
		quantisedMaximumValue := int(math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5))))
		maximumValue = float64(quantisedMaximumValue+1) / 166.0
		b.WriteString(encode83(quantisedMaximumValue, 1))
	} else {
		b.WriteString(encode83(0, 1))
	}

	b.WriteString(encode83(encodeDC(dc), 4))
	for _, c := range ac {
		b.WriteString(encode83(encodeAC(c, maximumValue), 2))
	}

	return b.String(), nil
}

func autoComponents(width, height int) (int, int) {
	if width <= 0 || height <= 0 {
		return 4, 3
	}
	if width >= height {
		x := 4
		y := int(math.Round((float64(height) / float64(width)) * float64(x)))
		if y < 1 {
			y = 1
		}
		if y > 9 {
			y = 9
		}
		return x, y
	}
	y := 4
	x := int(math.Round((float64(width) / float64(height)) * float64(y)))
	if x < 1 {
		x = 1
	}
	if x > 9 {
		x = 9
	}
	return x, y
}

func prepareRGB(img image.Image, maxSize int) ([]uint8, int, int, error) {
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()
	if width <= 0 || height <= 0 {
		return nil, 0, 0, ErrMismatchedImageSize
	}

	if maxSize > 0 && (width > maxSize || height > maxSize) {
		scale := float64(maxSize) / float64(width)
		if height > width {
			scale = float64(maxSize) / float64(height)
		}
		nw := int(math.Round(float64(width) * scale))
		nh := int(math.Round(float64(height) * scale))
		if nw < 1 {
			nw = 1
		}
		if nh < 1 {
			nh = 1
		}
		width = nw
		height = nh
	}

	rgb := make([]uint8, width*height*3)
	for y := 0; y < height; y++ {
		srcY := b.Min.Y + y*b.Dy()/height
		for x := 0; x < width; x++ {
			srcX := b.Min.X + x*b.Dx()/width
			r, g, bl, _ := img.At(srcX, srcY).RGBA()
			o := 3 * (x + y*width)
			rgb[o+0] = uint8(r >> 8)
			rgb[o+1] = uint8(g >> 8)
			rgb[o+2] = uint8(bl >> 8)
		}
	}

	return rgb, width, height, nil
}
