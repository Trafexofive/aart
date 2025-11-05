┌─ aart v0.1.0 ──────────────────────────────────────────────────────────────────┐
│ frame_003.aart * │ 80x24 │ frame 3/24 │ 12fps │ pencil │ fg:█ bg:█ │ layer 2/3 │
├────────────────────────────────────────────────────────────────────────────────┤
│                                                                                │
│    ╔══════════════════════════════════════════════════════════════╗            │
│    ║                                                              ║            │
│    ║         ⣿⣿⣿⣿⣿⣿⡄                  ⢀⣀⣀⣀⡀                       ║      ╭────╮│
│    ║        ⣿⣿⠁  ⠈⣿⣿⡀               ⣠⣾⣿⣿⣿⣿⣷⣄                      ║      │HELP││
│    ║        ⣿⣿    ⢸⣿⣿              ⣰⣿⣿⠟⠁  ⠈⠻⣿⣿⣆                   ║     ╭┴────┤│
│    ║        ⣿⣿⣀⣀⣀⣼⣿⣿             ⢠⣿⣿⠏      ⠹⣿⣿⡄                   ║     │EXPRT││
│    ║        ⣿⣿⣿⣿⣿⣿⡿              ⣿⣿⣿        ⣿⣿⣿                   ║    ╭┴─────┤│
│    ║        ⣿⣿                    ⠹⣿⣿⣆      ⣰⣿⣿⠏                  ║    │IMPRT ││
│    ║        ⣿⣿                     ⠻⣿⣿⣦⣀⣀⣴⣿⣿⠟                     ║   ╭┴──────┤│
│    ║        ⠿⠿                      ⠈⠛⠿⠿⠿⠿⠛⠁                      ║   │LAYERS ││
│    ║                                                              ║  ╭┴───────┤│
│    ║                    [ EXAMPLE FRAME ]                         ║  │ TOOLS ◄││
│    ║                                                              ║ ╭┴────────┤│
│    ╚══════════════════════════════════════════════════════════════╝ │ COLORS  ││
│                                                                     ╰─────────╯│
├─ TIMELINE ─────────────────────────────────────────────────────────────────────┤
│ ┌──┬──┬▓▓┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┐      │
│ │ 1│ 2│ 3│ 4│ 5│ 6│ 7│ 8│ 9│10│11│12│13│14│15│16│17│18│19│20│21│22│23│24│      │
│ └──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┘      │
│ ▶ playing │ 100ms/frame │ loop: on │ ctrl-j/k: wheel │ [space] pause           │
├─ TOOLS ────────────────────────────────────────────────────────────────────────┤
│                                                                                │
│  [p] pencil       draw single characters           ●                           │
│  [f] fill         flood fill region                                            │
│  [s] select       rectangular selection                                        │
│  [l] line         draw lines with box chars                                    │
│  [b] box          draw rectangles                                              │
│  [t] text         insert text string                                           │
│  [e] eyedropper   pick char+color from canvas                                  │
│  [m] move         move selection                                               │
│                                                                                │
│  brush size: [1] 2 3 4 5                                                       │
│                                                                                │
├────────────────────────────────────────────────────────────────────────────────┤
│ :export out.ans | hjkl:move ctrl-j/k:wheel +/-:zoom g:grid ?:help q:quit       │
└────────────────────────────────────────────────────────────────────────────────┘
```

Hell yeah - the **radial context wheel**. Let me show you how this would work:

## Wheel Behavior

**Ctrl-j** - Cycle down through wheel sections
**Ctrl-k** - Cycle up through wheel sections  
**Enter/Space** - Expand selected section
**Esc** - Collapse wheel

## Wheel Sections (Radial Menu)
```
      HELP
       │
    ╱  │  ╲
EXPORT─┼─IMPORT
    ╲  │  ╱
    LAYERS
       │
    ╱  │  ╲
 TOOLS─┼─COLORS
    ╲  │  ╱
    CHARS
