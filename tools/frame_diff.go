package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type AartFile struct {
	Version string `json:"version"`
	Canvas  struct {
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
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.aart> [frame1] [frame2]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Defaults: frame1=0, frame2=1\n")
		os.Exit(1)
	}

	filename := os.Args[1]
	frame1 := 0
	frame2 := 1

	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &frame1)
	}
	if len(os.Args) > 3 {
		fmt.Sscanf(os.Args[3], "%d", &frame2)
	}

	// Load file
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var aart AartFile
	if err := json.Unmarshal(data, &aart); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if frame1 >= len(aart.Frames) || frame2 >= len(aart.Frames) {
		fmt.Fprintf(os.Stderr, "Error: file only has %d frames\n", len(aart.Frames))
		os.Exit(1)
	}

	fmt.Printf("Comparing Frame %d vs Frame %d\n", frame1, frame2)
	fmt.Printf("Canvas: %dx%d\n\n", aart.Canvas.Width, aart.Canvas.Height)

	f1 := aart.Frames[frame1]
	f2 := aart.Frames[frame2]

	// Count differences
	totalCells := 0
	diffCells := 0
	charDiff := 0
	colorDiff := 0

	for y := 0; y < len(f1.Cells) && y < len(f2.Cells); y++ {
		for x := 0; x < len(f1.Cells[y]) && x < len(f2.Cells[y]); x++ {
			totalCells++
			c1 := f1.Cells[y][x]
			c2 := f2.Cells[y][x]

			isDiff := false
			if c1.Char != c2.Char {
				charDiff++
				isDiff = true
			}
			if c1.Foreground != c2.Foreground || c1.Background != c2.Background {
				colorDiff++
				isDiff = true
			}
			if isDiff {
				diffCells++
			}
		}
	}

	fmt.Printf("Statistics:\n")
	fmt.Printf("  Total cells:      %d\n", totalCells)
	fmt.Printf("  Different cells:  %d (%.1f%%)\n", diffCells, float64(diffCells)*100/float64(totalCells))
	fmt.Printf("  Character diffs:  %d (%.1f%%)\n", charDiff, float64(charDiff)*100/float64(totalCells))
	fmt.Printf("  Color diffs:      %d (%.1f%%)\n\n", colorDiff, float64(colorDiff)*100/float64(totalCells))

	// Show side-by-side comparison (first 40 rows)
	showRows := 40
	if len(f1.Cells) < showRows {
		showRows = len(f1.Cells)
	}

	fmt.Printf("Side-by-side comparison (first %d rows):\n", showRows)
	fmt.Printf("Frame %d                                    | Frame %d\n", frame1, frame2)
	fmt.Println(strings.Repeat("─", 45) + "┼" + strings.Repeat("─", 45))

	for y := 0; y < showRows; y++ {
		// Frame 1
		if y < len(f1.Cells) {
			line := ""
			for x := 0; x < len(f1.Cells[y]) && x < 43; x++ {
				line += f1.Cells[y][x].Char
			}
			fmt.Printf("%-45s", line)
		} else {
			fmt.Printf("%-45s", "")
		}
		
		fmt.Print("│ ")

		// Frame 2
		if y < len(f2.Cells) {
			line := ""
			for x := 0; x < len(f2.Cells[y]) && x < 43; x++ {
				line += f2.Cells[y][x].Char
			}
			fmt.Printf("%s\n", line)
		} else {
			fmt.Println()
		}
	}

	// Show diff visualization
	fmt.Printf("\nDifference map (X = different, · = same):\n")
	for y := 0; y < showRows; y++ {
		if y >= len(f1.Cells) || y >= len(f2.Cells) {
			break
		}
		for x := 0; x < len(f1.Cells[y]) && x < len(f2.Cells[y]) && x < 80; x++ {
			c1 := f1.Cells[y][x]
			c2 := f2.Cells[y][x]
			if c1.Char != c2.Char || c1.Foreground != c2.Foreground {
				fmt.Print("X")
			} else {
				fmt.Print("·")
			}
		}
		fmt.Println()
	}
}
