-- aart-startup.nvim - 3-panel startup with animated center and selectable buttons
local M = {}

M.config = {
  animation_file = vim.fn.expand("~/.config/aart/test.aa"),
  aart_binary = vim.fn.expand("~/repos/aart/aart"),
  show_on_startup = true,
  auto_restore_session = true,
  orientation = "default",  -- "default" or "centered"
}

function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", M.config, opts or {})
  if M.config.show_on_startup then
    M.setup_autocommands()
  end
end

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
    table.insert(info, string.format("󰏗 %d plugins", #lazy.plugins()))
  end
  
  local git_branch = vim.fn.system("git branch --show-current 2>/dev/null"):gsub("\n", "")
  if git_branch and git_branch ~= "" then
    table.insert(info, "")
    table.insert(info, " " .. git_branch)
  end
  
  return info
end

-- Button definitions with actions
local buttons = {
  {key = 'f', icon = '󰈞', text = 'Find file', cmd = function() vim.cmd('only | Telescope find_files') end},
  {key = 'r', icon = '', text = 'Recent files', cmd = function() vim.cmd('only | Telescope oldfiles') end},
  {key = 's', icon = '', text = 'Sessions', cmd = function() vim.cmd('only | Telescope session-lens') end},
  {key = 'S', icon = '', text = 'Restore Session', cmd = function() vim.cmd('only'); pcall(vim.cmd, 'SessionRestore') end},
  {key = 't', icon = '', text = 'TUI Commands', cmd = function() vim.cmd('only'); require("mlamkadm.core.terminal").show_tui_registry() end},
  {key = 'l', icon = '󰒲', text = 'Lazy', cmd = function() vim.cmd('only | Lazy') end},
  {key = 'q', icon = '', text = 'Quit', cmd = function() vim.cmd('qa') end},
}

local function create_buttons_buffer()
  local buf = vim.api.nvim_create_buf(false, true)
  vim.bo[buf].bufhidden = 'wipe'
  vim.bo[buf].filetype = 'aart-buttons'
  
  -- Build button lines
  local lines = {"", "Quick Actions", "─────────────", ""}
  for i, btn in ipairs(buttons) do
    table.insert(lines, string.format("  [%s] %s  %s", btn.key, btn.icon, btn.text))
  end
  table.insert(lines, "")
  table.insert(lines, '"The Great Work Continues." - L\'bro')
  
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
  vim.bo[buf].modifiable = false
  
  return buf, lines
end

local function setup_buttons_navigation(buf)
  local selected = 1
  local max_items = #buttons
  local ns_id = vim.api.nvim_create_namespace('aart_buttons_highlight')
  
  -- Highlight selected button
  local function update_highlight()
    vim.api.nvim_buf_clear_namespace(buf, ns_id, 0, -1)
    local line_idx = selected + 3  -- offset for header lines
    vim.api.nvim_buf_add_highlight(buf, ns_id, 'Visual', line_idx, 0, -1)
    
    -- Move cursor to that line
    local wins = vim.fn.win_findbuf(buf)
    if #wins > 0 then
      vim.api.nvim_win_set_cursor(wins[1], {line_idx + 1, 0})
    end
  end
  
  -- Navigation
  vim.keymap.set('n', 'j', function()
    selected = selected % max_items + 1
    update_highlight()
  end, {buffer = buf, silent = true, desc = 'Next button'})
  
  vim.keymap.set('n', 'k', function()
    selected = selected - 1
    if selected < 1 then selected = max_items end
    update_highlight()
  end, {buffer = buf, silent = true, desc = 'Previous button'})
  
  vim.keymap.set('n', '<Down>', function()
    selected = selected % max_items + 1
    update_highlight()
  end, {buffer = buf, silent = true})
  
  vim.keymap.set('n', '<Up>', function()
    selected = selected - 1
    if selected < 1 then selected = max_items end
    update_highlight()
  end, {buffer = buf, silent = true})
  
  -- Execute selected
  vim.keymap.set('n', '<CR>', function()
    local btn = buttons[selected]
    if btn and btn.cmd then
      btn.cmd()
    end
  end, {buffer = buf, silent = true, desc = 'Execute button'})
  
  -- Direct key shortcuts
  for _, btn in ipairs(buttons) do
    vim.keymap.set('n', btn.key, function()
      if btn.cmd then
        btn.cmd()
      end
    end, {buffer = buf, silent = true, desc = btn.text})
  end
  
  -- Initial highlight
  vim.schedule(function()
    update_highlight()
  end)
end

function M.show()
  if vim.fn.argc() > 0 then return end
  
  if vim.fn.filereadable(M.config.animation_file) == 0 then
    vim.notify("Animation not found: " .. M.config.animation_file, vim.log.levels.WARN)
    return
  end
  
  vim.cmd('only')
  
  if M.config.orientation == "centered" then
    -- CENTERED MODE: Animation centered with margins, buttons below
    local total_width = vim.o.columns
    local margin = math.floor(total_width * 0.15)
    
    -- Create left margin
    vim.cmd('topleft vnew')
    local left_buf = vim.api.nvim_create_buf(false, true)
    vim.api.nvim_win_set_buf(0, left_buf)
    vim.api.nvim_win_set_width(0, margin)
    vim.bo[left_buf].bufhidden = 'wipe'
    
    -- Center column
    vim.cmd('wincmd l')
    local anim_buf = vim.api.nvim_create_buf(false, true)
    vim.api.nvim_set_current_buf(anim_buf)
    vim.bo[anim_buf].bufhidden = 'wipe'
    local anim_win = vim.api.nvim_get_current_win()
    
    -- Right margin
    vim.cmd('rightbelow vnew')
    local right_buf = vim.api.nvim_create_buf(false, true)
    vim.api.nvim_win_set_buf(0, right_buf)
    vim.api.nvim_win_set_width(0, margin)
    vim.bo[right_buf].bufhidden = 'wipe'
    
    -- Go back to center and split for buttons
    vim.api.nvim_set_current_win(anim_win)
    vim.cmd('rightbelow split')
    
    local buttons_buf, button_lines = create_buttons_buffer()
    local buttons_win = vim.api.nvim_get_current_win()
    vim.api.nvim_win_set_buf(buttons_win, buttons_buf)
    vim.api.nvim_win_set_height(buttons_win, #button_lines + 1)
    
    setup_buttons_navigation(buttons_buf)
    
    -- Clean window options
    for _, win in ipairs(vim.api.nvim_list_wins()) do
      vim.api.nvim_win_set_option(win, 'number', false)
      vim.api.nvim_win_set_option(win, 'relativenumber', false)
      vim.api.nvim_win_set_option(win, 'signcolumn', 'no')
      vim.api.nvim_win_set_option(win, 'cursorline', false)
      vim.api.nvim_win_set_option(win, 'colorcolumn', '')
      vim.api.nvim_win_set_option(win, 'foldcolumn', '0')
    end
    
    -- Start animation in the anim window
    vim.api.nvim_set_current_win(anim_win)
    local cmd = string.format("%s --raw %s", M.config.aart_binary, vim.fn.shellescape(M.config.animation_file))
    vim.fn.termopen(cmd, {
      on_exit = function()
        vim.schedule(function()
          pcall(vim.cmd, 'only')
          if M.config.auto_restore_session then
            local cwd = vim.fn.getcwd()
            if cwd ~= vim.fn.expand('~') and cwd ~= '/' then
              pcall(vim.cmd, 'SessionRestore')
            end
          end
        end)
      end
    })
    
    -- Switch focus to buttons window (not animation)
    vim.api.nvim_set_current_win(buttons_win)
    
    return
  end
  
  -- DEFAULT MODE: 3-panel layout
  local total_width = vim.o.columns
  local side_width = math.floor(total_width * 0.25)
  
  -- Left panel (calendar)
  vim.cmd('topleft vnew')
  local left_buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_win_set_buf(0, left_buf)
  vim.api.nvim_win_set_width(0, side_width)
  vim.api.nvim_buf_set_lines(left_buf, 0, -1, false, get_calendar())
  vim.bo[left_buf].modifiable = false
  vim.bo[left_buf].bufhidden = 'wipe'
  
  -- Center column
  vim.cmd('wincmd l')
  local center_win = vim.api.nvim_get_current_win()
  
  -- Right panel (system info)
  vim.cmd('rightbelow vnew')
  local right_buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_win_set_buf(0, right_buf)
  vim.api.nvim_win_set_width(0, side_width)
  vim.api.nvim_buf_set_lines(right_buf, 0, -1, false, get_system_info())
  vim.bo[right_buf].modifiable = false
  vim.bo[right_buf].bufhidden = 'wipe'
  
  -- Back to center, split for animation and buttons
  vim.api.nvim_set_current_win(center_win)
  
  local anim_buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_current_buf(anim_buf)
  vim.bo[anim_buf].bufhidden = 'wipe'
  local anim_win = vim.api.nvim_get_current_win()
  
  vim.cmd('rightbelow split')
  local buttons_buf, button_lines = create_buttons_buffer()
  local buttons_win = vim.api.nvim_get_current_win()
  vim.api.nvim_win_set_buf(buttons_win, buttons_buf)
  vim.api.nvim_win_set_height(buttons_win, #button_lines + 1)
  
  setup_buttons_navigation(buttons_buf)
  
  -- Clean window options
  for _, win in ipairs(vim.api.nvim_list_wins()) do
    vim.api.nvim_win_set_option(win, 'number', false)
    vim.api.nvim_win_set_option(win, 'relativenumber', false)
    vim.api.nvim_win_set_option(win, 'signcolumn', 'no')
    vim.api.nvim_win_set_option(win, 'cursorline', false)
    vim.api.nvim_win_set_option(win, 'colorcolumn', '')
    vim.api.nvim_win_set_option(win, 'foldcolumn', '0')
  end
  
  -- Start animation in anim window
  vim.api.nvim_set_current_win(anim_win)
  local cmd = string.format("%s --raw %s", M.config.aart_binary, vim.fn.shellescape(M.config.animation_file))
  vim.fn.termopen(cmd, {
    on_exit = function()
      vim.schedule(function()
        pcall(vim.cmd, 'only')
        if M.config.auto_restore_session then
          local cwd = vim.fn.getcwd()
          if cwd ~= vim.fn.expand('~') and cwd ~= '/' then
            pcall(vim.cmd, 'SessionRestore')
          end
        end
      end)
    end
  })
  
  -- Switch focus to buttons window (not animation in terminal mode)
  vim.api.nvim_set_current_win(buttons_win)
end

function M.setup_autocommands()
  local group = vim.api.nvim_create_augroup('AartStartup', {clear = true})
  vim.api.nvim_create_autocmd('VimEnter', {
    group = group,
    callback = function()
      vim.defer_fn(function()
        if vim.fn.argc() == 0 and vim.bo.filetype == '' then
          M.show()
        end
      end, 10)
    end,
  })
end

function M.setup_command()
  vim.api.nvim_create_user_command('AartStartup', M.show, {desc = 'Show aart startup'})
end

return M
