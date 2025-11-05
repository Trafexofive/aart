package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mlamkadm/aart/internal/converter"
	"github.com/mlamkadm/aart/internal/ui"
)

var (
	importGif    = flag.String("import-gif", "", "Import GIF file (URL or local path)")
	outputFile   = flag.String("output", "", "Output file (default: opens in editor)")
	width        = flag.Int("width", 80, "Canvas width for import")
	height       = flag.Int("height", 24, "Canvas height for import")
	fps          = flag.Int("fps", 12, "Target FPS for imported animation")
	method       = flag.String("method", "luminosity", "Conversion method: luminosity, edge, block, dither")
	chars        = flag.String("chars", "", "Custom character set for conversion (default: auto)")
	showHelp     = flag.Bool("help", false, "Show help message")
	version      = flag.Bool("version", false, "Show version")
)

const versionString = "aart v0.1.0"

func main() {
	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	if *version {
		fmt.Println(versionString)
		return
	}

	// Handle GIF import
	if *importGif != "" {
		if err := handleGifImport(); err != nil {
			fmt.Fprintf(os.Stderr, "Error importing GIF: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Start interactive editor
	var model tea.Model
	
	// Check if a file was passed as argument
	args := flag.Args()
	if len(args) > 0 {
		// TODO: Load file
		model = ui.New()
	} else {
		model = ui.New()
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleGifImport() error {
	fmt.Printf("🎨 aart - GIF to ASCII Converter\n\n")
	fmt.Printf("Source: %s\n", *importGif)
	fmt.Printf("Target: %dx%d @ %dfps\n", *width, *height, *fps)
	fmt.Printf("Method: %s\n\n", *method)

	// Progress tracking
	lastPercent := 0
	progressCallback := func(current, total int, message string) {
		percent := current
		if percent != lastPercent {
			// Simple progress bar
			barWidth := 40
			filled := percent * barWidth / 100
			bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
			fmt.Printf("\r%s %3d%% %s", bar, percent, message)
			if percent >= 100 {
				fmt.Println()
			}
			lastPercent = percent
		}
	}

	frames, err := converter.ConvertGifToFrames(*importGif, converter.Options{
		Width:            *width,
		Height:           *height,
		FPS:              *fps,
		Method:           *method,
		Chars:            *chars,
		ProgressCallback: progressCallback,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Converted %d frames successfully!\n\n", len(frames))

	// If output file specified, save it
	if *outputFile != "" {
		fmt.Printf("💾 Saving to %s...\n", *outputFile)
		if err := converter.SaveFrames(frames, *outputFile); err != nil {
			return err
		}
		fmt.Printf("✓ Saved!\n")
		return nil
	}

	// Convert to UI format
	uiFrames := make([]ui.ImportedFrame, len(frames))
	for i, f := range frames {
		cells := make([][]ui.ImportedCell, f.Height)
		for y := 0; y < f.Height; y++ {
			cells[y] = make([]ui.ImportedCell, f.Width)
			for x := 0; x < f.Width; x++ {
				cells[y][x] = ui.ImportedCell{
					Char: f.Cells[y][x].Char,
					FG:   f.Cells[y][x].FG,
					BG:   f.Cells[y][x].BG,
				}
			}
		}
		uiFrames[i] = ui.ImportedFrame{
			Width:  f.Width,
			Height: f.Height,
			Cells:  cells,
			Delay:  f.Delay,
		}
	}

	// Otherwise, open in editor
	fmt.Println("🖼️  Opening in editor...\n")
	p := tea.NewProgram(
		ui.NewWithFrames(uiFrames),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func printHelp() {
	fmt.Printf(`%s - ASCII Art Animation Editor

USAGE:
    aart [options] [file]
    aart --import-gif <url|path> [options]

OPTIONS:
    --import-gif <source>    Import GIF from URL or local path
    --output <file>          Save imported frames to file (default: open editor)
    --width <int>            Canvas width (default: 80)
    --height <int>           Canvas height (default: 24)
    --fps <int>              Target FPS (default: 12)
    --method <string>        Conversion method (default: luminosity)
                             Options: luminosity, edge, block, dither
    --chars <string>         Custom character set for conversion
    --help                   Show this help message
    --version                Show version

CONVERSION METHODS:
    luminosity    Convert based on brightness (default)
    edge          Edge detection with line characters
    block         Block characters (░▒▓█)
    dither        Dithered output

EXAMPLES:
    # Start editor
    aart

    # Import GIF from URL
    aart --import-gif https://example.com/animation.gif

    # Import local GIF with custom size
    aart --import-gif ./animation.gif --width 120 --height 40

    # Import and save to file
    aart --import-gif animation.gif --output animation.aart

    # Import with edge detection
    aart --import-gif animation.gif --method edge

    # Open existing file
    aart animation.aart

KEYBOARD SHORTCUTS (in editor):
    hjkl / arrows    Navigate cursor
    i                Insert mode
    d                Draw character
    space            Play/pause
    , .              Seek frames
    ctrl-j/k         Cycle wheel menu
    :                Command mode
    q                Quit

For more information, see README.md
`, versionString)
}
