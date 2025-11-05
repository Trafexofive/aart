package converter

import (
	"fmt"
	"image"
	"image/gif"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

type Options struct {
	Width  int
	Height int
	FPS    int
	Method string
	Chars  string
}

type Frame struct {
	Width  int
	Height int
	Cells  [][]Cell
	Delay  int // milliseconds
}

type Cell struct {
	Char rune
	FG   string
	BG   string
}

// ConvertGifToFrames converts a GIF (from URL or file) to ASCII frames
func ConvertGifToFrames(source string, opts Options) ([]*Frame, error) {
	// Load GIF
	gifData, err := loadGif(source)
	if err != nil {
		return nil, fmt.Errorf("failed to load GIF: %w", err)
	}

	// Convert each frame with proper composition
	frames := make([]*Frame, len(gifData.Image))
	var previousFrame image.Image
	
	for i, img := range gifData.Image {
		// Handle GIF disposal and composition
		var composited image.Image
		
		if i == 0 || gifData.Disposal[i-1] == gif.DisposalBackground {
			// First frame or previous frame disposed - use frame as-is
			composited = img
		} else {
			// Composite current frame over previous
			composited = compositeImages(previousFrame, img)
		}
		
		// Save for next iteration
		previousFrame = composited
		
		// Resize image to target dimensions
		resized := resize.Resize(uint(opts.Width), uint(opts.Height), composited, resize.Lanczos3)
		
		// Convert to ASCII
		frame := convertImageToASCII(resized, opts)
		
		// Calculate delay in milliseconds
		delay := gifData.Delay[i] * 10 // GIF delay is in 100ths of a second
		if delay == 0 {
			delay = 1000 / opts.FPS // Use target FPS if no delay specified
		}
		frame.Delay = delay
		
		frames[i] = frame
	}

	return frames, nil
}

// compositeImages composites src over dst, handling transparency
func compositeImages(dst, src image.Image) image.Image {
	if dst == nil {
		return src
	}
	
	bounds := dst.Bounds()
	srcBounds := src.Bounds()
	
	// Create RGBA image for composition
	composited := image.NewRGBA(bounds)
	
	// Copy dst to composited
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			composited.Set(x, y, dst.At(x, y))
		}
	}
	
	// Composite src on top, respecting transparency
	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
			srcColor := src.At(x, y)
			_, _, _, a := srcColor.RGBA()
			
			// Only draw non-transparent pixels
			if a > 0 {
				composited.Set(x, y, srcColor)
			}
		}
	}
	
	return composited
}

// loadGif loads a GIF from URL or local file
func loadGif(source string) (*gif.GIF, error) {
	var reader io.ReadCloser
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		// Load from URL
		resp, err := http.Get(source)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL: %w", err)
		}
		reader = resp.Body
	} else {
		// Load from file
		reader, err = os.Open(source)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
	}
	defer reader.Close()

	gifData, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode GIF: %w", err)
	}

	return gifData, nil
}

// convertImageToASCII converts a single image to ASCII
func convertImageToASCII(img image.Image, opts Options) *Frame {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	cells := make([][]Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			
			// Convert to 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			char, fg, bg := convertPixel(r8, g8, b8, a8, opts)
			cells[y][x] = Cell{
				Char: char,
				FG:   fg,
				BG:   bg,
			}
		}
	}

	return &Frame{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// convertPixel converts a pixel to ASCII character and colors
func convertPixel(r, g, b, a uint8, opts Options) (rune, string, string) {
	// Handle transparency
	if a < 128 {
		return ' ', "#000000", "#000000"
	}

	// Calculate luminosity
	luminosity := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)))

	var char rune
	
	// Select character based on method
	switch opts.Method {
	case "edge":
		char = getEdgeChar(luminosity)
	case "block":
		char = getBlockChar(luminosity)
	case "dither":
		char = getDitherChar(luminosity, r, g, b)
	case "luminosity":
		fallthrough
	default:
		if opts.Chars != "" {
			char = getCustomChar(luminosity, opts.Chars)
		} else {
			char = getLuminosityChar(luminosity)
		}
	}

	// Convert RGB to hex color
	fg := fmt.Sprintf("#%02X%02X%02X", r, g, b)
	bg := "#000000"

	return char, fg, bg
}

// getLuminosityChar returns character based on brightness
func getLuminosityChar(lum uint8) rune {
	// Extended ramp for better gradation
	extendedChars := []rune{
		' ', '·', '`', '.', ',', ':', ';', '-', '~', '=', '+', 
		'*', 'o', 'x', 'O', 'X', '#', '%', '&', '@', '█',
	}
	
	idx := int(lum) * len(extendedChars) / 256
	if idx >= len(extendedChars) {
		idx = len(extendedChars) - 1
	}
	
	return extendedChars[idx]
}

// getBlockChar returns block drawing character
func getBlockChar(lum uint8) rune {
	blocks := []rune{' ', '░', '▒', '▓', '█'}
	idx := int(lum) * len(blocks) / 256
	if idx >= len(blocks) {
		idx = len(blocks) - 1
	}
	return blocks[idx]
}

// getEdgeChar returns line drawing character (simplified)
func getEdgeChar(lum uint8) rune {
	if lum < 64 {
		return ' '
	} else if lum < 128 {
		return '·'
	} else if lum < 192 {
		return '─'
	} else {
		return '━'
	}
}

// getDitherChar returns dithered character
func getDitherChar(lum uint8, r, g, b uint8) rune {
	// Simple dithering using brightness
	patterns := []rune{' ', '·', ':', '░', '▒', '▓', '█'}
	idx := int(lum) * len(patterns) / 256
	if idx >= len(patterns) {
		idx = len(patterns) - 1
	}
	return patterns[idx]
}

// getCustomChar uses custom character set
func getCustomChar(lum uint8, chars string) rune {
	runes := []rune(chars)
	if len(runes) == 0 {
		return ' '
	}
	idx := int(lum) * len(runes) / 256
	if idx >= len(runes) {
		idx = len(runes) - 1
	}
	return runes[idx]
}

// SaveFrames saves frames to a file
func SaveFrames(frames []*Frame, filename string) error {
	// TODO: Implement file format
	// For now, save as JSON or custom format
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write metadata
	fmt.Fprintf(f, "aart v0.1.0\n")
	fmt.Fprintf(f, "frames: %d\n", len(frames))
	if len(frames) > 0 {
		fmt.Fprintf(f, "dimensions: %dx%d\n", frames[0].Width, frames[0].Height)
	}
	fmt.Fprintf(f, "---\n")

	// Write frames
	for i, frame := range frames {
		fmt.Fprintf(f, "frame: %d\n", i)
		fmt.Fprintf(f, "delay: %d\n", frame.Delay)
		for y := 0; y < frame.Height; y++ {
			for x := 0; x < frame.Width; x++ {
				cell := frame.Cells[y][x]
				fmt.Fprintf(f, "%c", cell.Char)
			}
			fmt.Fprintf(f, "\n")
		}
		fmt.Fprintf(f, "---\n")
	}

	fmt.Printf("Saved to: %s\n", filename)
	return nil
}