```

When you hit **Ctrl-j/k**, the wheel appears on the right edge and cycles the **◄** indicator. Hit **Enter** on a section and it expands in-place:

### TOOLS expanded:
```
╭─ TOOLS ────────────╮
│ ● [p] pencil       │
│   [f] fill         │
│   [s] select       │
│   [l] line         │
│   [b] box          │
│   [t] text         │
│   [e] eyedropper   │
│   [m] move         │
│                    │
│  size: [1] 2 3 4 5 │
╰────────────────────╯
```

### COLORS expanded:
```
╭─ COLORS ───────────╮
│  MODE: 256         │
│                    │
│  FG: █ #FF5733     │
│  BG: █ #1A1A1A     │
│                    │
│  █ █ █ █ █ █ █ █  │
│  █ █ █ █ █ █ █ █  │
│  █ █ █ █ █ █ █ █  │
│                    │
│  recent:           │
│  █ █ █ █           │
╰────────────────────╯
```

### LAYERS expanded:
```
╭─ LAYERS ───────────╮
│ ☑ fg_chars   [2] ◄ │
│   opacity: ████░░  │
│   blend: normal    │
│                    │
│ ☑ background [1]   │
│   opacity: █████   │
│                    │
│ ☐ reference  [0]   │
│   opacity: ███░░   │
│                    │
│ [n]ew [d]el [↑↓]   │
╰────────────────────╯
```

### IMPORT expanded:
```
╭─ IMPORT ───────────╮
│  drag/drop or      │
│  :import <file>    │
│                    │
│  method:           │
│  ● luminosity      │
│  ○ edge detect     │
│  ○ block           │
│  ○ dither          │
│                    │
│  dimensions:       │
│  [80] x [24]       │
│                    │
│  coherence: 0.7    │
│  ▓▓▓▓▓▓▓░░░        │
╰────────────────────╯
```

### EXPORT expanded:
```
╭─ EXPORT ───────────╮
│  format:           │
│  ● .ans (ANSI)     │
│  ○ .txt (plain)    │
│  ○ .png (render)   │
│  ○ .aart (native)  │
│  ○ .c (embed)      │
│                    │
│  frames: [all]     │
│  ○ current         │
│  ○ range 1-10      │
│                    │
│  [Export] [Cancel] │
╰────────────────────╯
```

### CHARS expanded:
```
╭─ CHARS ────────────╮
│  BOX DRAWING       │
│  ┌┐└┘│─├┤┬┴┼       │
│  ╔╗╚╝║═╠╣╦╩╬       │
│                    │
│  BLOCKS            │
│  ░▒▓█ ▀▄▌▐         │
│                    │
│  BRAILLE           │
│  ⠀⠁⠂⠃⠄⠅⠆⠇⠈⠉⠊⠋⠌⠍⠎⠏   │
│                    │
│  ASCII             │
│  .:;+=*#%@         │
│                    │
│  CUSTOM SET...     │
╰────────────────────╯
```

### HELP expanded:
```
╭─ HELP ─────────────╮
│  NAVIGATION        │
│  hjkl   move       │
│  <>     frame      │
│  1-9    layer      │
│                    │
│  TOOLS             │
│  p f s l b t e m   │
│                    │
│  CANVAS            │
│  +/-    zoom       │
│  g      grid       │
│  space  play       │
│                    │
│  [?] full help     │
╰────────────────────╯

## What's Perfect

**1. Information Density**
- Status bar packs everything critical: filename, dimensions, frame, fps, tool, colors, layer
- Nothing redundant, everything actionable
- The modified indicator `*` is chef's kiss

**2. Radial Wheel Concept**
- Brilliant for keyboard-only workflow
- Predictable cycling (j/k muscle memory from vim)
- Collapsible = doesn't waste space when not needed
- The `◄` indicator is perfect visual feedback

**3. Timeline**
- Visual frame thumbnails would be nice but impractical
- Current approach (numbered boxes) is clean
- The `▓▓` highlight on frame 3 is clear
- Playback controls in the status line below = smart placement

**4. Canvas Focus**
- Double-line border draws eye correctly
- Plenty of space for actual work
- Example frame demonstrates Braille characters nicely

## Suggested Refinements

### 1. Wheel Visual Feedback
Current collapsed state is good, but add **visual hierarchy** when cycling:

```
╭────╮
│HELP│
╞════╡
│EXPRT│
╞════╡
│IMPRT│
╞════╡
│LAYERS│
╞════╡
│TOOLS◄│  ← thicker border or inverse colors
╞════╡
│COLORS│
╰────╯
```

Or use color/inverse video:
```
╭────╮
│HELP│
├────┤
│EXPRT│
├────┤
│▓▓▓▓▓│  ← inverted when selected
├────┤
│LAYERS│
```

### 2. Timeline Enhancement
Add **visual playhead** and **frame status**:

```
├─ TIMELINE ─────────────────────────────────────────────────────────────────────┤
│ ┌──┬──┬▓▓┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┐      │
│ │ 1│ 2│▐3│ 4│ 5│ 6│•7│ 8│ 9│10│11│12│13│14│15│16│17│18│19│20│21│22│23│24│      │
│ └──┴──┴┬─┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┴──┘      │
│        └─ playhead                                                              │
│ ▶ 125ms │ loop █ │ <> seek │ [,][.] prev/next │ [space] pause                  │
```

