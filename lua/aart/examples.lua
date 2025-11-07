-- Simple standalone example - use aart animations without alpha-nvim
-- Just require and use directly in your config

local aart = require('aart')

-- Setup the plugin
aart.setup({
    binary_path = "aart",  -- Path to aart binary
    cache_enabled = true,   -- Cache loaded animations
    auto_start = true,      -- Auto-start animations
    frame_delay = 83,       -- ~12 fps
})

-- Register commands
aart.setup_commands()

-- Example 1: Open an animation in a floating window
-- Usage in Neovim: :AartOpen ~/.config/nvim/animations/mylogo.aart
vim.keymap.set('n', '<leader>aa', function()
    aart.open_animation(vim.fn.expand("~/animations/logo.aart"))
end, { desc = "Open ASCII animation" })

-- Example 2: Use in a custom dashboard (without alpha-nvim)
local function show_custom_dashboard()
    -- Create a scratch buffer
    local bufnr = vim.api.nvim_create_buf(false, true)
    vim.api.nvim_buf_set_option(bufnr, 'buftype', 'nofile')
    vim.api.nvim_buf_set_option(bufnr, 'bufhidden', 'wipe')
    
    -- Load animation (just first frame for static display)
    local animation_file = vim.fn.expand("~/animations/logo.aart")
    local ascii_art = aart.get_static_frame(animation_file, 0)
    
    -- Add some menu items below
    local content = vim.list_extend({}, ascii_art)
    table.insert(content, "")
    table.insert(content, "Welcome to Neovim!")
    table.insert(content, "")
    table.insert(content, "Press 'f' to find files")
    table.insert(content, "Press 'r' for recent files")
    table.insert(content, "Press 'q' to quit")
    
    -- Center the content
    local padding = math.floor((vim.o.lines - #content) / 2)
    for i = 1, padding do
        table.insert(content, 1, "")
    end
    
    -- Set buffer content
    vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, content)
    vim.api.nvim_buf_set_option(bufnr, 'modifiable', false)
    
    -- Display in current window
    vim.api.nvim_set_current_buf(bufnr)
    
    -- Set up keymaps for the dashboard
    local opts = { buffer = bufnr, silent = true }
    vim.keymap.set('n', 'f', ':Telescope find_files<CR>', opts)
    vim.keymap.set('n', 'r', ':Telescope oldfiles<CR>', opts)
    vim.keymap.set('n', 'q', ':qa<CR>', opts)
end

-- Example 3: Create an animated buffer programmatically
vim.keymap.set('n', '<leader>ap', function()
    local filepath = vim.fn.expand("~/animations/spinner.aart")
    aart.create_animated_buffer(filepath, {
        open_window = true,
        auto_start = true,
        loop = true,
        width = 60,
        height = 20,
    })
end, { desc = "Play ASCII animation" })

-- Example 4: Export frames for inspection
vim.keymap.set('n', '<leader>ae', function()
    local filepath = vim.fn.input("Animation file: ", "", "file")
    if filepath ~= "" then
        local output_dir = vim.fn.input("Output directory: ", vim.fn.getcwd() .. "/frames")
        aart.export_frames(filepath, output_dir)
    end
end, { desc = "Export animation frames" })

-- Example 5: Show animation on VimEnter (startup)
vim.api.nvim_create_autocmd('VimEnter', {
    callback = function()
        if vim.fn.argc() == 0 then
            -- No files opened, show animation
            show_custom_dashboard()
        end
    end
})

return {
    show_dashboard = show_custom_dashboard,
}
