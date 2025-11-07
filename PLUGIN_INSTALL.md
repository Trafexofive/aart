# Quick Install: aart.nvim Plugin

3-step process to get ASCII animations in your Neovim alpha dashboard.

## Step 1: Add Plugin to Neovim

Add this to your lazy.nvim plugins (e.g., `~/.config/nvim/lua/plugins/aart.lua`):

```lua
return {
    dir = vim.fn.expand("~/repos/aart/lua/aart"),  -- Adjust path to where you cloned aart
    name = "aart",
    lazy = false,
    config = function()
        require('aart').setup({
            cache_enabled = true,
        })
        require('aart').setup_commands()
    end
}
```

## Step 2: Update Alpha-nvim Config

In your alpha-nvim config, replace the header section:

```lua
-- Before (static ASCII art):
local header_lines = {
    '    ███╗   ██╗██╗   ██╗██╗███╗   ███╗',
    '    ████╗  ██║██║   ██║██║████╗ ████║',
    -- ...
}

-- After (dynamic from .aart file):
local aart = require('aart')
local animation_file = vim.fn.expand("~/.config/nvim/logo.aart")
local header_lines = {}

if vim.fn.filereadable(animation_file) == 1 then
    header_lines = aart.get_alpha_animation(animation_file, {
        frame = 0,    -- First frame
        padding = 0,  -- No extra padding
    })
else
    -- Fallback to your static art
    header_lines = {
        '    ███╗   ██╗██╗   ██╗██╗███╗   ███╗',
        '    ████╗  ██║██║   ██║██║████╗ ████║',
        -- ...
    }
end
```

## Step 3: Create Your Animation

Option A: Import from GIF
```bash
aart --import-gif ~/Downloads/mylogo.gif --output ~/.config/nvim/logo.aart
```

Option B: Use example from aart repo
```bash
cp ~/repos/aart/examples/logo.aart ~/.config/nvim/logo.aart
```

Option C: Create in TUI
```bash
aart
# Create your art, then :w ~/.config/nvim/logo.aart
```

## Test It

```vim
# Restart Neovim
nvim

# Or test manually:
:lua vim.print(require('aart').get_static_frame(vim.fn.expand("~/.config/nvim/logo.aart"), 0))
```

## Bonus: Add Animation Viewer

Add a button to your alpha dashboard to view the full animation:

```lua
button('a', '🎨  View Animation', ':AartOpen ~/.config/nvim/logo.aart<CR>'),
```

## Full Example

See `NEOVIM_PLUGIN.md` for the complete integration example with your 3-column layout.

## Troubleshooting

**"Module not found"**
- Make sure the path in `dir = ` points to where you cloned aart
- Try `:lua vim.print(vim.fn.expand("~/repos/aart/lua/aart"))`

**"Animation file not found"**
- Check: `:lua vim.print(vim.fn.filereadable(vim.fn.expand("~/.config/nvim/logo.aart")))`
- Should return `1`

**"Failed to parse"**
- Verify JSON: `jq . < ~/.config/nvim/logo.aart`
- Make sure file was created with aart binary

**No animation shows**
- Check fallback ASCII art is working
- Add debug: `:lua vim.print(require('aart').get_static_frame("~/.config/nvim/logo.aart", 0))`

## Next Steps

- Customize colors in alpha-nvim highlights
- Try different frames: `frame = 5` instead of `frame = 0`
- Create multiple animations for different themes
- Use `:AartOpen` command to view animations anytime
