package blurhash

const (
	defaultMaxSize = 32
	defaultPunch   = 1.0
)

type mode uint8

const (
	modeManual mode = iota
	modeAuto
)

type options struct {
	compX   int
	compY   int
	mode    mode
	maxSize int
	punch   float64
}

func defaultOptions() options {
	return options{
		compX:   4,
		compY:   3,
		mode:    modeManual,
		maxSize: defaultMaxSize,
		punch:   defaultPunch,
	}
}

// Option configures encoding/decoding behavior.
type Option func(*options)

// WithComponents sets explicit BlurHash component counts for encoding.
func WithComponents(x, y int) Option {
	return func(o *options) {
		o.compX = x
		o.compY = y
		o.mode = modeManual
	}
}

// WithAutoComponents enables automatic component selection for encoding.
func WithAutoComponents() Option {
	return func(o *options) {
		o.mode = modeAuto
	}
}

// WithMaxSize sets the maximum image side used during encoding downscale.
func WithMaxSize(size int) Option {
	return func(o *options) {
		if size > 0 {
			o.maxSize = size
		}
	}
}

// WithPunch sets the decode punch factor (contrast multiplier for AC terms).
func WithPunch(p float64) Option {
	return func(o *options) {
		if p > 0 {
			o.punch = p
		}
	}
}
