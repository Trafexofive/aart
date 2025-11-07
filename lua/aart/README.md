# aart.nvim - ASCII Art Animations for Neovim

Neovim plugin to integrate ASCII art animations into your editor using the `aart` binary.

## Features

- 🎨 **Static frames** - Display any frame from `.aart` animations
- 🎬 **Animated buffers** - Play full animations in Neovim buffers
- 🖼️ **Alpha-nvim integration** - Use animations in your dashboard
- ⚡ **Cached loading** - Fast performance with animation caching
- 🎯 **Floating windows** - Open animations in centered floating windows
- 📦 **Export frames** - Extract frames as text files

## Installation

### Prerequisites

1. Install the `aart` binary and ensure it's in your PATH:
   ```bash
   go install github.com/mlamkadm/aart/cmd/aart@latest
   # or build from source
   ```

### Using lazy.nvim

Add to your Neovim config:

```lua
{
    dir = "~/path/to/aart/lua/aart",  -- or use git URL when published
    lazy = false,
    config = function()
        require('aart').setup({
            binary_path = "aart",  -- or full path
            cache_enabled = true,
            auto_start = true,
            frame_delay = 83,  -- ~12 fps
        })
        require('aart').setup_commands()
    end
}
```

### Using packer.nvim

```lua
use {
    '~/path/to/aart/lua/aart',  -- or git URL
    config = function()
        require('aart').setup()
        require('aart').setup_commands()
    end
}
```

## Usage

### Commands

```vim
" Open animation in floating window
:AartOpen ~/animations/logo.aart

" Export frames to directory
:AartExport ~/animations/logo.aart ./output_frames
```

### Lua API

#### Setup

```lua
require('aart').setup({
    binary_path = "aart",      -- Path to aart binary
    default_animation = nil,    -- Default animation file
    frame_delay = 83,          -- Frame delay in ms (~12fps)
    auto_start = true,         -- Auto-start animations
    cache_enabled = true,      -- Cache loaded animations
})
```

#### Load and Display Static Frame

Perfect for alpha-nvim or other dashboards:

```lua
local aart = require('aart')

-- Get first frame as table of strings
local ascii_art = aart.get_static_frame("~/logo.aart", 0)

-- Use in alpha-nvim
local header_lines = aart.get_alpha_animation("~/logo.aart", {
    frame = 0,      -- Frame index (0-based)
    padding = 1,    -- Add blank lines above/below
})
```

#### Create Animated Buffer

```lua
local aart = require('aart')

-- Create and open animated buffer in floating window
local bufnr = aart.create_animated_buffer("~/spinner.aart", {
    open_window = true,
    auto_start = true,
    loop = true,
    width = 60,
    height = 20,
})
```

#### Manual Animation Control

```lua
local aart = require('aart')

-- Load animation
local animation = aart.load_animation("~/logo.aart")
-- Returns: { frames = {...}, fps = 12, frame_count = 150 }

-- Play in existing buffer
aart.play_animation(bufnr, animation, {
    frame_delay = 83,  -- Override delay
    loop = true,       -- Loop animation
})

-- Stop animation
aart.stop_animation(bufnr)
```

## Integration Examples

### Alpha-nvim Dashboard

See `alpha-integration-example.lua` for a complete example. Quick snippet:

```lua
-- In your alpha-nvim config
local aart = require('aart')
aart.setup()

local function get_animated_header()
    local animation_file = vim.fn.expand("~/.config/nvim/logo.aart")
    return aart.get_alpha_animation(animation_file, { 
        frame = 0,
        padding = 1 
    })
end

-- Use in alpha dashboard layout
{
    type = 'text',
    val = get_animated_header(),
    opts = { hl = 'AlphaHeader', position = 'center' }
}
```

### Custom Startup Screen

```lua
vim.api.nvim_create_autocmd('VimEnter', {
    callback = function()
        if vim.fn.argc() == 0 then
            -- Show animation on startup
            local aart = require('aart')
            aart.create_animated_buffer("~/startup.aart", {
                open_window = false,  -- Use current window
                auto_start = true,
                loop = true,
            })
        end
    end
})
```

### Keybindings

```lua
local aart = require('aart')

vim.keymap.set('n', '<leader>aa', function()
    aart.open_animation("~/logo.aart")
end, { desc = "Show ASCII animation" })

vim.keymap.set('n', '<leader>ae', function()
    aart.export_frames("~/logo.aart", "./frames")
end, { desc = "Export animation frames" })
```

## Creating Animations

Use the `aart` CLI to create and edit animations:

```bash
# Create new animation
aart new --width 80 --height 24 --fps 12 --frames 30

# Import from GIF
aart import gif input.gif --output logo.aart

# Edit existing
aart edit logo.aart
```

## Configuration

### Default Configuration

```lua
{
    binary_path = "aart",           -- Path to aart binary
    default_animation = nil,         -- Default .aart file
    frame_delay = 83,               -- ms between frames (~12fps)
    auto_start = true,              -- Auto-start on buffer create
    cache_enabled = true,           -- Cache loaded animations
}
```

### Performance Tips

1. **Enable caching** - Animations are loaded once and cached
2. **Use static frames** for dashboards - Much faster than animation
3. **Optimize frame rate** - 12fps (83ms) is usually sufficient
4. **Pre-export frames** - For very large animations

## API Reference

### Setup Functions

- `setup(opts)` - Configure plugin
- `setup_commands()` - Register `:AartOpen` and `:AartExport` commands

### Loading Functions

- `load_animation(filepath)` - Load full animation with all frames
- `get_static_frame(filepath, frame_index)` - Get single frame (fast)
- `get_alpha_animation(filepath, opts)` - Get frame formatted for alpha-nvim

### Display Functions

- `create_animated_buffer(filepath, opts)` - Create new buffer with animation
- `play_animation(bufnr, animation, opts)` - Play animation in buffer
- `stop_animation(bufnr)` - Stop playing animation
- `open_animation(filepath)` - Open in floating window (command wrapper)

### Utility Functions

- `export_frames(filepath, output_dir)` - Export frames as .txt files

## Troubleshooting

### Animation not loading

1. Check file exists: `:echo filereadable('~/logo.aart')`
2. Verify aart binary: `:!which aart`
3. Test export manually: `aart export --input logo.aart --output /tmp/test --format txt`

### Slow performance

1. Enable caching: `cache_enabled = true`
2. Use static frames for dashboards instead of animations
3. Reduce frame rate or frame count

### Binary not found

Set full path in setup:
```lua
require('aart').setup({
    binary_path = "/full/path/to/aart"
})
```

## Examples

See the `examples.lua` file for more detailed usage examples.

## License

Same as the aart project.

## Contributing

Contributions welcome! Please test with various animation files and Neovim versions.
