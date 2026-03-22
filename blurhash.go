package blurhash

import "image"

// Encode produces a BlurHash string for an image.
func Encode(img image.Image, opts ...Option) (string, error) {
	enc := NewEncoder(opts...)
	return enc.Encode(img)
}

// Decode reconstructs an image from a BlurHash string.
func Decode(hash string, width, height int, opts ...Option) (image.Image, error) {
	dec := NewDecoder(opts...)
	return dec.Decode(hash, width, height)
}
