package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mlamkadm/aart/internal/config"
	"github.com/mlamkadm/aart/internal/converter"
	"github.com/mlamkadm/aart/internal/fileformat"
	"github.com/mlamkadm/aart/internal/ui"
)

var (
	importGif    = flag.String("import-gif", "", "Import GIF file (URL or local path)")
	outputFile   = flag.String("output", "", "Output file (default: opens in editor)")
	width        = flag.Int("width", 0, "Canvas width for import (0 = auto-detect from terminal)")
	height       = flag.Int("height", 0, "Canvas height for import (0 = auto-detect from terminal)")
	fps          = flag.Int("fps", 12, "Target FPS for imported animation")
	method       = flag.String("method", "luminosity", "Conversion method: luminosity, average, block, dither")
	ratio        = flag.String("ratio", "fill", "Aspect ratio handling: fill, fit, original")
	chars        = flag.String("chars", "", "Custom character set for conversion (default: auto)")
	useColors    = flag.Bool("colors", false, "Use RGB colors (default: monochrome grayscale)")
	showHelp     = flag.Bool("help", false, "Show help message")
	version      = flag.Bool("version", false, "Show version")
	initConfig   = flag.Bool("init", false, "Initialize configuration directory")
	showConfig   = flag.Bool("show-config", false, "Show current configuration")
	configPath   = flag.Bool("config-path", false, "Show configuration file path")
	rawMode      = flag.Bool("raw", false, "Raw playback mode (no UI, just animation)")
	centerMode   = flag.Bool("center", false, "Center the animation in terminal (works with --raw)")
	onceMode     = flag.Bool("once", false, "Play animation once then exit (works with --raw)")
	
	// Export options
	exportFile   = flag.String("export", "", "Export file to format (specify output path)")
	exportFormat = flag.String("export-format", "json", "Export format: json, csv, ansi, txt, html, svg")
	exportFrame  = flag.Int("export-frame", -1, "Export specific frame (-1 for all)")
	exportColors = flag.Bool("export-colors", true, "Include colors in export")
)

const versionString = "aart v0.1.0"

