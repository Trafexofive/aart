# File Formats Documentation

## Native Format: .aart

The `.aart` format is aart's native JSON-based file format for storing ASCII art animations.

### Structure

```json
{
  "version": "1.0",
  "metadata": {
    "title": "My Animation",
    "author": "Artist Name",
    "description": "Description of the animation",
    "created": "2025-11-05T02:23:25+01:00",
    "modified": "2025-11-05T02:23:25+01:00",
    "tags": ["tag1", "tag2"],
    "source": "original_file.gif"
  },
  "canvas": {
    "width": 80,
    "height": 24
  },
  "frames": [
    {
      "index": 0,
      "duration": 100,
      "name": "Frame 1",
      "cells": [
        [
          {
            "char": "█",
            "fg": "#FFFFFF",
            "bg": "#000000",
            "bold": false,
            "italic": false,
            "underline": false
          }
        ]
      ]
    }
  ],
  "layers": [
    {
      "name": "Background",
      "visible": true,
      "opacity": 1.0,
      "blend_mode": "normal"
    }
  ],
  "palette": [
    {
      "name": "White",
      "hex": "#FFFFFF"
    }
  ],
  "audio": {
    "path": "audio.mp3",
    "offset": 0,
    "loop": true
  }
}
```

### Features

- **Metadata**: Title, author, description, timestamps, tags
- **Canvas**: Width and height definition
- **Frames**: Array of animation frames with duration
- **Cells**: Individual character cells with styling
- **Layers**: Layer support for complex compositions
- **Palette**: Color palette definition
- **Audio**: Optional audio track (future feature)

### Cell Properties

- `char`: UTF-8 character (string)
- `fg`: Foreground color (hex format)
- `bg`: Background color (hex format)
- `bold`: Bold styling (optional)
- `italic`: Italic styling (optional)
- `underline`: Underline styling (optional)

## Export Formats

### JSON (.json)

Exports the complete .aart structure as JSON.

```bash
aart --export output.json --export-format json input.aart
```

**Use cases:**
- Data processing
- Web applications
- API integration
- Programmatic access

### CSV (.csv)

Exports frame data as comma-separated values.

```bash
aart --export output.csv --export-format csv input.aart
```

**Format:**
```csv
frame,x,y,char,fg,bg
0,0,0,█,#FFFFFF,#000000
0,1,0,▒,#CCCCCC,#000000
```

**Use cases:**
- Spreadsheet analysis
- Data science / ML training data
- Bulk processing
- Statistical analysis

### Plain Text (.txt)

Exports as plain text with optional metadata.

```bash
aart --export output.txt --export-format txt input.aart
```

**Features:**
- Optional metadata header
- Pure character output
- No formatting/colors
- Human-readable

**Use cases:**
- Simple viewing
- Text-based terminals
- Documentation
- Archival

### ANSI (.ansi / .ans)

Exports with ANSI color codes.

```bash
aart --export output.ans --export-format ansi input.aart
```

**Features:**
- ANSI escape sequences
- 256-color support
- Terminal-compatible
- Preserves colors

**Use cases:**
- Terminal display
- BBS systems
- ANSI art sharing
- Retro computing

### HTML (.html)

Exports as standalone HTML file.

```bash
aart --export output.html --export-format html input.aart
```

**Features:**
- Self-contained HTML
- Inline CSS styling
- Color preservation
- Web-ready

**Use cases:**
- Web pages
- Email signatures
- Blogs
- Social sharing

### SVG (.svg)

Exports as Scalable Vector Graphics.

```bash
aart --export output.svg --export-format svg input.aart
```

**Features:**
- Vector format
- Scalable
- Color support
- Monospace font

**Use cases:**
- High-quality printing
- Scalable graphics
- Professional design
- Logos/branding

## Export Options

### Frame Selection

Export specific frame or all frames:

```bash
# Export all frames (default)
aart --export output.csv --export-format csv input.aart

# Export frame 0 only
aart --export output.txt --export-format txt --export-frame 0 input.aart

# Export frame 5
aart --export frame5.html --export-format html --export-frame 5 input.aart
```

### Color Control

Enable or disable colors in export:

```bash
# With colors (default)
aart --export output.html --export-format html input.aart

# Without colors
aart --export output.txt --export-format txt --export-colors=false input.aart
```

## Import Formats

### GIF (.gif)

Import animated GIFs to .aart format:

