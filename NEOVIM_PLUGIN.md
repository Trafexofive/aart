# aart.nvim - Neovim Plugin Integration

Complete guide for integrating ASCII art animations into Neovim using the aart Lua plugin.

## Overview

The aart Lua plugin allows you to:
- Display ASCII animations in Neovim buffers
- Integrate static or animated frames into alpha-nvim dashboards
- Play looping animations in floating windows
- Export frames programmatically

## Quick Start

### Installation

Using lazy.nvim:

```lua
{
    -- Add to your plugins
    dir = vim.fn.expand("~/repos/aart/lua/aart"),  -- or wherever you cloned aart
    name = "aart",
    lazy = false,
    config = function()
        require('aart').setup({
            binary_path = "aart",  -- Will use from PATH
            cache_enabled = true,   -- Cache loaded animations
        })
        require('aart').setup_commands()
    end
}
```

### Basic Usage

```lua
local aart = require('aart')

-- Get first frame of animation for dashboard
local logo = aart.get_static_frame("~/animations/logo.aart", 0)
-- Returns: { "line1", "line2", ... }

-- Load full animation
local animation = aart.load_animation("~/animations/spinner.aart")
-- Returns: { frames = {...}, fps = 12, frame_count = 150 }

-- Open in floating window
aart.open_animation("~/animations/logo.aart")
```

## Alpha-nvim Integration

Perfect for creating animated startup screens. Here's your exact use case:

