# aart - ASCII Art Animation Editor

A terminal-based ASCII art animation editor with GIF import capabilities built with Go and Bubbletea.

## Features

### Core Editor ✅
- **Interactive Canvas**: Full-featured drawing canvas with character-level editing
- **Timeline**: Multi-frame animation timeline with visual indicators
- **Playback**: Smooth animation playback with configurable FPS
- **Navigation**: Vim-style hjkl navigation with arrow key support
- **Drawing Tools**: 
  - Pencil mode for direct character placement
  - Insert mode for text entry
  - Multiple tool selections (fill, select, line, box, text, eyedropper, move)
- **Frame Management**: Navigate, create, and edit animation frames
- **View Controls**: 
  - Grid overlay toggle
  - Zoom in/out with reset
  - Zen mode for distraction-free editing
- **Radial Wheel UI**: Intuitive ctrl-j/k navigation through tool sections
- **Command Mode**: Vim-style command interface for advanced operations
- **Status Bar**: Real-time display of file info, frame, tool, FPS, and layer data

### GIF Import ✅
- **URL Support**: Direct import from web URLs
- **Local File Import**: Load GIFs from filesystem
- **Multiple Conversion Methods**:
  - **luminosity**: Brightness-based conversion (default)
  - **block**: Block character style (░▒▓█)
  - **edge**: Edge detection wireframe
  - **dither**: Dithered gradients
- **Smart Sizing**: Auto-detect terminal dimensions or specify custom size
- **FPS Control**: Match original or set custom frame rate
- **Aspect Ratio Modes**: Fill, original, or custom ratios

### Configuration System ✅
- **Config File**: `~/.config/aart/config.yml` for persistent settings
- **Project Storage**: `~/.config/aart/` directory for recent files and cache
- **Custom File Format**: `.aa` format for native ASCII art storage
- **Export Formats**: CSV, JSON, ANSI, plain text

### UI Features ✅
- **Startup Page**: Quick access to common actions
- **Recent Files**: Navigate and open recent projects
- **Theme System**: Configurable color schemes
- **Loader Animations**: Zen-like loading indicators
- **Help System**: Comprehensive keyboard shortcut guide

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/aart.git
cd aart

# Install dependencies
go mod download

# Build
make build

# Run
./aart
```

### Quick Start

```bash
# Start with blank canvas
./aart

# Import a GIF from URL
./aart --import-gif https://media.giphy.com/media/example.gif

# Import with custom settings
./aart --import-gif animation.gif --fps 60 --width 200 --height 60 --method dither

# Auto-size to terminal
./aart --import-gif animation.gif --ratio fill
```

## Keyboard Shortcuts

### Normal Mode
- `hjkl` / arrows - Move cursor
- `d` - Draw current character
- `i` - Insert mode
- `space` - Play/pause animation
- `,` `.` - Previous/next frame
- `+` `-` - Zoom in/out
- `0` - Reset zoom
- `g` - Toggle grid
- `z` - Toggle zen mode
- `p` `f` `s` `l` `b` `t` `e` `m` - Select tool
- `ctrl-j` `ctrl-k` - Cycle wheel sections
- `enter` - Expand wheel
- `esc` - Collapse wheel / Return to normal mode
- `:` - Command mode
- `?` - Show help
- `q` - Quit

### Insert Mode
- Type normally - Characters appear at cursor
- Arrow keys - Navigate while typing
- `esc` - Return to normal mode

### Command Mode
Commands available:
- `:export <file>` - Export animation
- `:import <file>` - Import file
- `:new` - New animation
- `:save <file>` - Save to .aa format
- `:quit` or `:q` - Exit
- `:help` - Show help

## Configuration

aart uses a YAML configuration file located at `~/.config/aart/config.yml`.

### Initialize Configuration

```bash
# Create config directory and default config.yml
./aart --init

# View current config path
./aart --config-path

# View current configuration
./aart --show-config
```

### Config File Structure

```yaml
version: 0.1.0

editor:
  default_width: 80          # Default canvas width
  default_height: 24         # Default canvas height
  default_fps: 12            # Default animation FPS
  auto_save: false           # Auto-save enabled
  auto_save_interval: 300    # Auto-save interval (seconds)
  tab_size: 4                # Tab width
  show_grid: false           # Show grid by default
  show_line_numbers: false   # Show line numbers
  zen_mode: false            # Start in zen mode

