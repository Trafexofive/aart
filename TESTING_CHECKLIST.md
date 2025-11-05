# aart TUI Testing Checklist

## Test Date: 2024
## Version: Post 10x Redesign

---

## Startup Page Tests

### Visual Rendering
- [x] Logo renders centered and aligned
- [x] Tagline appears below logo
- [x] Version shows correctly
- [x] No ANSI escape sequences visible in output
- [x] Active panel has inverse title bar (▓▓▓▓)
- [x] Inactive panel has normal title
- [x] Selected menu item has inverse video
- [x] Active panel has thick border
- [x] Inactive panel has rounded border

### Navigation - Menu Panel
- [ ] 'j' or Down arrow moves selection down
- [ ] 'k' or Up arrow moves selection up
- [ ] Selection wraps at top/bottom
- [ ] Selected item shows ▸ indicator
- [ ] Selected item has inverse background

### Navigation - Tab Switching
- [ ] Tab key switches from menu to recent files
- [ ] Tab again switches back to menu
- [ ] Active panel visual indicator changes
- [ ] Border style changes (thick vs rounded)

### Navigation - Recent Files Panel
- [ ] 'j'/'k' moves through recent files list
- [ ] Number keys (1-8) directly select files
- [ ] Selected file shows ▸ indicator
- [ ] File info displays correctly (frames, time ago)

### Menu Actions
- [ ] 'n' - Creates new animation
- [ ] 'o' - Opens file picker
- [ ] 'i' - Opens GIF import screen
- [ ] 't' - Opens theme selector
- [ ] 's' - Opens settings
- [ ] 'c' - Edit config (launches $EDITOR)
- [ ] 'e' - Shows examples
- [ ] '?' - Shows help
- [ ] 'q' - Quits application

### Recent Files Actions
- [ ] Enter on selected file opens it
- [ ] Number key (1-8) opens corresponding file
- [ ] All 8 recent files display (if available)
- [ ] Time formatting works (1h, 13m, 7d)
- [ ] Frame count shows correctly

### Visual Polish
- [ ] Tip rotation counter shows (1/8, 2/8, etc.)
- [ ] Tip rotation icon animates (⟳ ⟲)
- [ ] Tips change every 5 seconds
- [ ] Config stats display (dimensions, FPS, method)
- [ ] Hints at bottom show all keys
- [ ] No visual glitches or artifacts

---

## Editor View Tests

### Visual Rendering
- [ ] Status bar shows filename, size, frame, FPS, tool, layer
- [ ] Canvas has single-line borders
- [ ] Wheel shows with compact 4-letter codes
- [ ] Timeline shows dots for frames
- [ ] Bottom status shows mode (NORMAL/INSERT/COMMAND)
- [ ] No escape sequences or glitches

### Canvas Rendering
- [ ] ASCII art displays correctly
- [ ] Cursor visible (┃ character)
- [ ] Cursor moves with hjkl
- [ ] Canvas content doesn't overflow borders
- [ ] Multiple frames render correctly

### Wheel Navigation
- [ ] Ctrl+J cycles wheel down
- [ ] Ctrl+K cycles wheel up
- [ ] Selected section highlighted with ◄
- [ ] Selected section is bold/bright
- [ ] All 6 sections accessible (HELP/EXPT/IMPT/LAYR/TOOL/COLR)
- [ ] Enter expands selected section
- [ ] Esc collapses wheel

### Wheel Sections
- [ ] HELP section shows shortcuts
- [ ] TOOL section shows all tools with [key]
- [ ] COLR section shows color picker
- [ ] LAYR section shows layers
- [ ] Selected tool marked with ●
- [ ] Tool section keyboard shortcuts work (p/f/s/l/b/t/e/m)

### Timeline
- [ ] Dots show for each frame (● for current, · for others)
- [ ] Current frame highlighted
- [ ] Scroll indicators (‹ ›) appear when >30 frames
- [ ] Window centers on current frame
- [ ] Modified frames show ◉ marker
- [ ] Frame counter shows (N/total)

### Playback
- [ ] Space toggles play/pause
- [ ] Status shows ▸ when playing
- [ ] FPS displayed correctly
- [ ] Frame advances during playback
- [ ] Loop indicator (↻) visible

### Frame Navigation
- [ ] ',' goes to previous frame
- [ ] '.' goes to next frame
- [ ] Frame counter updates
- [ ] Timeline updates current position
- [ ] Wraps at first/last frame

### Mode Switching
- [ ] 'i' enters INSERT mode
- [ ] Status shows "INSERT"
- [ ] Typing adds characters
- [ ] Arrow keys move cursor in insert mode
- [ ] Esc returns to NORMAL mode

### Command Mode
- [ ] ':' enters command mode
- [ ] Status shows "COMMAND"
- [ ] Command prompt visible with cursor
- [ ] Commands work (:quit, :save, :export, :help)
- [ ] Esc cancels command mode

### Drawing Tools
- [ ] 'p' selects pencil tool
- [ ] 'f' selects fill tool
- [ ] 's' selects select tool
- [ ] 'l' selects line tool
- [ ] 'b' selects box tool
- [ ] 't' selects text tool
- [ ] 'e' selects eyedropper
- [ ] 'm' selects move tool
- [ ] Selected tool shows in status bar
- [ ] Wheel TOOL section updates

### View Controls
- [ ] '+' zooms in
- [ ] '-' zooms out
- [ ] '0' resets zoom
- [ ] 'g' toggles grid
- [ ] 'z' toggles zen mode