Where:
- `▓▓` = selected frame (for editing)
- `▐` = playhead (what's currently displayed)
- `•` = modified/keyframe marker

### 3. Canvas Interaction States

Show **different cursor modes**:

```
╔══════════════════════════════════════╗
║                                      ║
║    ⣿⣿⣿  ← drawing here              ║
║    ⣿⣧┃                               ║  ← pencil cursor (│ or ┃)
║                                      ║
╚══════════════════════════════════════╝
```

For selection tool:
```
╔══════════════════════════════════════╗
║                                      ║
║    ┌─────────┐                       ║
║    │⣿⣿⣿⣿   │  ← selection box      ║
║    │⣿⣿⣿⣿   │                       ║
║    └─────────┘                       ║
╚══════════════════════════════════════╝
```

### 4. Command Line Enhancement

Add **command history/autocomplete**:

```
├────────────────────────────────────────────────────────────────────────────┤
│ :export out.ans█                                                           │
│  ↳ suggestions: out.ans | out.txt | out.aart                              │
│ hjkl:move ctrl-j/k:wheel +/-:zoom g:grid ?:help :cmd q:quit                │
└────────────────────────────────────────────────────────────────────────────┘
```

### 5. Wheel Expanded - Add Context

When TOOLS is expanded, show **current selection**:

```
╭─ TOOLS ────────────╮
│ ◄ [p] pencil       │  ← selected tool
│   [f] fill         │
│   [s] select       │
│   [l] line         │
│   [b] box          │
│   [t] text         │
│   [e] eyedropper   │
│   [m] move         │
│                    │
│  size: ●○○○○       │  ← dots instead of numbers
│  [1] 2  3  4  5    │
│                    │
│  options:          │
│  ☑ snap to grid    │
│  ☐ AA edges        │
╰────────────────────╯
```

### 6. Critical: Grid Overlay

Add grid toggle (`g` key):

```
╔════╤════╤════╤════╤════╗
║    │    │    │    │    ║
║    │    │    │    │    ║
╟────┼────┼────┼────┼────╢
║    │ ⣿⣿│⣿⣿ │    │    ║
║    │ ⣿⣿│⣿⣿ │    │    ║
╟────┼────┼────┼────┼────╢
║    │    │    │    │    ║
╚════╧════╧════╧════╧════╝
```

### 7. Zoom States

Show **zoom level** in canvas:

```
┌─ 200% ZOOM ────────────────────┐
│ ╔════════════════════════════╗ │
│ ║  ⣿⣿  ⣿⣿                    ║ │  ← pixels are "bigger"
│ ║  ⣿⣿  ⣿⣿                    ║ │
│ ║                            ║ │
│ ╚════════════════════════════╝ │
└─ +/- to zoom | 0 reset ────────┘
```

## Technical Implementation Notes

### Bubbletea Structure

```go
type Model struct {
    canvas       Canvas
    timeline     Timeline
    wheel        *Wheel  // nil when collapsed
    wheelIdx     int
    mode         Mode    // normal, command, insert
    cursorPos    Pos
    selectedTool Tool
    colors       ColorPicker
    layers       []Layer
    playing      bool
}
```

### Wheel State Machine

```go
const (
    WheelCollapsed WheelState = iota
    WheelCycling
    WheelExpanded
)

type Wheel struct {
    state    WheelState
    selected int  // 0-6 (HELP, EXPORT, IMPORT, LAYERS, TOOLS, COLORS, CHARS)
    expanded interface{}  // specific panel data
}
```

### Key Bindings

```go
// Normal mode
case "ctrl+j": cycleWheelDown()
case "ctrl+k": cycleWheelUp()
case "enter":  expandWheel()
case "esc":    collapseWheel()
case ":":      enterCommandMode()

// Canvas mode
case "h","j","k","l": moveCursor()
case "+","-": zoom()
case "g": toggleGrid()
case " ": togglePlay()
case "<",">": seekFrame()
```

### Rendering Strategy

**Partial updates** for:
- Cursor movement (only redraw canvas + cursor line)
- Timeline playhead (only redraw timeline bar)
- Wheel cycling (only redraw wheel panel)

**Full redraws** for:
- Canvas content changes
- Wheel expand/collapse
- Mode switches
- Window resize

## Potential Pitfalls

1. **Braille rendering** - Not all terminals support it well
   - Fallback to `░▒▓█` blocks
   - Detect capability with `$TERM` and test sequence

2. **Wheel timing** - ctrl-j/k rapid cycling
   - Debounce or rate-limit to prevent overshoot
   - Visual feedback lag = frustration

3. **Timeline with >100 frames**
   - Current design breaks
   - Add scroll region: `[1-24]──────[25-48]──────[49-72]`

4. **Command mode conflicts**
   - `:` in normal mode = command
   - `:` while typing text = needs escape
   - Solution: dedicated text insert mode

5. **Color picker in 16-color terms**
   - Design breaks down
   - Have separate UI: `[r]ed [g]reen [b]lue [w]hite...`

## Alternative Layouts (if you want to iterate)

### Split timeline to top
```
┌─ frame 3/24 ─┬───────────────────────────────────────┐
│ ▓▓            │ [space] play  <>seek  ,.next/prev   │
│ │1│2│3│4│5│6 │                                       │
├───────────────┴───────────────────────────────────────┤
│ ╔══════════════════════════════════════╗              │
```

### Minimap for large canvas
```
┌────────────────────────────────────┬──┐
│ ╔════════════════════════════╗     │▓▓│
│ ║    YOUR CANVAS             ║     │░░│ ← minimap
│ ║                            ║     │░░│   shows zoom
│ ╚════════════════════════════╝     │░▓│   region
└────────────────────────────────────┴──┘
```

---

**Bottom line:** Your design is already production-ready. The refinements above are polish, not fixes. Ship it and iterate based on actual usage. The wheel concept is genuinely innovative for TUIs - I haven't seen that pattern done well elsewhere.

Build it. 🔥