ui:
  theme: tokyo-night         # Theme name (tokyo-night, gruvbox, monokai, dracula)
  show_status_bar: true      # Show status bar
  show_timeline: true        # Show timeline
  show_wheel_by_default: false  # Show radial wheel on startup
  cursor_style: line         # Cursor style (block, line, underline)
  animation_smooth: true     # Smooth animations
  progress_style: bar        # Progress indicator (bar, spinner, minimal)
  border_style: rounded      # Border style (rounded, thick, double, ascii)
  timeline_style: detailed   # Timeline style (compact, detailed, minimal)
  status_bar_position: bottom  # Status bar position (top, bottom)

colors:
  name: default              # Color scheme name
  foreground: '#FFFFFF'      # Primary text color
  background: '#000000'      # Background color
  cursor: '#FFFF00'          # Cursor color
  selection: '#444444'       # Selection highlight
  status_bar: '#333333'      # Status bar background
  timeline: '#222222'        # Timeline background
  border: '#666666'          # Border color

recent:
  files: []                  # Auto-populated recent files
  max_entries: 10            # Maximum recent files to track

converter:
  default_method: luminosity # Default GIF conversion (luminosity, block, edge, dither)
  default_chars: ""          # Custom character ramp (empty = auto)
  preserve_aspect: true      # Preserve aspect ratio
  quality: high              # Conversion quality (low, medium, high)

startup:
  show_startup_page: true    # Show startup page on launch
  artwork_file: ""           # Custom ASCII art file path
  artwork_inline: ""         # Inline ASCII art (multiline YAML string)
  artwork_border: true       # Show border around artwork
  artwork_offset_x: 0        # X offset for artwork positioning
  artwork_offset_y: 0        # Y offset for artwork positioning
  artwork_width: 0           # Max width (0 = auto)
  artwork_height: 0          # Max height (0 = auto)
  show_recent_files: true    # Show recent files panel
  show_tips: true            # Show rotating tips
  tip_rotation_seconds: 5    # Seconds between tip rotation
  breathing_effect: true     # Enable breathing animation
```

### Custom Startup Artwork

You can customize the startup logo by either providing a file path or inline content:

```yaml
# Option 1: External file
startup:
  artwork_file: "~/.config/aart/my_logo.txt"
  artwork_border: true
  
# Option 2: Inline YAML
startup:
  artwork_inline: |
    ╔════════════════╗
    ║   MY EDITOR   ║
    ╚════════════════╝
  artwork_border: false
  
# Option 3: Relative to config directory
startup:
  artwork_file: "startup_art.txt"  # Looks in ~/.config/aart/startup_art.txt
  artwork_border: true
  artwork_offset_x: 2
  artwork_offset_y: 1
```

### Edit Config

```bash
# Edit config in $EDITOR
./aart --init  # Creates config if it doesn't exist
# Then select "Edit Config" from startup menu, or:
$EDITOR ~/.config/aart/config.yml

# Quick command from startup page
# Press 'c' key to open config in $EDITOR
```

## File Formats

### Native Format (.aa)

aart's native format stores complete animation data including layers, metadata, and frame information.

```bash
# Save
:save myanimation.aa

# Export to various formats
:export output.ans     # ANSI art
:export output.txt     # Plain text
:export output.json    # JSON metadata + frames
:export output.csv     # CSV format
```

### Import/Export

- **Import**: GIF, ANSI (.ans), plain text (.txt)
- **Export**: ANSI (.ans), plain text (.txt), JSON, CSV, native (.aa)

## Architecture

```
cmd/aart/           - Main entry point and CLI argument parsing
internal/
  ├── ui/           - Bubbletea UI components
  │   ├── model.go  - Main UI model and state
  │   ├── startup.go- Startup page
  │   ├── timeline.go- Timeline rendering
  │   └── wheel.go  - Radial wheel navigation
  ├── canvas/       - Canvas and frame management
  ├── config/       - Configuration system
  ├── converter/    - GIF to ASCII conversion
  └── export/       - Export format handlers
pkg/
  └── asciiart/     - Core ASCII art utilities
