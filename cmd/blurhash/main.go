package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/moeeinaali/go-blurhash"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "encode":
		cmdEncode()
	case "decode":
		cmdDecode()
	case "validate":
		cmdValidate()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func cmdEncode() {
	fs := flag.NewFlagSet("encode", flag.ExitOnError)
	compXPtr := fs.Int("x", 4, "X component count (1-9)")
	compYPtr := fs.Int("y", 3, "Y component count (1-9)")
	maxSizePtr := fs.Int("maxsize", 32, "Maximum image dimension before downscaling")
	autoPtr := fs.Bool("auto", false, "Auto-detect component counts based on aspect ratio")
	outputPtr := fs.String("out", "", "Output hash to file (optional)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: blurhash encode [OPTIONS] <image_path> [output_hash_file]\n\n")
		fmt.Fprintf(os.Stderr, "Encode an image to a BlurHash\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  blurhash encode photo.jpg\n")
		fmt.Fprintf(os.Stderr, "  blurhash encode -x 6 -y 4 photo.jpg hash.txt\n")
		fmt.Fprintf(os.Stderr, "  blurhash encode -auto photo.jpg\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(os.Args[2:])

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing image path\n")
		fs.Usage()
		os.Exit(1)
	}

	imagePath := fs.Arg(0)
	outputFile := ""
	if fs.NArg() > 1 {
		outputFile = fs.Arg(1)
	}
	if *outputPtr != "" {
		outputFile = *outputPtr
	}

	// Open and decode image
	f, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening image: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
		os.Exit(1)
	}

	// Encode
	var opts []blurhash.Option
	opts = append(opts, blurhash.WithMaxSize(*maxSizePtr))

	if *autoPtr {
		opts = append(opts, blurhash.WithAutoComponents())
	} else {
		opts = append(opts, blurhash.WithComponents(*compXPtr, *compYPtr))
	}

	hash, err := blurhash.Encode(img, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding: %v\n", err)
		os.Exit(1)
	}

	// Output
	bounds := img.Bounds()
	info := fmt.Sprintf("Image: %s (%s)\nDimensions: %dx%d\nHash: %s\nLength: %d\n",
		filepath.Base(imagePath), format, bounds.Dx(), bounds.Dy(), hash, len(hash))

	fmt.Print(info)

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(hash), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Hash saved to: %s\n", outputFile)
	}
}

func cmdDecode() {
	fs := flag.NewFlagSet("decode", flag.ExitOnError)
	widthPtr := fs.Int("w", 256, "Output image width")
	heightPtr := fs.Int("h", 256, "Output image height")
	punchPtr := fs.Float64("punch", 1.0, "Contrast multiplier (>0)")
	outputPtr := fs.String("out", "output.png", "Output image file")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: blurhash decode [OPTIONS] <hash>\n\n")
		fmt.Fprintf(os.Stderr, "Decode a BlurHash to an image\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  blurhash decode \"B~LrYI~c{H?b=::k\"\n")
		fmt.Fprintf(os.Stderr, "  blurhash decode -w 512 -h 512 \"B~LrYI~c{H?b=::k\" -out decoded.png\n")
		fmt.Fprintf(os.Stderr, "  blurhash decode -punch 1.5 \"hash\" -out out.png\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(os.Args[2:])

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing hash\n")
		fs.Usage()
		os.Exit(1)
	}

	hash := fs.Arg(0)

	// Validate hash
	if valid, reason := blurhash.IsValid(hash); !valid {
		fmt.Fprintf(os.Stderr, "Invalid hash: %s\n", reason)
		os.Exit(1)
	}

	// Decode
	img, err := blurhash.Decode(hash, *widthPtr, *heightPtr, blurhash.WithPunch(*punchPtr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
		os.Exit(1)
	}

	// Save image
	f, err := os.Create(*outputPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Decoded %dx%d image\n", *widthPtr, *heightPtr)
	fmt.Printf("✓ Saved to: %s\n", *outputPtr)
}

func cmdValidate() {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: blurhash validate <hash1> [hash2] ...\n\n")
		fmt.Fprintf(os.Stderr, "Validate one or more BlurHash strings\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  blurhash validate \"B~LrYI~c{H?b=::k\"\n")
		fmt.Fprintf(os.Stderr, "  blurhash validate \"hash1\" \"hash2\" \"hash3\"\n\n")
	}

	fs.Parse(os.Args[2:])

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing hash argument\n")
		fs.Usage()
		os.Exit(1)
	}

	allValid := true
	for _, hash := range fs.Args() {
		valid, reason := blurhash.IsValid(hash)
		status := "✓ VALID"
		if !valid {
			status = "✗ INVALID"
			allValid = false
		}
		fmt.Printf("%s  %s  (%s)\n", status, hash, reason)
	}

	if !allValid {
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `BlurHash CLI Tool v1.0.0

Usage: blurhash <command> [options] [arguments]

Commands:
  encode     Encode an image to a BlurHash
  decode     Decode a BlurHash to an image
  validate   Validate BlurHash string(s)
  help       Show this help message

Examples:
  blurhash encode photo.jpg                     # Encode with defaults (4x3 components)
  blurhash encode -x 6 -y 5 photo.jpg out.txt  # Custom components
  blurhash decode "hash_string" -out out.png   # Decode to 256x256 PNG
  blurhash validate "hash_string"               # Check if hash is valid
  
Options vary by command. Use 'blurhash <command> -h' for help on a specific command.

For more information, visit: https://github.com/moeeinaali/go-blurhash
`)
}
