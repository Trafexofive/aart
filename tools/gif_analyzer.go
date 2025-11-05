package main

import (
	"fmt"
	"image/gif"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gif_analyzer <gif_file_or_url>")
		os.Exit(1)
	}

	source := os.Args[1]
	
	// Load GIF
	var file *os.File
	var err error
	
	if len(source) > 4 && source[:4] == "http" {
		fmt.Println("URL import not implemented in analyzer yet")
		os.Exit(1)
	}
	
	file, err = os.Open(source)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	gifData, err := gif.DecodeAll(file)
	if err != nil {
		fmt.Printf("Error decoding GIF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("GIF Analysis\n")
	fmt.Printf("============\n\n")
	fmt.Printf("Total Frames: %d\n", len(gifData.Image))
	fmt.Printf("Loop Count: %d (0 = infinite)\n", gifData.LoopCount)
	fmt.Printf("\nFrame Details:\n")
	
	for i := 0; i < min(10, len(gifData.Image)); i++ {
		img := gifData.Image[i]
		delay := gifData.Delay[i]
		disposal := gifData.Disposal[i]
		
		bounds := img.Bounds()
		
		fmt.Printf("\nFrame %d:\n", i)
		fmt.Printf("  Size: %dx%d\n", bounds.Dx(), bounds.Dy())
		fmt.Printf("  Delay: %d (%.1fms)\n", delay, float64(delay*10))
		fmt.Printf("  Disposal: %d ", disposal)
		
		switch disposal {
		case 0:
			fmt.Printf("(no disposal)\n")
		case 1:
			fmt.Printf("(do not dispose)\n")
		case 2:
			fmt.Printf("(restore to background)\n")
		case 3:
			fmt.Printf("(restore to previous)\n")
		}
		
		fmt.Printf("  Palette colors: %d\n", len(img.Palette))
		
		// Check for transparency
		hasTransparent := false
		for _, c := range img.Palette {
			_, _, _, a := c.RGBA()
			if a == 0 {
				hasTransparent = true
				break
			}
		}
		fmt.Printf("  Has transparency: %v\n", hasTransparent)
		
		// Sample some pixels
		if i < 3 {
			fmt.Printf("  Sample pixels (first row):\n")
			for x := 0; x < min(10, bounds.Dx()); x++ {
				idx := img.ColorIndexAt(x, 0)
				r, g, b, a := img.Palette[idx].RGBA()
				fmt.Printf("    [%d,%d]: idx=%d rgba=(%d,%d,%d,%d)\n", 
					x, 0, idx, r>>8, g>>8, b>>8, a>>8)
			}
		}
	}
	
	if len(gifData.Image) > 10 {
		fmt.Printf("\n... and %d more frames\n", len(gifData.Image)-10)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
