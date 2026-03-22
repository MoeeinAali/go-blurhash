package blurhash

import (
	"image"
	"image/png"
	"io"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// EncodeReader decodes an image from a reader and returns its BlurHash.
func (e *Encoder) EncodeReader(r io.Reader) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}
	return e.Encode(img)
}

// DecodeToWriter decodes a BlurHash and writes it as PNG to a writer.
func (d *Decoder) DecodeToWriter(w io.Writer, hash string, width, height int) error {
	img, err := d.Decode(hash, width, height)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

// EncodeReader decodes an image from a reader and returns its BlurHash.
func EncodeReader(r io.Reader, opts ...Option) (string, error) {
	return NewEncoder(opts...).EncodeReader(r)
}

// DecodeToWriter decodes a BlurHash and writes it as PNG to a writer.
func DecodeToWriter(w io.Writer, hash string, width, height int, opts ...Option) error {
	return NewDecoder(opts...).DecodeToWriter(w, hash, width, height)
}
