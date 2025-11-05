# Configuration System Documentation

## Overview

aart uses a comprehensive YAML-based configuration system stored in `~/.config/aart/`. The system supports themes, templates, recent files tracking, and extensive customization options.

## Directory Structure

```
~/.config/aart/
├── config.yml        # Main configuration file
├── themes/           # Color themes
│   └── gruvbox.yml   # Example theme
├── templates/        # Animation templates
├── recent/           # Recent files cache
└── backups/          # Auto-save backups
```

## Initialization

```bash
# Initialize configuration directory
aart --init

# Show config file path
aart --config-path

# Display current configuration
aart --show-config
```

## Configuration File

Location: `~/.config/aart/config.yml`

### Full Example

```yaml
version: 0.1.0

editor:
  default_width: 80
  default_height: 24
  default_fps: 12
  auto_save: false
  auto_save_interval: 300  # seconds
  tab_size: 4
  show_grid: false
  show_line_numbers: false
  zen_mode: false

ui:
  theme: dark  # dark, light, custom
  show_status_bar: true
  show_timeline: true
  show_wheel_by_default: false
  cursor_style: line  # block, line, underline
  animation_smooth: true
  progress_style: bar  # bar, spinner, minimal

colors:
  name: default
  foreground: '#FFFFFF'
  background: '#000000'
  cursor: '#FFFF00'
  selection: '#444444'
  status_bar: '#333333'
  timeline: '#222222'
  border: '#666666'
  custom:
    accent: '#00FF00'
    error: '#FF0000'

recent:
  files: []
  max_entries: 10

converter:
  default_method: luminosity  # luminosity, block, edge, dither
  default_chars: ""
  preserve_aspect: true
  quality: high  # low, medium, high

keybindings:
  play: " "
  quit: "q"
  save: "ctrl+s"
  zen_mode: "z"
  next_frame: "."
  prev_frame: ","
```

## Settings Reference

### Editor Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `default_width` | int | 80 | Default canvas width |
| `default_height` | int | 24 | Default canvas height |
| `default_fps` | int | 12 | Default animation FPS |
| `auto_save` | bool | false | Enable auto-save |
| `auto_save_interval` | int | 300 | Auto-save interval (seconds) |
| `tab_size` | int | 4 | Tab size for indentation |
| `show_grid` | bool | false | Show grid overlay |
| `show_line_numbers` | bool | false | Show line numbers |
| `zen_mode` | bool | false | Start in zen mode |

### UI Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `theme` | string | "dark" | Color theme (dark/light/custom) |
| `show_status_bar` | bool | true | Display status bar |
| `show_timeline` | bool | true | Display timeline |
| `show_wheel_by_default` | bool | false | Show radial wheel on start |
| `cursor_style` | string | "line" | Cursor style (block/line/underline) |
| `animation_smooth` | bool | true | Smooth animations |
| `progress_style` | string | "bar" | Progress style (bar/spinner/minimal) |

### Color Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `name` | string | "default" | Color scheme name |
| `foreground` | string | "#FFFFFF" | Default foreground color |
| `background` | string | "#000000" | Default background color |
| `cursor` | string | "#FFFF00" | Cursor color |
| `selection` | string | "#444444" | Selection color |
| `status_bar` | string | "#333333" | Status bar color |
| `timeline` | string | "#222222" | Timeline color |
| `border` | string | "#666666" | Border color |
| `custom.*` | string | - | Custom color keys |

### Converter Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `default_method` | string | "luminosity" | Default conversion method |
| `default_chars` | string | "" | Default character set |
| `preserve_aspect` | bool | true | Preserve aspect ratio |
| `quality` | string | "high" | Conversion quality |

## Recent Files

Recent files are automatically tracked when:
- Importing and saving GIFs
- Opening existing files (TODO)
- Saving files (TODO)

Format:
```yaml
recent:
  files:
    - path: /path/to/file.aart
      timestamp: 2025-11-05T02:12:00Z
      frames: 8
  max_entries: 10
```

