package blurhash

import "math"

type linearColor struct {
	r float64
	g float64
	b float64
}

func sRGBToLinear(v uint8) float64 {
	fv := float64(v) / 255.0
	if fv <= 0.04045 {
		return fv / 12.92
	}
	return math.Pow((fv+0.055)/1.055, 2.4)
}

func linearTosRGB(value float64) uint8 {
	v := value
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	if v <= 0.0031308 {
		return uint8(v*12.92*255 + 0.5)
	}
	return uint8((1.055*math.Pow(v, 1.0/2.4)-0.055)*255 + 0.5)
}

func signPow(v, exp float64) float64 {
	if v < 0 {
		return -math.Pow(-v, exp)
	}
	return math.Pow(v, exp)
}

func encodeDC(c linearColor) int {
	r := int(linearTosRGB(c.r))
	g := int(linearTosRGB(c.g))
	b := int(linearTosRGB(c.b))
	return (r << 16) + (g << 8) + b
}

func decodeDC(v int) linearColor {
	r := uint8(v >> 16)
	g := uint8((v >> 8) & 255)
	b := uint8(v & 255)
	return linearColor{r: sRGBToLinear(r), g: sRGBToLinear(g), b: sRGBToLinear(b)}
}

func encodeAC(c linearColor, maximumValue float64) int {
	quantR := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(c.r/maximumValue, 0.5)*9+9.5)))))
	quantG := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(c.g/maximumValue, 0.5)*9+9.5)))))
	quantB := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(c.b/maximumValue, 0.5)*9+9.5)))))
	return quantR*19*19 + quantG*19 + quantB
}

func decodeAC(v int, maximumValue float64) linearColor {
	quantR := v / (19 * 19)
	quantG := (v / 19) % 19
	quantB := v % 19
	return linearColor{
		r: signPow(float64(quantR-9)/9.0, 2) * maximumValue,
		g: signPow(float64(quantG-9)/9.0, 2) * maximumValue,
		b: signPow(float64(quantB-9)/9.0, 2) * maximumValue,
	}
}