```

## Development Roadmap

### Phase 1: Core Editor ✅
- [x] Basic canvas with cursor navigation
- [x] Timeline and playback system
- [x] Radial wheel UI
- [x] Drawing tools foundation
- [x] Command mode
- [x] Status bar and help system

### Phase 2: GIF Import ✅
- [x] URL and file support
- [x] Multiple conversion methods
- [x] Terminal size auto-detection
- [x] FPS control
- [x] Aspect ratio modes

### Phase 3: Configuration ✅
- [x] Config file system (~/.config/aart/)
- [x] Recent files management
- [x] Theme system
- [x] Startup page with customization
- [x] Custom startup artwork support
- [x] Terminal size auto-detection
- [x] Aspect ratio modes (fill, fit, original)
- [x] Edit config from startup menu

### Phase 4: Current Focus 🚧
- [x] GIF import with URL and local file support
- [x] Multiple conversion methods (luminosity, block, edge, dither)
- [x] Auto-sizing to terminal dimensions
- [ ] Layer system implementation
- [ ] Advanced tool implementations (fill, line, box)
- [ ] Undo/redo stack
- [ ] Color picker improvements
- [ ] GIF conversion quality improvements

### Phase 5: Polish
- [ ] Mouse support
- [ ] Selection and clipboard operations
- [ ] Grid overlay rendering
- [ ] Minimap for large canvases
- [ ] Plugin system
- [ ] Tutorial/onboarding

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Areas for Contribution

- **Conversion Methods**: New algorithms for GIF to ASCII conversion
- **Tools**: Advanced drawing tools (bezier curves, spray paint, etc.)
- **Export Formats**: Additional output formats
- **Themes**: Color schemes and UI themes
- **Performance**: Optimization for large canvases and long animations
- **Documentation**: Tutorials, examples, and guides

## Design Philosophy

aart follows these principles:

1. **Keyboard-First**: All features accessible via keyboard shortcuts
2. **Vim-Inspired**: Modal editing with familiar keybindings
3. **Terminal Native**: Works beautifully in any terminal emulator
4. **Fast & Responsive**: Immediate visual feedback
5. **Zen Interface**: Clean, distraction-free UI that gets out of your way

## Credits

Built with:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style and layout
- [Bubble Tea Components](https://github.com/charmbracelet/bubbles) - UI components

Inspired by:
- vim's modal editing
- asciinema's terminal recording
- The demoscene ASCII art community

## Design Reference

See `notes/designs/design_reference.md` for the complete UI design specification.

## License

MIT License - see LICENSE file for details.

---

Made with ❤️ for terminal enthusiasts and ASCII artists

## GIF Import

Import animated GIFs and convert them to ASCII art animations with multiple conversion methods.

### Basic Usage

```bash
# Import from local file
./aart --import-gif animation.gif

# Import from URL
./aart --import-gif https://example.com/animation.gif

# Custom dimensions
./aart --import-gif animation.gif --width 120 --height 40

# Auto-size to terminal
./aart --import-gif animation.gif --ratio fill

# High FPS conversion
./aart --import-gif animation.gif --fps 60 --method dither
```

### Conversion Methods

- **luminosity** (default): Brightness-based using extended character set ` .:-=+*#%@`
- **block**: Block characters (░▒▓█) for solid, blocky appearance
- **edge**: Edge detection with line characters for wireframe effect
- **dither**: Floyd-Steinberg dithering for smoother gradients

### Options

- `--import-gif <url|file>`: Source GIF file or URL
- `--width <int>`: Canvas width (auto-detected from terminal if not specified)
- `--height <int>`: Canvas height (auto-detected from terminal if not specified)
- `--fps <int>`: Target FPS (default: matches source GIF)
- `--method <string>`: Conversion method (luminosity|block|edge|dither)
- `--ratio <string>`: Aspect ratio mode (fill|original)
- `--chars <string>`: Custom character ramp for luminosity method
- `--output <file>`: Save to .aa file instead of opening editor

### Examples

```bash
# Small blocky animation
./aart --import-gif dance.gif --width 40 --height 20 --method block

# Large detailed conversion that fills terminal
./aart --import-gif movie.gif --ratio fill --method luminosity

# High FPS with dithering
./aart --import-gif smooth.gif --fps 60 --method dither --width 200 --height 60

# Custom character ramp
./aart --import-gif art.gif --chars " .'`^\",:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
```