```lua
return {
    'goolord/alpha-nvim',
    event = "VimEnter",
    dependencies = { 'nvim-tree/nvim-web-devicons' },
    config = function()
        local alpha = require('alpha')
        local dashboard = require('alpha.themes.dashboard')
        local aart = require('aart')
        
        -- Initialize aart
        aart.setup({ cache_enabled = true })

        -- Set theme colors (your existing colors)
        local highlights = {
            AlphaHeader   = { fg = '#00FFFF', bold = true },
            AlphaButtons  = { fg = '#FFFFFF' },
            AlphaFooter   = { fg = '#FF0000', italic = true },
            AlphaShortcut = { fg = '#FF00FF', bold = true },
            AlphaLeftPanel = { fg = '#8ec07c' },
            AlphaRightPanel = { fg = '#fabd2f' },
        }
        for group, hl in pairs(highlights) do
            vim.api.nvim_set_hl(0, group, hl)
        end

        -- Your existing helper functions
        local function get_calendar()
            local handle = io.popen('cal')
            if not handle then return {} end
            local result = handle:read("*a")
            handle:close()
            local lines = {}
            for line in result:gmatch("[^\r\n]+") do
                table.insert(lines, line)
            end
            return lines
        end

        local function get_system_info()
            local info = {}
            local version = vim.version()
            table.insert(info, string.format("Neovim v%d.%d.%d", version.major, version.minor, version.patch))
            table.insert(info, "")
            
            local lazy_ok, lazy = pcall(require, "lazy")
            if lazy_ok then
                local plugins = lazy.plugins()
                table.insert(info, string.format("󰏗 %d plugins", #plugins))
            end
            
            local git_branch = vim.fn.system("git branch --show-current 2>/dev/null"):gsub("\n", "")
            if git_branch and git_branch ~= "" then
                table.insert(info, "")
                table.insert(info, " " .. git_branch)
            end
            
            return info
        end

        -- Create fancy buttons (your existing function)
        local function button(sc, txt, cmd)
            local b = dashboard.button(sc, txt, cmd)
            b.opts.hl = "AlphaButtons"
            b.opts.hl_shortcut = "AlphaShortcut"
            return b
        end

        -- Build the 3-column layout with ANIMATED HEADER!
        local function create_layout()
            local calendar = get_calendar()
            local sysinfo = get_system_info()
            
            -- NEW: Get ASCII animation from aart
            local animation_file = vim.fn.expand("~/.config/nvim/logo.aart")
            local header_lines = {}
            
            -- Try to load animation, fallback to default if not found
            if vim.fn.filereadable(animation_file) == 1 then
                header_lines = aart.get_alpha_animation(animation_file, {
                    frame = 0,      -- Use first frame (or cycle through frames!)
                    padding = 0,    -- Add blank lines if needed
                })
            else
                -- Fallback ASCII art
                header_lines = {
                    '',
                    '    ███╗   ██╗██╗   ██╗██╗███╗   ███╗',
                    '    ████╗  ██║██║   ██║██║████╗ ████║',
                    '    ██╔██╗ ██║██║   ██║██║██╔████╔██║',
                    '    ██║╚██╗██║╚██╗ ██╔╝██║██║╚██╔╝██║',
                    '    ██║ ╚████║ ╚████╔╝ ██║██║ ╚═╝ ██║',
                    '    ╚═╝  ╚═══╝  ╚═══╝  ╚═╝╚═╝     ╚═╝',
                    '',
                }
            end
            
            -- Calculate max widths
            local cal_width = 0
            for _, line in ipairs(calendar) do
                cal_width = math.max(cal_width, vim.fn.strdisplaywidth(line))
            end
            
            local center_width = 0
            for _, line in ipairs(header_lines) do
                center_width = math.max(center_width, vim.fn.strdisplaywidth(line))
            end
            
            -- Combine into 3-column layout
            local max_rows = math.max(#calendar, #header_lines, #sysinfo)
            local combined = {}
            
            for i = 1, max_rows do
                local left = calendar[i] or string.rep(" ", cal_width)
                local center = header_lines[i] or string.rep(" ", center_width)
                local right = sysinfo[i] or ""
                
                left = left .. string.rep(" ", cal_width - vim.fn.strdisplaywidth(left) + 8)
                center = center .. string.rep(" ", 8)
                
                table.insert(combined, left .. center .. right)
            end
            
            return {
                { type = 'padding', val = 2 },
                {
                    type = 'text',
                    val = combined,
                    opts = { hl = 'AlphaHeader', position = 'center' }
                },
                { type = 'padding', val = 2 },
                {
                    type = 'group',
                    val = {
                        button('f', '󰈞  Find file', ':Telescope find_files<CR>'),
                        button('n', '  Sessions', ':Telescope session-lens<CR>'),
                        button('s', '  Restore Session', ':SessionRestore<CR>'),
                        button('r', '  Recent files', ':Telescope oldfiles<CR>'),
                        button('g', '󰊄  Find text', ':Telescope live_grep<CR>'),
                        button('a', '🎨  View Animation', ':AartOpen ~/.config/nvim/logo.aart<CR>'),  -- NEW!
                        button('t', '  TUI Commands', ':lua require("mlamkadm.core.terminal").show_tui_registry()<CR>'),
                        button('l', '󰒲  Lazy', '<cmd>Lazy<cr>'),
                        button('q', '  Quit', '<cmd>qa<CR>')
                    },
                    opts = { spacing = 1 }
                },
                { type = 'padding', val = 1 },
                {
                    type = 'text',
                    val = '"The Great Work Continues." - L\'bro',
                    opts = { hl = 'AlphaFooter', position = 'center' }
                },
            }
        end

        local config = {
            layout = create_layout(),
            opts = { margin = 5 }
        }

        alpha.setup(config)
        aart.setup_commands()  -- Adds :AartOpen command

        -- Your existing autocmds...
        vim.api.nvim_create_autocmd('VimEnter', {
            callback = function()
                if vim.fn.argc() == 0 then
                    local cwd = vim.fn.getcwd()
                    if cwd ~= vim.fn.expand('~') and cwd ~= '/' then
                        pcall(vim.cmd, 'SessionRestore')
                    end
                end
            end
        })

        vim.api.nvim_create_autocmd('BufEnter', {
            pattern = '*',
            callback = function()
                local bufnr = vim.api.nvim_get_current_buf()
                if vim.bo[bufnr].buftype == '' and vim.bo[bufnr].filetype ~= 'alpha' then
                    for _, win in ipairs(vim.api.nvim_list_wins()) do
                        local buf = vim.api.nvim_win_get_buf(win)
                        if vim.bo[buf].filetype == 'alpha' then
                            if #vim.api.nvim_list_wins() > 1 then
                                vim.api.nvim_win_close(win, true)
                            end
                            break
                        end
                    end
                end
            end
        })
    end
}
```

## Creating Animations

Use the aart TUI to create animations for your dashboard:

```bash
# Import from GIF
aart --import-gif ~/Downloads/logo.gif --output ~/.config/nvim/logo.aart

# Or create from scratch
aart  # Opens TUI editor
# Then save to ~/.config/nvim/logo.aart
```

## Advanced: Animated Cycling

