package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/moeeinaali/go-blurhash"
)

const (
	readmePath  = "README.md"
	startMarker = "<!-- GENERATED_OUTPUT_MATRIX_START -->"
	endMarker   = "<!-- GENERATED_OUTPUT_MATRIX_END -->"
	outputDir   = "docs/generated/output-matrix"
)

type componentPair struct {
	x int
	y int
}

type renderedCell struct {
	Label string
	Data  string
}

func main() {
	imagePaths := []string{"testdata/1.jpg", "testdata/2.jpg"}
	sizes := []int{8, 128}
	components := []componentPair{
		{x: 1, y: 1},
		{x: 9, y: 1},
		{x: 1, y: 9},
		{x: 9, y: 9},
	}

	content, err := generateContent(imagePaths, sizes, components)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate content failed: %v\n", err)
		os.Exit(1)
	}

	if err := upsertReadme(readmePath, content); err != nil {
		fmt.Fprintf(os.Stderr, "update README failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("README output matrix updated successfully")
}

func generateContent(imagePaths []string, sizes []int, components []componentPair) (string, error) {
	var b strings.Builder
	// b.WriteString(fmt.Sprintf("Generated at: %s UTC\n\n", time.Now().UTC().Format(time.RFC3339)))
	b.WriteString("Inputs:\n")
	b.WriteString("- files: testdata/1.jpg, testdata/2.jpg\n")
	b.WriteString("- sizes (width=height): 8, 128\n")
	b.WriteString("- components: x=1 y=1, x=9 y=1, x=1 y=9, x=9 y=9\n\n")
	b.WriteString("Rendering notes:\n")
	b.WriteString("- each image is decoded at its real size (8x8 or 128x128)\n")
	b.WriteString("- each image is displayed in README at 300x300 pixels\n")
	b.WriteString("- image src points to generated PNG files in docs/generated/output-matrix\n")

	for _, imgPath := range imagePaths {
		img, err := readImage(imgPath)
		if err != nil {
			return "", err
		}

		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("### %s\n\n", imgPath))
		bounds := img.Bounds()
		b.WriteString(fmt.Sprintf("Original image (real size: %dx%d, display: 300x300)\n\n", bounds.Dx(), bounds.Dy()))
		b.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"original %s\" width=\"300\" height=\"300\" style=\"object-fit:contain;\" />\n\n", imgPath, imgPath))

		name := strings.TrimSuffix(filepath.Base(imgPath), filepath.Ext(imgPath))

		for _, size := range sizes {
			row := make([]renderedCell, 0, len(components))
			for _, c := range components {
				hash, err := blurhash.Encode(img, blurhash.WithComponents(c.x, c.y))
				if err != nil {
					return "", fmt.Errorf("encode failed for %s size=%d x=%d y=%d: %w", imgPath, size, c.x, c.y, err)
				}

				decoded, err := blurhash.Decode(hash, size, size)
				if err != nil {
					return "", fmt.Errorf("decode failed for %s size=%d x=%d y=%d: %w", imgPath, size, c.x, c.y, err)
				}

				relPath := filepath.ToSlash(filepath.Join(outputDir, name, fmt.Sprintf("size-%d-x-%d-y-%d.png", size, c.x, c.y)))
				if err := writePNG(relPath, decoded); err != nil {
					return "", fmt.Errorf("write png failed for %s size=%d x=%d y=%d: %w", imgPath, size, c.x, c.y, err)
				}

				_, err = os.Stat(relPath)
				if err != nil {
					return "", fmt.Errorf("verify png failed for %s size=%d x=%d y=%d: %w", imgPath, size, c.x, c.y, err)
				}

				row = append(row, renderedCell{
					Label: fmt.Sprintf("x=%d y=%d", c.x, c.y),
					Data:  relPath,
				})
			}

			b.WriteString(fmt.Sprintf("#### size=%d (real output: %dx%d, display: 300x300)\n\n", size, size, size))
			b.WriteString(renderHTMLGrid([][]renderedCell{row}, imgPath, size))
			b.WriteString("\n")
		}
	}

	return b.String(), nil
}

func renderHTMLGrid(rows [][]renderedCell, imagePath string, size int) string {
	var b strings.Builder
	b.WriteString("<table>\n")
	b.WriteString("  <thead>\n")
	b.WriteString("    <tr><th>Case</th>")
	for _, c := range rows[0] {
		b.WriteString(fmt.Sprintf("<th>%s</th>", c.Label))
	}
	b.WriteString("</tr>\n")
	b.WriteString("  </thead>\n")
	b.WriteString("  <tbody>\n")

	for rowIdx, row := range rows {
		b.WriteString(fmt.Sprintf("    <tr><th>row-%d</th>", rowIdx+1))
		for _, cell := range row {
			alt := fmt.Sprintf("%s size %d %s", imagePath, size, cell.Label)
			b.WriteString(fmt.Sprintf("<td><img src=\"%s\" alt=\"%s\" width=\"300\" height=\"300\" style=\"object-fit:contain; image-rendering: pixelated;\" /></td>", cell.Data, alt))
		}
		b.WriteString("</tr>\n")
	}

	b.WriteString("  </tbody>\n")
	b.WriteString("</table>\n")

	return b.String()
}

func writePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}

func upsertReadme(path, generated string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(raw)

	section := startMarker + "\n" + strings.TrimSpace(generated) + "\n" + endMarker

	start := strings.Index(content, startMarker)
	end := strings.Index(content, endMarker)

	if start >= 0 && end >= 0 && end > start {
		end += len(endMarker)
		updated := content[:start] + section + content[end:]
		return os.WriteFile(path, []byte(updated), 0644)
	}

	var out strings.Builder
	out.WriteString(strings.TrimRight(content, "\n"))
	out.WriteString("\n\n## Generated Output Matrix\n\n")
	out.WriteString("This section is auto-generated by `go run ./scripts/generate_readme_output_matrix.go`.\n\n")
	out.WriteString(section)
	out.WriteString("\n")

	return os.WriteFile(path, []byte(out.String()), 0644)
}