func main() {
	flag.Parse()
	
	// Track which flags were explicitly set by the user
	flagsSet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})

	if *showHelp {
		printHelp()
		return
	}

	if *version {
		fmt.Println(versionString)
		return
	}

	// Handle config commands
	if *initConfig {
		if err := config.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}
		dir, _ := config.ConfigDir()
		fmt.Printf("✓ Configuration initialized at: %s\n", dir)
		fmt.Println("✓ Created directories: themes, templates, recent, backups")
		fmt.Println("✓ Created default config.yml")
		return
	}

	if *configPath {
		path, err := config.ConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(path)
		return
	}

	if *showConfig {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		printConfig(cfg)
		return
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config, using defaults: %v\n", err)
		cfg = &config.DefaultConfig
	}

	// Handle export
	if *exportFile != "" {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: specify input .aart file for export\n")
			fmt.Fprintf(os.Stderr, "Usage: aart <file.aart> --export output.json --export-format json\n")
			os.Exit(1)
		}
		if err := handleExport(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle GIF import
	if *importGif != "" {
		if err := handleGifImport(cfg, flagsSet); err != nil {
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
		// Load the specified file
		filepath := args[0]
		
		// Check if file exists and is .aa/.aart
		if _, err := os.Stat(filepath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", filepath)
			os.Exit(1)
		}
		
		// Load the .aa file
		aartFile, err := fileformat.Load(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading file: %v\n", err)
			os.Exit(1)
		}
		
		// Raw mode: just play the animation without UI
		if *rawMode {
			playRawAnimation(aartFile)
			return
		}
		
		// Create editor model with loaded file
		model = ui.NewWithFile(cfg, filepath, aartFile)
	} else {
		// Show startup page
		model = ui.NewStartupPage(cfg)
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

// playRawAnimation plays animation in raw mode (no UI, just frames)
func playRawAnimation(aartFile *fileformat.AartFile) {
	if len(aartFile.Frames) == 0 {
		return
	}
	
	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[H\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit
	
	// Calculate frame delay from first frame duration
	frameDuration := time.Duration(aartFile.Frames[0].Duration) * time.Millisecond
	if frameDuration == 0 {
		frameDuration = 83 * time.Millisecond // Default ~12fps
	}
	
	// Play animation in loop
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()
	
	frameIdx := 0
	
	// Get terminal size for centering if needed
	var termWidth, termHeight int
	if *centerMode {
		termWidth, termHeight = getTerminalSize()
	}
	
	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	loopCount := 0
	maxLoops := -1 // Infinite by default
	if *onceMode {
		maxLoops = 1
	}
	
	for maxLoops < 0 || loopCount < maxLoops {
		select {
		case <-sigChan:
			return
		case <-ticker.C:
			// Move cursor to home
			fmt.Print("\033[H")
			
			// Render current frame
			frame := aartFile.Frames[frameIdx]
			
			if *centerMode {
				// Calculate padding for centering
				frameHeight := len(frame.Cells)
				frameWidth := 0
				if frameHeight > 0 {
					frameWidth = len(frame.Cells[0])
				}
				
				verticalPadding := (termHeight - frameHeight) / 2
				horizontalPadding := (termWidth - frameWidth) / 2
				
				if verticalPadding < 0 {
					verticalPadding = 0
				}
				if horizontalPadding < 0 {
					horizontalPadding = 0
				}
				
				// Add vertical padding
				for i := 0; i < verticalPadding; i++ {
					fmt.Println()
				}
				
				// Render frame with horizontal padding
				for _, row := range frame.Cells {
					// Add horizontal padding
					for i := 0; i < horizontalPadding; i++ {
						fmt.Print(" ")
					}
					
					for _, cell := range row {
						char := cell.Char
						if char == "" {
							char = " "
						}
						
						// Apply colors if present
						if cell.Foreground != "" && cell.Foreground != "#000000" {
							fg := hexToRGB(cell.Foreground)
							// If bg is black (#000000), use transparent (don't set bg)
							if cell.Background != "" && cell.Background != "#000000" {
								bg := hexToRGB(cell.Background)
								fmt.Printf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%s\033[0m", 
									fg[0], fg[1], fg[2], bg[0], bg[1], bg[2], char)
							} else {
								// Foreground color only, transparent background
								fmt.Printf("\033[38;2;%d;%d;%dm%s\033[0m", 
									fg[0], fg[1], fg[2], char)
							}
						} else {
							fmt.Print(char)
						}
					}
					fmt.Println()
				}
			} else {
				// Normal rendering without centering
				for _, row := range frame.Cells {
					for _, cell := range row {
						char := cell.Char
						if char == "" {
							char = " "
						}
						
						// Apply colors if present
						if cell.Foreground != "" && cell.Foreground != "#000000" {
							fg := hexToRGB(cell.Foreground)
							// If bg is black (#000000), use transparent (don't set bg)
							if cell.Background != "" && cell.Background != "#000000" {
								bg := hexToRGB(cell.Background)
								fmt.Printf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%s\033[0m", 
									fg[0], fg[1], fg[2], bg[0], bg[1], bg[2], char)
							} else {
								// Foreground color only, transparent background
								fmt.Printf("\033[38;2;%d;%d;%dm%s\033[0m", 
									fg[0], fg[1], fg[2], char)
							}
						} else {
							fmt.Print(char)
						}
					}
					fmt.Println()
				}
			}
			
			// Next frame
			frameIdx = (frameIdx + 1) % len(aartFile.Frames)
			
			// Check if we completed a loop
			if frameIdx == 0 {
				loopCount++
			}
		}
	}
}

// hexToRGB converts hex color string to RGB values
func hexToRGB(hex string) [3]int {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	
	// Default to black if invalid
	if len(hex) != 6 {
		return [3]int{0, 0, 0}
	}
	
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return [3]int{r, g, b}
}

func handleGifImport(cfg *config.Config, flagsSet map[string]bool) error {
	// Auto-detect terminal size if width/height not specified
	convertWidth := *width
	convertHeight := *height
	
	if convertWidth == 0 || convertHeight == 0 {
		termWidth, termHeight := getTerminalSize()
		if convertWidth == 0 {
			// Use 80% of terminal width, accounting for UI chrome
			convertWidth = int(float64(termWidth) * 0.8)
			if convertWidth < 40 {
				convertWidth = cfg.Editor.DefaultWidth
			}
		}
		if convertHeight == 0 {
			// Use 80% of terminal height, accounting for UI chrome
			convertHeight = int(float64(termHeight-10) * 0.8)
			if convertHeight < 20 {
				convertHeight = cfg.Editor.DefaultHeight
			}
		}
	}
	
	convertFPS := *fps
	convertMethod := *method
	convertRatio := *ratio
	
	// Use config defaults ONLY if flags weren't explicitly set
	if !flagsSet["fps"] && cfg.Editor.DefaultFPS != 0 {
		convertFPS = cfg.Editor.DefaultFPS
	}
	if !flagsSet["method"] && cfg.Converter.DefaultMethod != "" {
		convertMethod = cfg.Converter.DefaultMethod
	}
	
	fmt.Printf("🎨 aart - GIF to ASCII Converter\n\n")
	fmt.Printf("Source: %s\n", *importGif)
	fmt.Printf("Target: %dx%d @ %dfps\n", convertWidth, convertHeight, convertFPS)
	fmt.Printf("Method: %s\n", convertMethod)
	fmt.Printf("Ratio: %s\n", convertRatio)
	fmt.Printf("Colors: %v\n\n", *useColors)

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
		Width:            convertWidth,
		Height:           convertHeight,
		FPS:              convertFPS,
		Method:           convertMethod,
		Ratio:            convertRatio,
		Chars:            *chars,
		UseColors:        *useColors,
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
		
		// Add to recent files
		cfg.AddRecentFile(*outputFile, len(frames))
		config.Save(cfg)
		
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
		ui.NewWithFramesAndConfig(uiFrames, cfg),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

// getTerminalSize returns the current terminal dimensions
func getTerminalSize() (width, height int) {
	// Try to get terminal size using TIOCGWINSZ ioctl
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	
	// This is a simplified version - in production you'd use golang.org/x/term
	// For now, use environment variables or defaults
	width = 120
	height = 40
	
	// Try to get from stty command
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err == nil {
		fmt.Sscanf(string(out), "%d %d", &height, &width)
	}
	
	// Fallback to reasonable defaults if detection failed
	if width == 0 {
		width = 120
	}
	if height == 0 {
		height = 40
	}
	
	return width, height
}

func printConfig(cfg *config.Config) {
	fmt.Printf("📝 aart Configuration\n\n")
	
	fmt.Printf("Version: %s\n\n", cfg.Version)
	
	fmt.Printf("Editor Settings:\n")
	fmt.Printf("  Default Size: %dx%d\n", cfg.Editor.DefaultWidth, cfg.Editor.DefaultHeight)
	fmt.Printf("  Default FPS: %d\n", cfg.Editor.DefaultFPS)
	fmt.Printf("  Auto Save: %v (interval: %ds)\n", cfg.Editor.AutoSave, cfg.Editor.AutoSaveInterval)
	fmt.Printf("  Show Grid: %v\n", cfg.Editor.ShowGrid)
	fmt.Printf("  Zen Mode: %v\n\n", cfg.Editor.ZenMode)
	
	fmt.Printf("UI Settings:\n")
	fmt.Printf("  Theme: %s\n", cfg.UI.Theme)
	fmt.Printf("  Cursor Style: %s\n", cfg.UI.CursorStyle)
	fmt.Printf("  Progress Style: %s\n", cfg.UI.ProgressStyle)
	fmt.Printf("  Animation Smooth: %v\n\n", cfg.UI.AnimationSmooth)
	
	fmt.Printf("Colors:\n")
	fmt.Printf("  Scheme: %s\n", cfg.Colors.Name)
	fmt.Printf("  Foreground: %s\n", cfg.Colors.Foreground)
	fmt.Printf("  Background: %s\n", cfg.Colors.Background)
	fmt.Printf("  Cursor: %s\n\n", cfg.Colors.Cursor)
	
	fmt.Printf("Converter:\n")
	fmt.Printf("  Default Method: %s\n", cfg.Converter.DefaultMethod)
	fmt.Printf("  Quality: %s\n", cfg.Converter.Quality)
	fmt.Printf("  Preserve Aspect: %v\n\n", cfg.Converter.PreserveAspect)
	
	fmt.Printf("Recent Files: (%d)\n", len(cfg.Recent.Files))
	for i, rf := range cfg.Recent.Files {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(cfg.Recent.Files)-5)
			break
		}
		fmt.Printf("  %s (%d frames) - %s\n", rf.Path, rf.Frames, rf.Timestamp.Format("2006-01-02 15:04"))
	}
	
	path, _ := config.ConfigPath()
	fmt.Printf("\nConfig File: %s\n", path)
}

func handleExport(inputFile string) error {
	fmt.Printf("📦 Exporting: %s\n", inputFile)
	fmt.Printf("Format: %s\n", *exportFormat)
	fmt.Printf("Output: %s\n\n", *exportFile)

	// Load .aart file
	aart, err := fileformat.Load(inputFile)
	if err != nil {
		return fmt.Errorf("failed to load .aart file: %w", err)
	}

	// Validate
	if err := aart.Validate(); err != nil {
		return fmt.Errorf("invalid .aart file: %w", err)
	}

	// Prepare export options
	opts := fileformat.ExportOptions{
		Format:      fileformat.ExportFormat(*exportFormat),
		FrameIndex:  *exportFrame,
		IncludeMeta: true,
		Colors:      *exportColors,
	}

	// Export
	if err := fileformat.Export(aart, *exportFile, opts); err != nil {
		return err
	}

	fmt.Printf("✓ Exported successfully!\n")
	fmt.Printf("  Frames: %d\n", aart.FrameCount())
	fmt.Printf("  Size: %dx%d\n", aart.Canvas.Width, aart.Canvas.Height)
	if *exportFrame >= 0 {
		fmt.Printf("  Frame: %d only\n", *exportFrame)
	}

	return nil
}

func printHelp() {
	fmt.Printf(`%s - ASCII Art Animation Editor

USAGE:
    aart [options] [file]
    aart --import-gif <url|path> [options]
    aart --init                    # Initialize configuration
    aart --show-config             # Show current configuration

OPTIONS:
    --import-gif <source>    Import GIF from URL or local path
    --output <file>          Save imported frames to file (default: open editor)
    --width <int>            Canvas width (default: auto from terminal)
    --height <int>           Canvas height (default: auto from terminal)
    --fps <int>              Target FPS (default: 12)
    --method <string>        Conversion method (default: luminosity)
                             Options: luminosity, edge, block, dither
    --ratio <string>         Aspect ratio (default: fill)
                             Options: fill, fit, original
    --chars <string>         Custom character set for conversion
    --colors                 Use RGB colors (default: monochrome grayscale)
    
CONFIGURATION:
    --init                   Initialize ~/.config/aart directory
    --show-config            Display current configuration
    --config-path            Show configuration file path
    
INFO:
    --help                   Show this help message
    --version                Show version

CONVERSION METHODS:
    luminosity    Convert based on brightness (default)
    edge          Edge detection with line characters
    block         Block characters (░▒▓█)
    dither        Dithered output

CONFIGURATION FILE:
    Location: ~/.config/aart/config.yml
    
    Settings include:
    - Editor defaults (width, height, FPS, auto-save)
    - UI preferences (theme, cursor style, progress style)
    - Color schemes (foreground, background, cursor colors)
    - Converter defaults (method, quality, aspect ratio)
    - Recent files tracking
    - Custom keybindings

EXAMPLES:
    # Initialize configuration
    aart --init

    # Start editor with default settings
    aart

    # Import GIF from URL
    aart --import-gif https://example.com/animation.gif

    # Import local GIF with custom size
    aart --import-gif ./animation.gif --width 120 --height 40

    # Import and save to file
    aart --import-gif animation.gif --output animation.aart

    # Import with specific method
    aart --import-gif animation.gif --method block

    # Show current configuration
    aart --show-config

    # Open existing file
    aart animation.aart

KEYBOARD SHORTCUTS (in editor):
    hjkl / arrows    Navigate cursor
    i                Insert mode
    d                Draw character
    space            Play/pause
    , .              Seek frames
    z                Toggle zen mode
    ctrl-j/k         Cycle wheel menu
    :                Command mode
    q                Quit

CONFIGURATION DIRECTORIES:
    ~/.config/aart/           # Main configuration directory
    ~/.config/aart/themes/    # Custom color themes
    ~/.config/aart/templates/ # Animation templates
    ~/.config/aart/recent/    # Recent files cache
    ~/.config/aart/backups/   # Auto-save backups

For more information, see README.md or visit:
https://github.com/mlamkadm/aart
`, versionString)
}
