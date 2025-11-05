package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mlamkadm/aart/internal/config"
	"github.com/mlamkadm/aart/internal/converter"
	"github.com/mlamkadm/aart/internal/fileformat"
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
	initConfig   = flag.Bool("init", false, "Initialize configuration directory")
	showConfig   = flag.Bool("show-config", false, "Show current configuration")
	configPath   = flag.Bool("config-path", false, "Show configuration file path")
	
	// Export options
	exportFile   = flag.String("export", "", "Export file to format (specify output path)")
	exportFormat = flag.String("export-format", "json", "Export format: json, csv, ansi, txt, html, svg")
	exportFrame  = flag.Int("export-frame", -1, "Export specific frame (-1 for all)")
	exportColors = flag.Bool("export-colors", true, "Include colors in export")
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
		if err := handleGifImport(cfg); err != nil {
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
		model = ui.NewWithConfig(cfg)
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

func handleGifImport(cfg *config.Config) error {
	fmt.Printf("🎨 aart - GIF to ASCII Converter\n\n")
	fmt.Printf("Source: %s\n", *importGif)
	fmt.Printf("Target: %dx%d @ %dfps\n", *width, *height, *fps)
	fmt.Printf("Method: %s\n\n", *method)

	// Use config defaults if not specified
	convertWidth := *width
	convertHeight := *height
	convertFPS := *fps
	convertMethod := *method
	
	if convertWidth == 80 {
		convertWidth = cfg.Editor.DefaultWidth
	}
	if convertHeight == 24 {
		convertHeight = cfg.Editor.DefaultHeight
	}
	if convertFPS == 12 {
		convertFPS = cfg.Editor.DefaultFPS
	}
	if convertMethod == "luminosity" && cfg.Converter.DefaultMethod != "" {
		convertMethod = cfg.Converter.DefaultMethod
	}

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
    --width <int>            Canvas width (default: from config or 80)
    --height <int>           Canvas height (default: from config or 24)
    --fps <int>              Target FPS (default: from config or 12)
    --method <string>        Conversion method (default: from config or luminosity)
                             Options: luminosity, edge, block, dither
    --chars <string>         Custom character set for conversion
    
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
