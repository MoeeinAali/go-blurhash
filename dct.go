package blurhash

import "math"

func computeFactors(rgb []uint8, width, height, compX, compY int) []linearColor {
	factors := make([]linearColor, 0, compX*compY)
	invW := 1.0 / float64(width)
	invH := 1.0 / float64(height)

	for y := 0; y < compY; y++ {
		for x := 0; x < compX; x++ {
			normalisation := 2.0
			if x == 0 && y == 0 {
				normalisation = 1.0
			}

			r, g, b := 0.0, 0.0, 0.0
			for j := 0; j < height; j++ {
				basisY := math.Cos(math.Pi * float64(y*j) * invH)
				row := j * width * 3
				for i := 0; i < width; i++ {
					basis := normalisation * math.Cos(math.Pi*float64(x*i)*invW) * basisY
					o := row + i*3
					r += basis * sRGBToLinear(rgb[o+0])
					g += basis * sRGBToLinear(rgb[o+1])
					b += basis * sRGBToLinear(rgb[o+2])
				}
			}

			scale := 1.0 / float64(width*height)
			factors = append(factors, linearColor{r: r * scale, g: g * scale, b: b * scale})
		}
	}

	return factors
}
