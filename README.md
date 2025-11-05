# aart - ASCII Art Animation Editor

A terminal-based ASCII art animation editor built with Go and Bubbletea.

## Features

### Currently Working ✅
- **Canvas**: 80x24 character canvas with Braille art support
- **Timeline**: 24-frame animation timeline with visual frame indicator
- **Playback**: Play/pause animations at 12fps
- **Navigation**: hjkl/arrow keys for cursor movement
- **Drawing**: 
  - Press `d` to place current character
  - Press `i` for insert mode - type characters directly
- **Frame Navigation**: `,` and `.` to seek between frames
- **Tool Selection**: p/f/s/l/b/t/e/m for different tools (pencil is active)
- **View Controls**: 
  - `g` for grid toggle
  - `+/-` for zoom
  - `0` to reset zoom
- **Radial Wheel**: `ctrl-j`/`ctrl-k` to cycle sections, `enter` to expand, `esc` to collapse
- **Command Mode**: `:` to enter command mode
- **Status Bar**: Shows file info, frame count, tool, FPS, layer info

### Radial Wheel Sections
- **HELP**: Keyboard shortcuts
- **TOOLS**: Pencil, fill, select, line, box, text, eyedropper, move
- **COLORS**: Color picker (stubbed)
- **EXPORT**: Export formats (stubbed)
- **IMPORT**: Import options (stubbed)
- **LAYERS**: Layer management (stubbed)

## Build & Run

```bash
# Install dependencies
go mod download

# Build
make build

# Run
make run
# or
./aart
```

## Keyboard Shortcuts

### Normal Mode
- `hjkl` or arrow keys - Move cursor
- `d` - Draw current character at cursor
- `i` - Enter insert mode
- `space` - Play/pause animation
- `,` `.` - Seek to previous/next frame
- `+` `-` - Zoom in/out
- `0` - Reset zoom
- `g` - Toggle grid
- `p` `f` `s` `l` `b` `t` `e` `m` - Select tool
- `ctrl-j` `ctrl-k` - Cycle wheel sections
- `enter` - Expand wheel section
- `esc` - Collapse wheel
- `:` - Enter command mode
- `q` - Quit

### Insert Mode
- Any character - Type at cursor position
- `esc` - Return to normal mode

### Command Mode
- Type command and press `enter`
- `esc` - Cancel and return to normal mode

## Architecture

```
cmd/aart/          - Main entry point
internal/ui/       - Bubbletea UI model and rendering
internal/canvas/   - Canvas and frame management (TODO)
internal/timeline/ - Timeline and playback (TODO)
```

## Next Steps

Based on STARTING_BUILD_REF.md:

1. **Drawing Engine**
   - Persistent drawing in insert mode
   - Tool implementations (fill, line, box, selection)
   
2. **Layers**
   - Multiple layers per frame
   - Opacity and blending
   - Layer visibility toggle

3. **Export/Import**
   - ANSI format (.ans)
   - Plain text (.txt)
   - Native format (.aart)
   - Image import with character conversion

4. **Advanced Canvas**
   - Actual zoom rendering
   - Grid overlay display
   - Undo/redo stack
   - Copy/paste

5. **UI Polish**
   - Wheel rendering fix for wide terminals
   - Color picker implementation
   - Better timeline for >24 frames
   - Mouse support

## Design Reference

See `notes/designs/design_reference.md` for the complete UI design specification.

## License

MIT

## GIF Import Feature

Import animated GIFs and convert them to ASCII art animations!

### Usage

```bash
# Import from local file
aart --import-gif animation.gif

# Import from URL
aart --import-gif https://example.com/animation.gif

# Customize conversion
aart --import-gif animation.gif --width 120 --height 40 --method block

# Save to file instead of opening editor
aart --import-gif animation.gif --output converted.aart
```

### Conversion Methods

- **luminosity** (default): Converts based on pixel brightness using extended character set
- **block**: Uses block characters (░▒▓█) for a solid, blocky look
- **edge**: Edge detection with line characters for wireframe effect
- **dither**: Dithered output for smoother gradients

### Options

- `--width <int>`: Canvas width (default: 80)
- `--height <int>`: Canvas height (default: 24)
- `--fps <int>`: Target FPS (default: 12)
- `--method <string>`: Conversion method
- `--chars <string>`: Custom character set for conversion
- `--output <file>`: Save to file instead of opening editor

### Examples

```bash
# Small, blocky animation
aart --import-gif dance.gif --width 40 --height 20 --method block

# Large, detailed conversion
aart --import-gif movie.gif --width 160 --height 60 --method luminosity

# Custom character ramp
aart --import-gif art.gif --chars " .:-=+*#%@"
```