### Zen Mode
- [ ] Only canvas visible (no borders/chrome)
- [ ] Cursor still visible
- [ ] Esc or 'z' exits zen mode
- [ ] Returns to full UI

### Help System
- [ ] '?' shows help screen
- [ ] All keyboard shortcuts listed
- [ ] Navigation works in help
- [ ] Esc or 'q' closes help

---

## File Operations Tests

### Save
- [ ] ':save filename.aa' saves file
- [ ] Modified indicator (*) appears in status
- [ ] Modified indicator clears after save
- [ ] Recent files list updates

### Export
- [ ] ':export output.ans' exports ANSI
- [ ] ':export output.txt' exports plain text
- [ ] ':export output.json' exports JSON
- [ ] ':export output.csv' exports CSV
- [ ] Files created successfully

### Open
- [ ] ':import file' loads file
- [ ] Canvas updates with loaded content
- [ ] Frames load correctly
- [ ] Animation plays correctly

---

## GIF Import Tests

### Basic Import
- [ ] '--import-gif file.gif' imports GIF
- [ ] '--import-gif URL' imports from web
- [ ] Frames convert correctly
- [ ] FPS matches or can be set

### Conversion Methods
- [ ] '--method luminosity' works
- [ ] '--method block' works
- [ ] '--method edge' works
- [ ] '--method dither' works
- [ ] Character sets correct for each method

### Sizing
- [ ] '--width N --height M' sets size
- [ ] Auto-size to terminal works
- [ ] '--ratio fill' stretches to fit
- [ ] '--ratio fit' preserves aspect
- [ ] '--ratio original' keeps original size

---

## Configuration Tests

### Config File
- [ ] '~/.config/aart/config.yml' exists
- [ ] '--init' creates config
- [ ] '--show-config' displays current config
- [ ] '--config-path' shows path
- [ ] Editing config works ('c' key on startup)

### Theme System
- [ ] 't' key shows theme selector
- [ ] tokyo-night theme works
- [ ] gruvbox theme works
- [ ] monokai theme works
- [ ] dracula theme works
- [ ] catppuccin theme works
- [ ] oceanic theme works
- [ ] Theme changes apply immediately

### Config Options
- [ ] default_width/height work
- [ ] default_fps works
- [ ] show_grid default works
- [ ] zen_mode default works
- [ ] converter.default_method works
- [ ] Custom startup artwork works

---

## Terminal Compatibility Tests

### Terminal Sizes
- [ ] Works on 80x24 (minimum)
- [ ] Works on 120x40 (medium)
- [ ] Works on 200x60 (large)
- [ ] Resize handling works
- [ ] Layout doesn't break on narrow terminals

### Terminal Types
- [ ] xterm-256color works
- [ ] screen-256color works
- [ ] tmux-256color works
- [ ] Basic 16-color terminal works (degraded but functional)

### Character Sets
- [ ] Box drawing characters render (─│┌┐└┘)
- [ ] Double borders render (═║╔╗╚╝)
- [ ] Braille characters render (⣿⠀⠁)
- [ ] Block characters render (░▒▓█)
- [ ] Unicode arrows render (▸◄●◉)

---

## Performance Tests

### Rendering Speed
- [ ] Startup page renders <100ms
- [ ] Editor view renders <100ms
- [ ] Navigation feels instant
- [ ] No visible lag on keypresses
- [ ] Timeline scrolls smoothly

### Large Files
- [ ] 150+ frame animation loads
- [ ] Timeline scroll window works
- [ ] Playback smooth with many frames
- [ ] Memory usage reasonable

---

## Regression Tests (Don't Break Existing)

### Existing Features
- [ ] All original keyboard shortcuts work
- [ ] File formats still load (.aa, .ans, .txt)
- [ ] Export formats still work
- [ ] GIF import still functional
- [ ] Config backward compatible

---

## Edge Cases

### Empty States
- [ ] No recent files shows helpful message
- [ ] Empty canvas works
- [ ] Single frame animation works
- [ ] Zero-size canvas prevented

### Error Handling
- [ ] Invalid file shows error
- [ ] Missing GIF shows error
- [ ] Bad command shows helpful message
- [ ] Network errors handled (URL import)
- [ ] Disk full handled gracefully

### Unusual Input
- [ ] Very long filename handled
- [ ] Special characters in filename
- [ ] Large frame count (1000+)
- [ ] Very small terminal handled
- [ ] Rapid keypress buffering works

---

## Final Checks

### Documentation
- [ ] README.md updated with screenshots
- [ ] Keyboard shortcuts documented
- [ ] Config options documented
- [ ] Examples work

### Code Quality
- [ ] No compiler warnings
- [ ] Build succeeds
- [ ] No obvious memory leaks
- [ ] Error messages helpful

### User Experience
- [ ] First-time user can figure it out
- [ ] Power user can work efficiently
- [ ] No confusing states
- [ ] Consistent behavior throughout

---

## Known Issues

### To Fix
1. Escape sequences sometimes visible in recent files panel (cosmetic)
2. Navigation (hjkl) in startup menu not visually updating (functional but not visible)

### To Investigate
1. Full redraw optimization for performance
2. Partial canvas updates
3. Mouse support

---

## Overall Assessment

**Status**: Ready for testing
**Blockers**: None critical
**Nice-to-haves**: Performance optimization, mouse support
**Ship-ready**: Yes (with minor cosmetic issues noted)
