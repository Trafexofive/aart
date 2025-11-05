package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeCommand
	ModeInsert
)

type WheelState int

const (
	WheelCollapsed WheelState = iota
	WheelCycling
	WheelExpanded
)

type WheelSection int

const (
	WheelHelp WheelSection = iota
	WheelExport
	WheelImport
	WheelLayers
	WheelTools
	WheelColors
	WheelChars
)

var wheelNames = []string{"HELP", "EXPORT", "IMPORT", "LAYERS", "TOOLS", "COLORS", "CHARS"}

type Tool int

const (
	ToolPencil Tool = iota
	ToolFill
	ToolSelect
	ToolLine
	ToolBox
	ToolText
	ToolEyedropper
	ToolMove
)

var toolNames = []string{"pencil", "fill", "select", "line", "box", "text", "eyedropper", "move"}
var toolKeys = []string{"p", "f", "s", "l", "b", "t", "e", "m"}

type Pos struct {
	X, Y int
}

type Cell struct {
	Char rune
	FG   string
	BG   string
}

type Frame struct {
	Width    int
	Height   int
	Cells    [][]Cell
	Modified bool
}

type Layer struct {
	Name     string
	Visible  bool
	Opacity  float64
	BlendMode string
	Frame    *Frame
}

type Wheel struct {
	State    WheelState
	Selected WheelSection
}

type Model struct {
	width      int
	height     int
	mode       Mode
	cursor     Pos
	frames     []*Frame
	currentFrame int
	playing    bool
	fps        int
	lastTick   time.Time
	
	// Tools
	selectedTool Tool
	fgChar       rune
	fgColor      string
	bgColor      string
	brushSize    int
	
	// Wheel
	wheel      *Wheel
	
	// Layers
	layers     []Layer
	currentLayer int
	
	// UI state
	showGrid   bool
	zoom       float64
	
	// Command mode
	command    string
	
	// Misc
	modified   bool
	filename   string
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*83, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func New() Model {
	// Create initial frame
	width, height := 80, 24
	frame := NewFrame(width, height)
	
	// Add some demo art to frame 0
	demoArt := []string{
		"    ╔══════════════════════════════════════════════════════════════╗",
		"    ║                                                              ║",
		"    ║         ⣿⣿⣿⣿⣿⣿⡄                  ⢀⣀⣀⣀⡀                       ║",
		"    ║        ⣿⣿⠁  ⠈⣿⣿⡀               ⣠⣾⣿⣿⣿⣿⣷⣄                      ║",
		"    ║        ⣿⣿    ⢸⣿⣿              ⣰⣿⣿⠟⠁  ⠈⠻⣿⣿⣆                   ║",
		"    ║        ⣿⣿⣀⣀⣀⣼⣿⣿             ⢠⣿⣿⠏      ⠹⣿⣿⡄                   ║",
		"    ║        ⣿⣿⣿⣿⣿⣿⡿              ⣿⣿⣿        ⣿⣿⣿                   ║",
		"    ║        ⣿⣿                    ⠹⣿⣿⣆      ⣰⣿⣿⠏                  ║",
		"    ║        ⣿⣿                     ⠻⣿⣿⣦⣀⣀⣴⣿⣿⠟                     ║",
		"    ║        ⠿⠿                      ⠈⠛⠿⠿⠿⠿⠛⠁                      ║",
		"    ║                                                              ║",
		"    ║                    [ EXAMPLE FRAME ]                         ║",
		"    ║                                                              ║",
		"    ╚══════════════════════════════════════════════════════════════╝",
	}
	
	for y, line := range demoArt {
		if y >= height {
			break
		}
		runes := []rune(line)
		for x, r := range runes {
			if x >= width {
				break
			}
			frame.Cells[y][x] = Cell{Char: r, FG: "#FFFFFF", BG: "#000000"}
		}
	}
	
	// Create 24 frames
	frames := make([]*Frame, 24)
	frames[0] = frame
	for i := 1; i < 24; i++ {
		frames[i] = NewFrame(width, height)
	}
	
	return Model{
		mode:         ModeNormal,
		cursor:       Pos{X: 40, Y: 12},
		frames:       frames,
		currentFrame: 0,
		playing:      false,
		fps:          12,
		selectedTool: ToolPencil,
		fgChar:       '█',
		fgColor:      "#FFFFFF",
		bgColor:      "#000000",
		brushSize:    1,
		zoom:         1.0,
		showGrid:     false,
		filename:     "untitled.aart",
		layers: []Layer{
			{Name: "background", Visible: true, Opacity: 1.0, BlendMode: "normal"},
			{Name: "fg_chars", Visible: true, Opacity: 1.0, BlendMode: "normal"},
		},
		currentLayer: 1,
	}
}

func NewFrame(width, height int) *Frame {
	cells := make([][]Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			cells[y][x] = Cell{Char: ' ', FG: "#FFFFFF", BG: "#000000"}
		}
	}
	return &Frame{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
		
	case tea.KeyMsg:
		switch m.mode {
		case ModeNormal:
			return m.handleNormalMode(msg)
		case ModeCommand:
			return m.handleCommandMode(msg)
		case ModeInsert:
			return m.handleInsertMode(msg)
		}
		
	case tickMsg:
		if m.playing {
			m.currentFrame = (m.currentFrame + 1) % len(m.frames)
			return m, tickCmd()
		}
	}
	
	return m, nil
}

