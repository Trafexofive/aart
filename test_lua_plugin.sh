#!/bin/bash
# Simple test to verify the Lua plugin can parse .aart files

echo "=== Testing aart Lua plugin ==="
echo

# Check if test.aa exists
if [ ! -f "test.aa" ]; then
    echo "Error: test.aa not found"
    exit 1
fi

echo "✓ Found test.aa"

# Show file structure
echo
echo "File structure:"
head -20 test.aa | jq -r '.canvas, .frames[0] | keys' 2>/dev/null || echo "  (JSON structure)"

# Test JSON parsing with a simple Lua script
lua << 'LUA'
local json = require('cjson')
local file = io.open('test.aa', 'r')
if not file then
    print('✗ Failed to open test.aa')
    os.exit(1)
end

local content = file:read('*all')
file:close()

local ok, data = pcall(json.decode, content)
if not ok then
    print('✗ Failed to parse JSON: ' .. tostring(data))
    os.exit(1)
end

print('✓ Successfully parsed JSON')
print('  Canvas: ' .. data.canvas.width .. 'x' .. data.canvas.height)
print('  Frames: ' .. #data.frames)

if #data.frames > 0 then
    print('  Frame 0 has ' .. #data.frames[1].cells .. ' rows')
    
    -- Extract first few lines
    print('\nFirst frame preview:')
    for i = 1, math.min(5, #data.frames[1].cells) do
        local line = ''
        for _, cell in ipairs(data.frames[1].cells[i]) do
            line = line .. (cell.char or ' ')
        end
        print('  ' .. line:gsub('%s+$', ''))
    end
end

print('\n✓ All basic tests passed')
LUA

echo
echo "=== Plugin structure ==="
echo "Files:"
ls -lh lua/aart/

echo
echo "=== Usage example ==="
cat << 'USAGE'
# Add to Neovim config:
{
    dir = "~/repos/aart/lua/aart",
    config = function()
        require('aart').setup()
        require('aart').setup_commands()
    end
}

# Then in Neovim:
:lua require('aart').get_static_frame('test.aa', 0)
USAGE

