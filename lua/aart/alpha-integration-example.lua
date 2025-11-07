-- Example integration of aart.nvim with alpha-nvim
-- Drop this in your Neovim config alongside the alpha-nvim plugin config

return {
    'goolord/alpha-nvim',
    event = "VimEnter",
    dependencies = { 
        'nvim-tree/nvim-web-devicons',
        -- Add aart as a dependency (assumes it's in runtimepath)
        -- You can also load it manually with require('aart').setup()
    },
    config = function()
        local alpha = require('alpha')
        local dashboard = require('alpha.themes.dashboard')
        
        -- Setup aart plugin
        local aart = require('aart')
        aart.setup({
            binary_path = "aart", -- or full path like "/usr/local/bin/aart"
            cache_enabled = true,
        })

        -- Set theme colors
        local highlights = {
            AlphaHeader   = { fg = '#7AA2F7', bold = true },  -- Tokyo Night blue
            AlphaButtons  = { fg = '#C0CAF5' },
            AlphaFooter   = { fg = '#BB9AF7', italic = true },
            AlphaShortcut = { fg = '#FF79C6', bold = true },
            AlphaLeftPanel = { fg = '#9ECE6A' },  -- green
            AlphaRightPanel = { fg = '#E0AF68' }, -- yellow
        }
        for group, hl in pairs(highlights) do
            vim.api.nvim_set_hl(0, group, hl)
        end

        -- Helper to get calendar output
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

        -- Helper to get system info
        local function get_system_info()
            local info = {}
            
            -- Get neovim version
            local version = vim.version()
            table.insert(info, string.format("Neovim v%d.%d.%d", version.major, version.minor, version.patch))
            table.insert(info, "")
            
            -- Get plugin count
            local lazy_ok, lazy = pcall(require, "lazy")
            if lazy_ok then
                local plugins = lazy.plugins()
                table.insert(info, string.format("󰏗 %d plugins", #plugins))
            end
            
            -- Get git status if in repo
            local git_branch = vim.fn.system("git branch --show-current 2>/dev/null"):gsub("\n", "")
            if git_branch and git_branch ~= "" then
                table.insert(info, "")
                table.insert(info, " " .. git_branch)
            end
            
            return info
        end

        -- Create fancy buttons with better styling
        local function button(sc, txt, cmd)
            local b = dashboard.button(sc, txt, cmd)
            b.opts.hl = "AlphaButtons"
            b.opts.hl_shortcut = "AlphaShortcut"
            return b
        end

        -- Get ASCII animation from aart
        -- You can specify any .aart file here
        local function get_animation_header()
            -- Option 1: Use a specific animation file
            local animation_file = vim.fn.expand("~/.config/nvim/animations/logo.aart")
            
            -- Check if file exists, otherwise use a default
            if vim.fn.filereadable(animation_file) == 0 then
                -- Fallback to default ASCII art
                return {
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
            
            -- Load animation from aart file
            -- This gets the first frame as static ASCII art
            local ascii_art = aart.get_alpha_animation(animation_file, {
                frame = 0,  -- Use first frame (or any frame number)
                padding = 1, -- Add padding above/below
            })
            
            return ascii_art
        end

        -- Build the 3-column layout with animated header
        local function create_layout()
            local calendar = get_calendar()
            local sysinfo = get_system_info()
            local header_lines = get_animation_header()
            
            -- Calculate max widths
            local cal_width = 0
            for _, line in ipairs(calendar) do
                cal_width = math.max(cal_width, vim.fn.strdisplaywidth(line))
            end
            
            local info_width = 0
            for _, line in ipairs(sysinfo) do
                info_width = math.max(info_width, vim.fn.strdisplaywidth(line))
            end
            
            -- Get center width from animation
            local center_width = 0
            for _, line in ipairs(header_lines) do
                center_width = math.max(center_width, vim.fn.strdisplaywidth(line))
            end
            
            local max_rows = math.max(#calendar, #header_lines, #sysinfo)
            local combined = {}
            
            for i = 1, max_rows do
                local left = calendar[i] or string.rep(" ", cal_width)
                local center = header_lines[i] or string.rep(" ", center_width)
                local right = sysinfo[i] or ""
                
                -- Pad left column to width + spacing
                left = left .. string.rep(" ", cal_width - vim.fn.strdisplaywidth(left) + 8)
                
                -- Pad center column
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
                        button('a', '🎨  ASCII Animations', ':AartOpen<CR>'),  -- Add aart command!
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
        
        -- Setup aart commands
        aart.setup_commands()

        -- Only show alpha when opening with no args
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

        -- Auto-close alpha when entering a real buffer
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