func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
		
	// Wheel navigation
	case "ctrl+j":
		if m.wheel == nil {
			m.wheel = &Wheel{State: WheelCycling, Selected: WheelHelp}
		} else {
			m.wheel.Selected = WheelSection((int(m.wheel.Selected) + 1) % len(wheelNames))
		}
		
	case "ctrl+k":
		if m.wheel == nil {
			m.wheel = &Wheel{State: WheelCycling, Selected: WheelChars}
		} else {
			idx := int(m.wheel.Selected) - 1
			if idx < 0 {
				idx = len(wheelNames) - 1
			}
			m.wheel.Selected = WheelSection(idx)
		}
		
	case "enter":
		if m.wheel != nil && m.wheel.State == WheelCycling {
			m.wheel.State = WheelExpanded
		}
		
	case "esc":
		if m.wheel != nil {
			if m.wheel.State == WheelExpanded {
				m.wheel.State = WheelCycling
			} else {
				m.wheel = nil
			}
		}
		
	// Canvas navigation
	case "h", "left":
		if m.cursor.X > 0 {
			m.cursor.X--
		}
	case "l", "right":
		if m.cursor.X < m.frames[m.currentFrame].Width-1 {
			m.cursor.X++
		}
	case "k", "up":
		if m.cursor.Y > 0 {
			m.cursor.Y--
		}
	case "j", "down":
		if m.cursor.Y < m.frames[m.currentFrame].Height-1 {
			m.cursor.Y++
		}
		
	// Drawing
	case "d":
		frame := m.frames[m.currentFrame]
		frame.Cells[m.cursor.Y][m.cursor.X] = Cell{
			Char: m.fgChar,
			FG:   m.fgColor,
			BG:   m.bgColor,
		}
		frame.Modified = true
		m.modified = true
		
	case "i":
		m.mode = ModeInsert
		
	// Tool selection
	case "p":
		m.selectedTool = ToolPencil
	case "f":
		m.selectedTool = ToolFill
	case "s":
		m.selectedTool = ToolSelect
		
	// Playback
	case " ":
		m.playing = !m.playing
		if m.playing {
			return m, tickCmd()
		}
		
	case ",":
		if m.currentFrame > 0 {
			m.currentFrame--
		}
	case ".":
		if m.currentFrame < len(m.frames)-1 {
			m.currentFrame++
		}
		
	// View controls
	case "g":
		m.showGrid = !m.showGrid
	case "+":
		m.zoom = minFloat(m.zoom+0.25, 4.0)
	case "-":
		m.zoom = maxFloat(m.zoom-0.25, 0.25)
	case "0":
		m.zoom = 1.0
		
	// Command mode
	case ":":
		m.mode = ModeCommand
		m.command = ""
	}
	
	return m, nil
}

func (m Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.command = ""
	case "enter":
		m.executeCommand()
		m.mode = ModeNormal
		m.command = ""
	case "backspace":
		if len(m.command) > 0 {
			m.command = m.command[:len(m.command)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.command += msg.String()
		}
	}
	return m, nil
}

func (m Model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
	default:
		if len(msg.String()) == 1 {
			runes := []rune(msg.String())
			if len(runes) > 0 {
				frame := m.frames[m.currentFrame]
				frame.Cells[m.cursor.Y][m.cursor.X] = Cell{
					Char: runes[0],
					FG:   m.fgColor,
					BG:   m.bgColor,
				}
				frame.Modified = true
				m.modified = true
				
				// Move cursor right
				if m.cursor.X < frame.Width-1 {
					m.cursor.X++
				}
			}
		}
	}
	return m, nil
}