Want to cycle through frames in your dashboard? Add this:

```lua
-- Cycle through frames on a timer
local current_frame = 0
local total_frames = 10  -- Adjust based on your animation

vim.fn.timer_start(100, function()  -- Update every 100ms
    current_frame = (current_frame + 1) % total_frames
    
    -- Recreate dashboard with new frame
    if vim.bo.filetype == 'alpha' then
        local new_layout = create_layout()  -- Will use updated current_frame
        require('alpha').setup({ layout = new_layout })
    end
end, { ['repeat'] = -1 })
```

## API Reference

### Setup

```lua
require('aart').setup({
    binary_path = "aart",      -- Path to binary (not needed for parsing)
    cache_enabled = true,       -- Cache loaded animations
    auto_start = true,          -- Auto-start animations in buffers
    frame_delay = 83,          -- Frame delay in ms (~12fps)
})
```

### Functions

#### `get_static_frame(filepath, frame_index)`
Get a single frame from animation (fast, for dashboards).

```lua
local frame = aart.get_static_frame("~/logo.aart", 0)
-- Returns: { "line1", "line2", ... }
```

#### `get_alpha_animation(filepath, opts)`
Get frame formatted for alpha-nvim.

```lua
local header = aart.get_alpha_animation("~/logo.aart", {
    frame = 0,      -- Frame index
    padding = 1,    -- Add blank lines
})
```

#### `load_animation(filepath)`
Load full animation with all frames.

```lua
local anim = aart.load_animation("~/spinner.aart")
-- Returns: { frames = {...}, fps = 12, frame_count = 150 }
```

#### `create_animated_buffer(filepath, opts)`
Create buffer with animation.

```lua
local bufnr = aart.create_animated_buffer("~/logo.aart", {
    open_window = true,   -- Open in floating window
    auto_start = true,    -- Start playing
    loop = true,          -- Loop forever
    width = 80,
    height = 24,
})
```

#### `open_animation(filepath)`
Quick command to open in floating window.

```lua
aart.open_animation("~/logo.aart")
```

#### `export_frames(filepath, output_dir)`
Export all frames as text files.

```lua
aart.export_frames("~/logo.aart", "/tmp/frames")
```

## Commands

After calling `setup_commands()`:

```vim
:AartOpen ~/logo.aart              " Open in floating window
:AartExport ~/logo.aart /tmp/out   " Export frames
```

## File Structure

```
~/.config/nvim/
├── lua/
│   └── plugins/
│       └── alpha.lua              # Your alpha config (updated)
└── animations/                     # Your animations
    ├── logo.aart
    ├── spinner.aart
    └── welcome.aart

~/repos/aart/
└── lua/
    └── aart/
        ├── init.lua                      # Main plugin
        ├── README.md                     # Plugin docs
        ├── examples.lua                  # Usage examples
        └── alpha-integration-example.lua # Full alpha example
```

## Troubleshooting

### Animation not loading

Check file exists and is valid JSON:
```lua
:lua vim.print(vim.fn.filereadable(vim.fn.expand("~/logo.aart")))
:lua vim.print(vim.json.decode(io.open(vim.fn.expand("~/logo.aart")):read("*all")))
```

### Slow performance

Enable caching:
```lua
require('aart').setup({ cache_enabled = true })
```

Use static frames for dashboards (much faster than loading full animations).

## Examples Directory

See `lua/aart/examples.lua` for more usage patterns.

## Creating Your First Animation

1. Create or import an animation:
   ```bash
   aart --import-gif ~/coollogo.gif --output ~/.config/nvim/logo.aart
   ```

2. Test in Neovim:
   ```vim
   :lua vim.print(require('aart').get_static_frame(vim.fn.expand("~/.config/nvim/logo.aart"), 0))
   ```

3. Add to alpha-nvim config (see example above)

4. Restart Neovim and enjoy your animated dashboard!

## Tips

- Keep animations under 100 frames for dashboards
- Use 80x24 or smaller for better compatibility
- Test frame 0 before committing to an animation
- Cache is automatic when enabled
- For static dashboards, use `get_static_frame()` instead of full animation

## Contributing

The plugin is designed to work without the aart binary by parsing JSON directly. This makes it lightweight and portable.

Improvements welcome:
- Better error handling
- Color support (if alpha-nvim supports it)
- Frame cycling for animated dashboards
- Performance optimizations