## Themes

Create custom themes in `~/.config/aart/themes/`:

### Example: Gruvbox Theme

```yaml
name: gruvbox
foreground: '#EBDBB2'
background: '#282828'
cursor: '#FE8019'
selection: '#504945'
status_bar: '#3C3836'
timeline: '#1D2021'
border: '#665C54'
custom:
  accent: '#B8BB26'
  error: '#FB4934'
  success: '#B8BB26'
  warning: '#FABD2F'
```

To use: Update `ui.theme: gruvbox` in config.yml (TODO: theme loading)

## CLI Integration

Configuration values are used as defaults when CLI flags are not provided:

```bash
# Uses config.editor.default_width (80)
aart --import-gif animation.gif

# Overrides config with explicit value
aart --import-gif animation.gif --width 120

# Uses config.converter.default_method
aart --import-gif animation.gif

# Overrides with explicit method
aart --import-gif animation.gif --method block
```

## Configuration Precedence

1. Command-line flags (highest priority)
2. config.yml settings
3. Default configuration (lowest priority)

## Customization Examples

### Large Canvas Default

```yaml
editor:
  default_width: 160
  default_height: 60
  default_fps: 24
```

### Block Art Workflow

```yaml
converter:
  default_method: block
  quality: high

editor:
  default_width: 100
  default_height: 40
```

### Minimalist UI

```yaml
ui:
  show_status_bar: false
  show_timeline: false
  progress_style: minimal

editor:
  zen_mode: true
```

### Auto-Save Enabled

```yaml
editor:
  auto_save: true
  auto_save_interval: 60  # Every minute
```

## Environment Variables

- `XDG_CONFIG_HOME`: Override config directory location
  - Default: `~/.config`
  - Custom: `XDG_CONFIG_HOME=/custom/path aart`

## Configuration Management

### Backup Config

```bash
cp ~/.config/aart/config.yml ~/.config/aart/config.yml.backup
```

### Reset to Defaults

```bash
rm ~/.config/aart/config.yml
aart --init
```

### Validate Config

```bash
aart --show-config
```

## Future Enhancements

- Theme loading from themes/ directory
- Template system for common animations
- Per-project configuration (.aart.yml)
- Config migration for version updates
- Theme marketplace/sharing
- Workspace presets
- Plugin configuration
- Keybinding customization UI
- Config validation with helpful errors

## Troubleshooting

### Config Not Loading

```bash
# Check config path
aart --config-path

# Verify file exists
cat $(aart --config-path)

# Re-initialize
aart --init
```

### Invalid YAML

- Use `aart --show-config` to validate
- Check indentation (2 spaces)
- Ensure proper quoting of strings
- Validate colors start with '#'

### Defaults Not Applying

- Verify config.yml syntax
- Check for typos in keys
- Ensure values are correct type (int vs string)
- Use --show-config to see active values

## API Usage (for developers)

```go
import "github.com/mlamkadm/aart/internal/config"

// Load configuration
cfg, err := config.Load()

// Initialize config directory
err := config.Init()

// Save configuration
err := config.Save(cfg)

// Add recent file
cfg.AddRecentFile("/path/to/file.aart", 24)

// Get config directory
dir, err := config.ConfigDir()
```

## File Format

The configuration uses YAML 1.2 format:
- Indentation: 2 or 4 spaces (consistent)
- Comments: Lines starting with `#`
- Strings: Quoted or unquoted
- Colors: Hex format `#RRGGBB`
- Booleans: `true` or `false`
- Numbers: Unquoted integers

## Best Practices

1. **Backup before editing**: Always backup config.yml
2. **Small changes**: Test one setting at a time
3. **Use --show-config**: Verify changes take effect
4. **Comment your changes**: Document customizations
5. **Keep defaults**: Start from default config
6. **Test thoroughly**: Ensure app works after changes

## Support

For issues with configuration:
1. Check this documentation
2. Validate with `aart --show-config`
3. Reset with `aart --init`
4. Report bugs with config file attached