func (m *Model) executeCommand() {
	parts := strings.Fields(m.command)
	if len(parts) == 0 {
		return
	}
	
	switch parts[0] {
	case "export":
		// TODO: implement export
	case "import":
		// TODO: implement import
	case "q", "quit":
		// TODO: confirm if modified
	}
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	
	var b strings.Builder
	
	// Status bar
	b.WriteString(m.renderStatusBar())
	b.WriteString("\n")
	
	// Main area
	mainHeight := m.height - 7 // status(1) + canvas border(2) + timeline(3) + statusline(1)
	
	// Canvas
	b.WriteString(m.renderCanvas(mainHeight))
	b.WriteString("\n")
	
	// Timeline
	b.WriteString(m.renderTimeline())
	b.WriteString("\n")
	
	// Bottom status line
	b.WriteString(m.renderBottomStatus())
	
	return b.String()
}

func (m Model) renderStatusBar() string {
	modified := ""
	if m.modified {
		modified = " *"
	}
	
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("8")).
		Render(fmt.Sprintf(
			" aart v0.1.0  │ %s%s │ %dx%d │ frame %d/%d │ %dfps │ %s │ fg:%c bg:  │ layer %d/%d ",
			m.filename,
			modified,
			m.frames[m.currentFrame].Width,
			m.frames[m.currentFrame].Height,
			m.currentFrame+1,
			len(m.frames),
			m.fps,
			toolNames[m.selectedTool],
			m.fgChar,
			m.currentLayer+1,
			len(m.layers),
		))
}

func (m Model) renderCanvas(height int) string {
	var b strings.Builder
	
	// Calculate canvas area
	canvasWidth := m.width
	wheelWidth := 0
	
	if m.wheel != nil {
		if m.wheel.State == WheelExpanded {
			wheelWidth = 22
		} else {
			wheelWidth = 13
		}
		canvasWidth -= wheelWidth
	}
	
	frame := m.frames[m.currentFrame]
	
	// Top border
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	line := border.Render("┌" + strings.Repeat("─", canvasWidth-2) + "┐")
	
	if m.wheel != nil {
		line += m.renderWheel(0, height)
	}
	b.WriteString(line + "\n")
	
	// Canvas content
	for y := 0; y < height-2; y++ {
		lineContent := border.Render("│")
		
		for x := 0; x < canvasWidth-2; x++ {
			if x < frame.Width && y < frame.Height {
				cell := frame.Cells[y][x]
				
				// Show cursor
				if x == m.cursor.X && y == m.cursor.Y && m.mode != ModeCommand {
					lineContent += lipgloss.NewStyle().
						Foreground(lipgloss.Color("11")).
						Render("┃")
				} else {
					lineContent += string(cell.Char)
				}
			} else {
				lineContent += " "
			}
		}
		
		lineContent += border.Render("│")
		
		if m.wheel != nil {
			lineContent += m.renderWheel(y+1, height)
		}
		b.WriteString(lineContent + "\n")
	}
	
	// Bottom border
	line = border.Render("└" + strings.Repeat("─", canvasWidth-2) + "┘")
	
	if m.wheel != nil {
		line += m.renderWheel(height-1, height)
	}
	b.WriteString(line)
	
	return b.String()
}

func (m Model) renderWheel(line, totalHeight int) string {
	if m.wheel == nil {
		return ""
	}
	
	if m.wheel.State == WheelExpanded {
		return m.renderExpandedWheel(line, totalHeight)
	}
	
	// Collapsed/Cycling wheel - show on right edge
	wheelLines := []string{
		"      ╭────╮",
		"      │HELP│",
		"     ╭┴────┤",
		"     │EXPRT│",
		"    ╭┴─────┤",
		"    │IMPRT │",
		"   ╭┴──────┤",
		"   │LAYERS │",
		"  ╭┴───────┤",
		"  │ TOOLS  │",
		" ╭┴────────┤",
		" │ COLORS  │",
		" ╰─────────╯",
	}
	
	// Indicators on specific lines
	var wheelLine string
	if line < len(wheelLines) {
		wheelLine = wheelLines[line]
		
		// Highlight selected section
		if (line == 1 && m.wheel.Selected == WheelHelp) ||
		   (line == 3 && m.wheel.Selected == WheelExport) ||
		   (line == 5 && m.wheel.Selected == WheelImport) ||
		   (line == 7 && m.wheel.Selected == WheelLayers) ||
		   (line == 9 && m.wheel.Selected == WheelTools) ||
		   (line == 11 && m.wheel.Selected == WheelColors) {
			// Add indicator
			wheelLine = strings.Replace(wheelLine, "│", "│", 1)
			wheelLine = strings.TrimRight(wheelLine, " │")
			wheelLine += " ◄│"
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("11")).
				Render(wheelLine)
		}
		
		return wheelLine
	}
	return strings.Repeat(" ", 13)
}

