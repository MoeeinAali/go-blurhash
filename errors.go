package blurhash

import "errors"

var (
	// ErrInvalidComponents indicates component counts outside valid range [1,9].
	ErrInvalidComponents = errors.New("blurhash: components must be between 1 and 9")
	// ErrMismatchedImageSize indicates inconsistent image bounds.
	ErrMismatchedImageSize = errors.New("blurhash: invalid image bounds")
	// ErrHashTooShort indicates an invalid hash shorter than minimum length.
	ErrHashTooShort = errors.New("blurhash: hash must be at least 6 characters")
	// ErrHashLengthMismatch indicates hash length does not match encoded component counts.
	ErrHashLengthMismatch = errors.New("blurhash: hash length mismatch for encoded component counts")
	// ErrInvalidDimensions indicates invalid output dimensions for decoding.
	ErrInvalidDimensions = errors.New("blurhash: width and height must be > 0")
	// ErrInvalidCharacter indicates unsupported characters in base83 strings.
	ErrInvalidCharacter = errors.New("blurhash: invalid base83 character")
}
