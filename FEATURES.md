# aart Feature Showcase

## Quick Start

```bash
make run
```

## Feature Demos

### 1. Navigation
```
Press: h j k l (or arrow keys)
Result: Yellow cursor (┃) moves around the canvas
```

### 2. Drawing - Single Character
```
Press: d d d d d
Result: Places █ characters at cursor position
```

### 3. Insert Mode
```
Press: i
Type: HELLO WORLD
Press: Esc
Result: Characters appear on canvas as you type
```

### 4. Animation Playback
```
Press: Space
Result: Timeline indicator (▓▓) cycles through frames 1-24
        Status shows "▶ playing"
Press: Space again
Result: Playback stops
```

### 5. Frame Seeking
```
Press: . . . (period)
Result: Jump forward 3 frames
Press: , , (comma)  
Result: Jump backward 2 frames
```

### 6. Tool Selection
```
Press: p (pencil)
Press: f (fill)
Press: s (select)
Result: Status bar updates to show current tool
```

### 7. Radial Wheel Menu
```
Press: Ctrl+j
Result: Cycles to next wheel section
Press: Ctrl+j Ctrl+j
Result: Cycles through HELP → EXPORT → IMPORT
Press: Enter
Result: Expands selected section showing details
Press: Esc
Result: Collapses back to cycling mode
Press: Esc again
Result: Hides wheel completely
```

### 8. Wheel Sections Detail

**HELP** (expanded):
- Shows all keyboard shortcuts
- Organized by category

**TOOLS** (expanded):
- Lists all 8 tools with keys
- Shows brush size selector
- Indicates active tool with ●

**COLORS** (stubbed):
- Color mode selector
- FG/BG swatches
- Palette grid
- Recent colors

### 9. Zoom Controls
```
Press: + + +
Result: Zoom level increases (1.0 → 1.25 → 1.5 → 1.75)
Press: -
Result: Zoom level decreases (1.75 → 1.5)
Press: 0
Result: Reset to 1.0x zoom
```

### 10. Command Mode
```
Press: :
Type: export out.ans
Press: Enter
Result: (Stubbed - command entered but not executed yet)

Press: :
Press: Esc
Result: Command cancelled, back to normal mode
```

### 11. View State
```
Press: g
Result: Grid toggle on/off (state tracked, overlay TBD)
```

## Demo Art

Frame 1 includes pre-loaded Braille art:
- Box drawing characters (╔═╗║╚╝)
- Braille patterns (⣿⣀⠁⠈⢸⣀⣼⡿⡄)
- Demonstrates Unicode support

## Visual Indicators

- **Cursor**: Yellow `┃` shows editing position
- **Current Frame**: `▓▓` in timeline
- **Modified Frame**: Frame number shown in green
- **Playing**: `▶ playing` in timeline status
- **Selected Wheel**: `◄` points to active section
- **Modified File**: `*` in status bar
- **Insert Mode**: `-- INSERT --` in bottom bar
- **Command Mode**: `:` prefix in bottom bar

## Status Bar Fields (Left to Right)

```
aart v0.1.0 │ filename.aart * │ 80x24 │ frame 3/24 │ 12fps │ pencil │ fg:█ bg:  │ layer 2/2
     │             │         │      │        │       │      │      │       │         │
  version      filename   modified size   current  fps   tool   colors  active/total
                                           frame                        layers
```

## Timeline Fields

```
▓▓  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24
││
│└─ Frame numbers (modified frames shown in green)
└── Current frame indicator (moves during playback)

▶ playing │ 83ms/frame │ loop: on │ ctrl-j/k: wheel │ [space] pause
│              │              │            │                  │
state      frame time      mode      help hint          next action
```

## Performance

- Smooth 12fps playback (83ms/frame)
- Instant cursor movement
- No lag on key input
- Efficient rendering (partial updates)

## Try This Sequence

```
1. make run
2. Press: i
3. Type: ASCII ART!
4. Press: Esc
5. Press: Space (start animation)
6. Press: , , , (seek back 3 frames)
7. Press: i
8. Type: FRAME 21
9. Press: Esc Space (toggle playback)
10. Press: Ctrl+j Enter (open TOOLS wheel)
11. Press: Esc Esc (close wheel)
12. Press: q (quit)
```

Result: You've navigated, drawn, animated, and explored the UI!

## Coming Soon

Check STARTING_BUILD_REF.md for planned features:
- Flood fill tool
- Line and box drawing
- Layer compositing
- Export to ANSI/PNG
- Import from images
- Undo/redo
- Copy/paste