```bash
# From file
aart --import-gif animation.gif --output animation.aart

# From URL
aart --import-gif https://example.com/animation.gif --output animation.aart

# With options
aart --import-gif animation.gif \
  --width 100 \
  --height 40 \
  --method block \
  --fps 24 \
  --output animation.aart
```

**Conversion Methods:**
- `luminosity`: Brightness-based (default)
- `block`: Block characters (░▒▓█)
- `edge`: Edge detection
- `dither`: Dithered patterns

## Workflow Examples

### GIF to Multiple Formats

```bash
# Import GIF
aart --import-gif animation.gif --output animation.aart

# Export to various formats
aart --export animation.json --export-format json animation.aart
aart --export animation.csv --export-format csv animation.aart
aart --export animation.txt --export-format txt animation.aart
aart --export animation.html --export-format html animation.aart
aart --export animation.svg --export-format svg animation.aart
```

### Single Frame Exports

```bash
# Import multi-frame GIF
aart --import-gif animation.gif --output animation.aart

# Export first frame as PNG preview
aart --export frame0.html --export-format html --export-frame 0 animation.aart

# Export each frame
for i in {0..7}; do
  aart --export frame${i}.txt --export-format txt --export-frame $i animation.aart
done
```

### Data Analysis Pipeline

```bash
# Import
aart --import-gif data.gif --output data.aart

# Export to CSV
aart --export data.csv --export-format csv data.aart

# Process with pandas/R/Excel
python analyze_ascii.py data.csv
```

## File Size Comparison

Typical sizes for 80x24 canvas, 24 frames:

| Format | Size | Ratio | Use Case |
|--------|------|-------|----------|
| .aart | 120KB | 1.0x | Native, editable |
| .json | 120KB | 1.0x | Data exchange |
| .csv | 200KB | 1.7x | Analysis |
| .txt | 48KB | 0.4x | Simple view |
| .html | 150KB | 1.25x | Web display |
| .svg | 180KB | 1.5x | Vector graphics |
| .ans | 52KB | 0.43x | Terminal |

## Best Practices

### Choosing Format

- **Editing**: Use .aart (native format)
- **Viewing**: Use .txt or .html
- **Web**: Use .html or .svg
- **Terminal**: Use .ansi
- **Analysis**: Use .csv or .json
- **Archival**: Use .aart (preserves all data)

### Performance

- **Large files**: Use --export-frame to export single frames
- **Batch processing**: Export to CSV then process
- **Web display**: Use HTML with lazy loading for animations
- **Storage**: .aart is most efficient for complete data

### Compatibility

- **.aart**: aart only (future: other tools may support)
- **JSON**: Universal
- **CSV**: Universal (Excel, pandas, R)
- **TXT**: Universal
- **HTML**: All web browsers
- **SVG**: All modern browsers, design tools
- **ANSI**: Terminal emulators, BBS

## Future Formats

Planned support:

- **PNG**: Export as rasterized image
- **GIF**: Export back to animated GIF
- **MP4**: Export as video
- **WEBP**: Modern web format
- **PDF**: Print-ready document
- **LaTeX**: Academic papers
- **Markdown**: Documentation

## API Usage

```go
import "github.com/mlamkadm/aart/internal/fileformat"

// Load .aart file
aart, err := fileformat.Load("animation.aart")

// Export to format
opts := fileformat.ExportOptions{
    Format:      fileformat.FormatJSON,
    FrameIndex:  -1, // all frames
    IncludeMeta: true,
    Colors:      true,
}
err = fileformat.Export(aart, "output.json", opts)

// Create new .aart file
aart := fileformat.NewAartFile(80, 24, "My Animation")
aart.AddFrame(cells, 100) // 100ms duration
err = fileformat.Save("output.aart", aart)
```

## Troubleshooting

### Invalid .aart File

```bash
# Validate structure
aart --export /dev/null --export-format json input.aart
```

### Large Files

```bash
# Export single frame to reduce size
aart --export frame.txt --export-format txt --export-frame 0 large.aart
```

### Color Issues

```bash
# Disable colors if not displaying correctly
aart --export output.txt --export-format txt --export-colors=false input.aart
```

## Migration

### From Old Format

If you have files in the old text-based format:

```bash
# Old format is no longer supported
# Re-import from original GIF
aart --import-gif original.gif --output new_format.aart
```

### Version Updates

The .aart format version is in the file. Future versions will be backwards compatible with automatic migration on load.