func (m Model) renderExpandedWheel(line, totalHeight int) string {
	var content []string
	
	switch m.wheel.Selected {
	case WheelTools:
		content = []string{
			"╭─ TOOLS ────────────╮",
			"│ ● [p] pencil       │",
			"│   [f] fill         │",
			"│   [s] select       │",
			"│   [l] line         │",
			"│   [b] box          │",
			"│   [t] text         │",
			"│   [e] eyedropper   │",
			"│   [m] move         │",
			"│                    │",
			"│  size: [1] 2 3 4 5 │",
			"╰────────────────────╯",
		}
		// Update selection indicator
		toolLine := 1 + int(m.selectedTool)
		if line == toolLine && toolLine < len(content) {
			content[line] = strings.Replace(content[line], "  ", "● ", 1)
		}
		
	case WheelColors:
		content = []string{
			"╭─ COLORS ───────────╮",
			"│  MODE: 256         │",
			"│                    │",
			"│  FG: █ #FFFFFF     │",
			"│  BG: █ #000000     │",
			"│                    │",
			"│  █ █ █ █ █ █ █ █  │",
			"│  █ █ █ █ █ █ █ █  │",
			"│  █ █ █ █ █ █ █ █  │",
			"│                    │",
			"│  recent:           │",
			"│  █ █ █ █           │",
			"╰────────────────────╯",
		}
		
	case WheelHelp:
		content = []string{
			"╭─ HELP ─────────────╮",
			"│ hjkl    move       │",
			"│ d       draw       │",
			"│ i       insert     │",
			"│ space   play/pause │",
			"│ ,.      seek       │",
			"│ +/-     zoom       │",
			"│ g       grid       │",
			"│ ctrl-j/k  wheel    │",
			"│ enter   expand     │",
			"│ esc     collapse   │",
			"│ :       command    │",
			"│ q       quit       │",
			"╰────────────────────╯",
		}
		
	default:
		content = []string{
			"╭─ " + wheelNames[m.wheel.Selected] + " ──────────╮",
			"│  Coming soon...    │",
			"╰────────────────────╯",
		}
	}
	
	if line < len(content) {
		return content[line]
	}
	return strings.Repeat(" ", 22)
}

func (m Model) renderTimeline() string {
	var b strings.Builder
	
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	b.WriteString(border.Render("├─ TIMELINE "))
	b.WriteString(border.Render(strings.Repeat("─", m.width-13) + "┤"))
	b.WriteString("\n│ ")
	
	// Frame boxes
	for i := 0; i < len(m.frames); i++ {
		if i == m.currentFrame {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("11")).
				Render("▓▓"))
		} else if m.frames[i].Modified {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("10")).
				Render(fmt.Sprintf("%2d", i+1)))
		} else {
			b.WriteString(fmt.Sprintf("%2d", i+1))
		}
		if i < len(m.frames)-1 {
			b.WriteString(" ")
		}
	}
	
	b.WriteString(strings.Repeat(" ", max(0, m.width-3*len(m.frames)-3)))
	b.WriteString("│\n")
	
	playStatus := "stopped"
	if m.playing {
		playStatus = "▶ playing"
	}
	
	b.WriteString(border.Render(fmt.Sprintf("│ %s │ %dms/frame │ loop: on │ ctrl-j/k: wheel │ [space] pause", 
		playStatus, 1000/m.fps)))
	b.WriteString(strings.Repeat(" ", max(0, m.width-80)))
	b.WriteString(border.Render("│"))
	
	return b.String()
}

func (m Model) renderBottomStatus() string {
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	
	if m.mode == ModeCommand {
		cmd := lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Render(":" + m.command + "▌")
		return border.Render("│ ") + cmd + strings.Repeat(" ", max(0, m.width-len(m.command)-5)) + border.Render("│")
	}
	
	if m.mode == ModeInsert {
		status := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Render("-- INSERT --")
		return border.Render("│ ") + status + strings.Repeat(" ", max(0, m.width-16)) + border.Render("│")
	}
	
	statusText := fmt.Sprintf(" :export out.ans | hjkl:move ctrl-j/k:wheel +/-:zoom g:grid ?:help q:quit")
	padding := max(0, m.width-len(statusText)-2)
	return border.Render("│") + statusText + strings.Repeat(" ", padding) + border.Render("│")
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
