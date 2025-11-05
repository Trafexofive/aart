package converter

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

type Options struct {
	Width            int
	Height           int
	FPS              int
	Method           string
	Ratio            string // "fill", "fit", "original"
	Chars            string
	UseColors        bool   // If true, include RGB colors; if false, monochrome
	ProgressCallback func(current, total int, message string)
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

// ImportGIF is a convenience wrapper for UI that returns frames directly
func ImportGIF(source string, width, height, fps int, method, ratio string) ([]*Frame, error) {
	return ConvertGifToFrames(source, Options{
		Width:  width,
		Height: height,
		FPS:    fps,
		Method: method,
		Ratio:  ratio,
	})
}

// ConvertGifToFrames converts a GIF (from URL or file) to ASCII frames
func ConvertGifToFrames(source string, opts Options) ([]*Frame, error) {
	// Report progress
	if opts.ProgressCallback != nil {
		opts.ProgressCallback(0, 100, "Loading GIF...")
	}
	
	// Load GIF
	gifData, err := loadGif(source)
	if err != nil {
		return nil, fmt.Errorf("failed to load GIF: %w", err)
	}

	if opts.ProgressCallback != nil {
		opts.ProgressCallback(10, 100, fmt.Sprintf("Processing %d frames...", len(gifData.Image)))
	}

	// Convert each frame with proper composition
	frames := make([]*Frame, len(gifData.Image))
	var previousFrame image.Image
	
	for i, img := range gifData.Image {
		// Report progress every 10 frames or for small GIFs
		if opts.ProgressCallback != nil && (i%10 == 0 || len(gifData.Image) < 50) {
			percent := 10 + (i * 80 / len(gifData.Image))
			opts.ProgressCallback(percent, 100, fmt.Sprintf("Converting frame %d/%d...", i+1, len(gifData.Image)))
		}
		
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

	if opts.ProgressCallback != nil {
		opts.ProgressCallback(100, 100, "Complete!")
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
		return ' ', "#FFFFFF", "#000000"
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

	// Return colors based on mode
	var fg, bg string
	if opts.UseColors {
		// Full RGB color mode
		fg = fmt.Sprintf("#%02X%02X%02X", r, g, b)
		bg = "#000000"
	} else {
		// Monochrome mode - quantize to 16 gray levels for consistency
		// This reduces color variations and makes frames more similar
		grayLevel := int(luminosity) / 16 // 0-15
		quantized := uint8(grayLevel * 17) // Map back to 0-255 in steps of 17
		fg = fmt.Sprintf("#%02X%02X%02X", quantized, quantized, quantized)
		bg = "#000000"
	}

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

// SaveFrames saves frames to .aart format
func SaveFrames(frames []*Frame, filename string) error {
	if len(frames) == 0 {
		return fmt.Errorf("no frames to save")
	}

	// Create .aart file structure
	aartFile := &struct {
		Version  string `json:"version"`
		Metadata struct {
			Title    string `json:"title"`
			Created  string `json:"created"`
			Modified string `json:"modified"`
			Source   string `json:"source"`
		} `json:"metadata"`
		Canvas struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"canvas"`
		Frames []struct {
			Index    int `json:"index"`
			Duration int `json:"duration"`
			Cells    [][]struct {
				Char       string `json:"char"`
				Foreground string `json:"fg"`
				Background string `json:"bg"`
			} `json:"cells"`
		} `json:"frames"`
	}{}

	aartFile.Version = "1.0"
	aartFile.Metadata.Title = filename
	aartFile.Metadata.Created = time.Now().Format(time.RFC3339)
	aartFile.Metadata.Modified = time.Now().Format(time.RFC3339)
	aartFile.Metadata.Source = "converted"
	
	aartFile.Canvas.Width = frames[0].Width
	aartFile.Canvas.Height = frames[0].Height
	
	// Convert frames
	aartFile.Frames = make([]struct {
		Index    int `json:"index"`
		Duration int `json:"duration"`
		Cells    [][]struct {
			Char       string `json:"char"`
			Foreground string `json:"fg"`
			Background string `json:"bg"`
		} `json:"cells"`
	}, len(frames))

	for i, frame := range frames {
		aartFile.Frames[i].Index = i
		aartFile.Frames[i].Duration = frame.Delay
		aartFile.Frames[i].Cells = make([][]struct {
			Char       string `json:"char"`
			Foreground string `json:"fg"`
			Background string `json:"bg"`
		}, frame.Height)

		for y := 0; y < frame.Height; y++ {
			aartFile.Frames[i].Cells[y] = make([]struct {
				Char       string `json:"char"`
				Foreground string `json:"fg"`
				Background string `json:"bg"`
			}, frame.Width)
			
			for x := 0; x < frame.Width; x++ {
				cell := frame.Cells[y][x]
				aartFile.Frames[i].Cells[y][x].Char = string(cell.Char)
				aartFile.Frames[i].Cells[y][x].Foreground = cell.FG
				aartFile.Frames[i].Cells[y][x].Background = cell.BG
			}
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(aartFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal frames: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Saved to: %s\n", filename)
	return nil
}
