-- aart.nvim - ASCII Art Animation integration for Neovim
-- Spawns the aart binary to PLAY animations, not just show static frames

local M = {}

-- Default configuration
M.config = {
  binary_path = "aart",
  default_animation = nil,
}

-- Setup function
function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", M.config, opts or {})
end

-- Open animation in floating terminal (PLAYS THE ACTUAL ANIMATION!)
function M.open_animation(filepath)
  filepath = filepath or M.config.default_animation
  
  if not filepath then
    vim.notify("No animation file specified", vim.log.levels.ERROR)
    return
  end
  
  if vim.fn.filereadable(filepath) == 0 then
    vim.notify("Animation file not found: " .. filepath, vim.log.levels.ERROR)
    return
  end
  
  -- Calculate floating window size
  local width = math.floor(vim.o.columns * 0.9)
  local height = math.floor(vim.o.lines * 0.9)
  local row = math.floor((vim.o.lines - height) / 2)
  local col = math.floor((vim.o.columns - width) / 2)
  
  -- Create buffer
  local buf = vim.api.nvim_create_buf(false, true)
  
  -- Create floating window
  local win = vim.api.nvim_open_win(buf, true, {
    relative = 'editor',
    width = width,
    height = height,
    row = row,
    col = col,
    style = 'minimal',
    border = 'rounded',
  })
  
  -- Run aart binary in terminal mode to PLAY the animation
  local cmd = string.format("%s %s", M.config.binary_path, vim.fn.shellescape(filepath))
  vim.fn.termopen(cmd, {
    on_exit = function()
      if vim.api.nvim_win_is_valid(win) then
        vim.api.nvim_win_close(win, true)
      end
    end
  })
  
  -- Set buffer options
  vim.bo[buf].filetype = 'aart'
  vim.wo[win].number = false
  vim.wo[win].relativenumber = false
  vim.wo[win].signcolumn = 'no'
  
  -- Add keymaps to close (q and Esc)
  vim.keymap.set('n', 'q', '<cmd>quit<CR>', { buffer = buf, silent = true })
  vim.keymap.set('n', '<Esc>', '<cmd>quit<CR>', { buffer = buf, silent = true })
  vim.keymap.set('t', '<Esc>', '<C-\\><C-n>:quit<CR>', { buffer = buf, silent = true })
  
  -- Enter terminal mode so the TUI works
  vim.cmd('startinsert')
end

-- Get static frame for alpha-nvim dashboard
function M.get_static_frame(filepath, frame_index)
  frame_index = frame_index or 0
  
  local file = io.open(filepath, "r")
  if not file then
    return {}
  end
  
  local content = file:read("*all")
  file:close()
  
  local ok, data = pcall(vim.json.decode, content)
  if not ok or not data or not data.frames or not data.frames[frame_index + 1] then
    return {}
  end
  
  -- Extract ASCII from cells
  local lines = {}
  for _, row in ipairs(data.frames[frame_index + 1].cells) do
    local line = ""
    for _, cell in ipairs(row) do
      line = line .. (cell.char or " ")
    end
    line = line:gsub("%s+$", "")
    table.insert(lines, line)
  end
  
  -- Trim empty lines
  while #lines > 0 and lines[#lines] == "" do
    table.remove(lines)
  end
  while #lines > 0 and lines[1] == "" do
    table.remove(lines, 1)
  end
  
  return lines
end

-- Get frame for alpha-nvim with padding
function M.get_alpha_animation(filepath, opts)
  opts = opts or {}
  local frame = M.get_static_frame(filepath, opts.frame or 0)
  
  if #frame == 0 then
    return { "Failed to load animation" }
  end
  
  if opts.padding then
    for i = 1, opts.padding do
      table.insert(frame, 1, "")
      table.insert(frame, "")
    end
  end
  
  return frame
end

-- Register commands
function M.setup_commands()
  vim.api.nvim_create_user_command('AartOpen', function(opts)
    M.open_animation(opts.args ~= "" and opts.args or nil)
  end, {
    nargs = '?',
    complete = 'file',
    desc = 'Open ASCII art animation (plays with aart binary)'
  })
end

return M
