package blurhash

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// Decoder decodes blurhash strings into images.
type Decoder struct {
	opts options
}

// NewDecoder constructs a decoder with options.
func NewDecoder(opts ...Option) *Decoder {
	o := defaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return &Decoder{opts: o}
}

// Decode decodes a blurhash string into an RGBA image.
func (d *Decoder) Decode(hash string, width, height int) (image.Image, error) {
	return decodeHash(hash, width, height, d.opts)
}

// IsValid returns validation state and reason for an invalid blurhash string.
func IsValid(hash string) (bool, string) {
	_, _, err := validateHash(hash)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func decodeHash(hash string, width, height int, opts options) (image.Image, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidDimensions
	}

	numX, numY, err := validateHash(hash)
	if err != nil {
		return nil, err
	}

	quantisedMaximumValue, err := decode83(hash[1:2])
	if err != nil {
		return nil, err
	}
	maximumValue := float64(quantisedMaximumValue+1) / 166.0

	punch := opts.punch
	if punch <= 0 {
		punch = defaultPunch
	}

	colors := make([]linearColor, numX*numY)
	for i := 0; i < len(colors); i++ {
		if i == 0 {
			v, err := decode83(hash[2:6])
			if err != nil {
				return nil, err
			}
			colors[i] = decodeDC(v)
			continue
		}
		start := 4 + i*2
		v, err := decode83(hash[start : start+2])
		if err != nil {
			return nil, err
		}
		colors[i] = decodeAC(v, maximumValue*punch)
	}

	cosX := make([][]float64, numX)
	for i := 0; i < numX; i++ {
		row := make([]float64, width)
		for x := 0; x < width; x++ {
			row[x] = math.Cos(math.Pi * float64(i*x) / float64(width))
		}
		cosX[i] = row
	}
	cosY := make([][]float64, numY)
	for j := 0; j < numY; j++ {
		row := make([]float64, height)
		for y := 0; y < height; y++ {
			row[y] = math.Cos(math.Pi * float64(j*y) / float64(height))
		}
		cosY[j] = row
	}

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b := 0.0, 0.0, 0.0
			for j := 0; j < numY; j++ {
				basisY := cosY[j][y]
				for i := 0; i < numX; i++ {
					basis := cosX[i][x] * basisY
					c := colors[i+j*numX]
					r += c.r * basis
					g += c.g * basis
					b += c.b * basis
				}
			}
			out.SetRGBA(x, y, color.RGBA{
				R: linearTosRGB(r),
				G: linearTosRGB(g),
				B: linearTosRGB(b),
				A: 255,
			})
		}
	}
	return out, nil
}

func validateHash(hash string) (int, int, error) {
	if len(hash) < 6 {
		return 0, 0, ErrHashTooShort
	}
	sizeFlag, err := decode83(hash[0:1])
	if err != nil {
		return 0, 0, err
	}
	numY := sizeFlag/9 + 1
	numX := sizeFlag%9 + 1
	expected := 4 + 2*numX*numY
	if len(hash) != expected {
		return 0, 0, fmt.Errorf("%w: got=%d expected=%d", ErrHashLengthMismatch, len(hash), expected)
	}
	return numX, numY, nil
}
