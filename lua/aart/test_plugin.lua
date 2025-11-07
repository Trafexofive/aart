-- Test script for aart.nvim plugin
-- Run with: nvim -l test_plugin.lua

-- Add the lua directory to package path
package.path = package.path .. ";./lua/?.lua;./lua/?/init.lua"

-- Mock vim API for testing outside Neovim
if not vim then
  _G.vim = {
    fn = {
      filereadable = function(path)
        local f = io.open(path, "r")
        if f then
          f:close()
          return 1
        end
        return 0
      end,
      shellescape = function(s) return "'" .. s .. "'" end,
      tempname = function() return "/tmp/test_aart_" .. os.time() end,
      mkdir = function(path) os.execute("mkdir -p " .. path) end,
      getcwd = function() return "." end,
      expand = function(s) return s:gsub("^~", os.getenv("HOME")) end,
    },
    notify = function(msg, level)
      print(string.format("[%s] %s", level or "INFO", msg))
    end,
    log = {
      levels = { ERROR = "ERROR", WARN = "WARN", INFO = "INFO" }
    },
    json = {
      decode = function(str)
        -- Simple JSON decode for testing
        -- In real Neovim, use vim.json.decode
        local json = require('cjson')  -- or use another JSON lib
        return json.decode(str)
      end
    },
    tbl_deep_extend = function(behavior, ...)
      local result = {}
      for _, tbl in ipairs({...}) do
        for k, v in pairs(tbl) do
          result[k] = v
        end
      end
      return result
    end,
  }
end

-- Load the plugin
local aart = require('aart')

print("=== Testing aart.nvim plugin ===\n")

-- Test 1: Setup
print("Test 1: Setup")
aart.setup({
  cache_enabled = true,
  auto_start = true,
})
print("✓ Setup complete\n")

-- Test 2: Load static frame
print("Test 2: Load static frame from test.aa")
local frame = aart.get_static_frame("test.aa", 0)
if frame and #frame > 0 then
  print("✓ Loaded " .. #frame .. " lines")
  print("First few lines:")
  for i = 1, math.min(5, #frame) do
    print("  " .. frame[i])
  end
else
  print("✗ Failed to load frame")
end
print()

-- Test 3: Load animation
print("Test 3: Load full animation")
local animation = aart.load_animation("test.aa")
if animation then
  print("✓ Loaded animation:")
  print("  Frames: " .. animation.frame_count)
  print("  FPS: " .. animation.fps)
  print("  Frame delay: " .. math.floor(1000 / animation.fps) .. "ms")
else
  print("✗ Failed to load animation")
end
print()

-- Test 4: Export frames
print("Test 4: Export frames")
local export_dir = "/tmp/aart_test_frames"
if aart.export_frames("test.aa", export_dir) then
  print("✓ Frames exported to " .. export_dir)
  -- List exported files
  local handle = io.popen("ls -1 " .. export_dir .. " | head -5")
  if handle then
    print("  Files:")
    for line in handle:lines() do
      print("    " .. line)
    end
    handle:close()
  end
else
  print("✗ Failed to export frames")
end
print()

-- Test 5: Alpha-nvim integration
print("Test 5: Get alpha-nvim formatted frame")
local alpha_frame = aart.get_alpha_animation("test.aa", {
  frame = 0,
  padding = 1,
})
if alpha_frame and #alpha_frame > 0 then
  print("✓ Got alpha frame (" .. #alpha_frame .. " lines with padding)")
  print("Preview (first 3 lines):")
  for i = 1, math.min(3, #alpha_frame) do
    print("  " .. (alpha_frame[i] or ""))
  end
else
  print("✗ Failed to get alpha frame")
end
print()

print("=== Tests complete ===")
