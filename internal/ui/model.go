package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
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
	Delay    int // milliseconds per frame
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
	zenMode    bool  // Minimal UI mode
	
	// Theme and styles
	theme      Theme
	styles     Styles
	breathing  *BreathingEffect
	
	// Command mode
	command    string
	statusMsg  string
	
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
	return NewWithConfig(nil)
}

// NewWithConfig creates a model with configuration
func NewWithConfig(cfg *config.Config) Model {
	// Use defaults if no config provided
	if cfg == nil {
		c := config.DefaultConfig
		cfg = &c
	}

	// Create initial frame
	width, height := cfg.Editor.DefaultWidth, cfg.Editor.DefaultHeight
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
			frame.Cells[y][x] = Cell{Char: r, FG: cfg.Colors.Foreground, BG: cfg.Colors.Background}
		}
	}
	
	// Create 24 frames
	frames := make([]*Frame, 24)
	frames[0] = frame
	for i := 1; i < 24; i++ {
		frames[i] = NewFrame(width, height)
	}
	
	return newModelWithConfig(frames, "untitled.aart", cfg)
}

// NewWithFrames creates a model with imported frames
func NewWithFrames(importedFrames []ImportedFrame) Model {
	return NewWithFramesAndConfig(importedFrames, nil)
}

// NewWithFramesAndConfig creates a model with imported frames and config
func NewWithFramesAndConfig(importedFrames []ImportedFrame, cfg *config.Config) Model {
	if len(importedFrames) == 0 {
		return NewWithConfig(cfg)
	}
	
	// Use defaults if no config provided
	if cfg == nil {
		c := config.DefaultConfig
		cfg = &c
	}
	
	// Convert imported frames to internal format
	frames := make([]*Frame, len(importedFrames))
	for i, imported := range importedFrames {
		frame := &Frame{
			Width:    imported.Width,
			Height:   imported.Height,
			Cells:    make([][]Cell, imported.Height),
			Modified: false,
		}
		
		for y := 0; y < imported.Height; y++ {
			frame.Cells[y] = make([]Cell, imported.Width)
			for x := 0; x < imported.Width; x++ {
				frame.Cells[y][x] = Cell{
					Char: imported.Cells[y][x].Char,
					FG:   imported.Cells[y][x].FG,
					BG:   imported.Cells[y][x].BG,
				}
			}
		}
		
		frames[i] = frame
	}
	
	return newModelWithConfig(frames, "imported.aart", cfg)
}

// ImportedFrame represents a frame from the converter
type ImportedFrame struct {
	Width  int
	Height int
	Cells  [][]ImportedCell
	Delay  int
}

// ImportedCell represents a cell from the converter
type ImportedCell struct {
	Char rune
	FG   string
	BG   string
}

func newModel(frames []*Frame, filename string) Model {
	return newModelWithConfig(frames, filename, nil)
}

func newModelWithConfig(frames []*Frame, filename string, cfg *config.Config) Model {
	// Use defaults if no config provided
	if cfg == nil {
		c := config.DefaultConfig
		cfg = &c
	}

	// Get theme from config or use default
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	styles := NewStyles(theme)

	return Model{
		mode:         ModeNormal,
		cursor:       Pos{X: 40, Y: 12},
		frames:       frames,
		currentFrame: 0,
		playing:      false,
		fps:          cfg.Editor.DefaultFPS,
		selectedTool: ToolPencil,
		fgChar:       '█',
		fgColor:      cfg.Colors.Foreground,
		bgColor:      cfg.Colors.Background,
		brushSize:    1,
		zoom:         1.0,
		showGrid:     cfg.Editor.ShowGrid,
		zenMode:      cfg.Editor.ZenMode,
		theme:        theme,
		styles:       styles,
		breathing:    NewBreathingEffect(3 * time.Second),
		filename:     filename,
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
	case "z":
		m.zenMode = !m.zenMode
		
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
	
	cmd := parts[0]
	
	switch cmd {
	case "save", "w":
		// Save to current filename or specify new one
		filename := m.filename
		if len(parts) > 1 {
			filename = parts[1]
		}
		if err := m.saveToFile(filename); err != nil {
			m.statusMsg = fmt.Sprintf("Error saving: %v", err)
		} else {
			m.statusMsg = fmt.Sprintf("Saved to %s", filename)
			m.modified = false
		}
	
	case "export":
		if len(parts) < 2 {
			m.statusMsg = "Usage: :export <file.ext> [frame]"
			return
		}
		filename := parts[1]
		frameIdx := -1 // All frames
		if len(parts) > 2 {
			fmt.Sscanf(parts[2], "%d", &frameIdx)
		}
		if err := m.exportToFile(filename, frameIdx); err != nil {
			m.statusMsg = fmt.Sprintf("Export error: %v", err)
		} else {
			m.statusMsg = fmt.Sprintf("Exported to %s", filename)
		}
	
	case "import":
		// TODO: implement import
		m.statusMsg = "Import not yet implemented"
	
	case "q", "quit":
		if m.modified {
			m.statusMsg = "Unsaved changes! Use :q! to force quit or :wq to save and quit"
		} else {
			// TODO: trigger quit
			m.statusMsg = "Quit"
		}
	
	case "q!":
		// Force quit
		// TODO: trigger quit
		m.statusMsg = "Force quit"
	
	case "wq":
		// Save and quit
		if err := m.saveToFile(m.filename); err != nil {
			m.statusMsg = fmt.Sprintf("Error saving: %v", err)
		} else {
			// TODO: trigger quit
			m.statusMsg = "Saved and quit"
		}
	
	case "new":
		// Create new animation
		m.statusMsg = "New animation created"
		// TODO: reset to blank canvas
	
	case "help":
		m.statusMsg = "Commands: :save :export :import :new :quit :help"
	
	default:
		m.statusMsg = fmt.Sprintf("Unknown command: %s", cmd)
	}
}

// saveToFile saves the current animation to a .aart file
func (m *Model) saveToFile(filename string) error {
	// Import the fileformat package at the top if not already
	aartFile := convertModelToAart(m, filename)
	
	// Use the fileformat package to save
	// For now, we'll implement a simple JSON save
	// TODO: Use fileformat.Save() once imported
	data, err := json.Marshal(aartFile)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	m.filename = filename
	return nil
}

// exportToFile exports to various formats
func (m *Model) exportToFile(filename string, frameIdx int) error {
	aartFile := convertModelToAart(m, filename)
	
	// Determine format from extension
	ext := ""
	if idx := strings.LastIndex(filename, "."); idx >= 0 {
		ext = filename[idx+1:]
	}
	
	// Simple export based on extension
	switch ext {
	case "json":
		data, err := json.MarshalIndent(aartFile, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(filename, data, 0644)
	
	case "txt", "ans":
		// Export as plain text
		frame := m.frames[0]
		if frameIdx >= 0 && frameIdx < len(m.frames) {
			frame = m.frames[frameIdx]
		}
		
		var output strings.Builder
		for y := 0; y < frame.Height; y++ {
			for x := 0; x < frame.Width; x++ {
				output.WriteRune(frame.Cells[y][x].Char)
			}
			output.WriteRune('\n')
		}
		
		return os.WriteFile(filename, []byte(output.String()), 0644)
	
	default:
		return fmt.Errorf("unsupported export format: %s", ext)
	}
}

// convertModelToAart converts the model to AartFile structure
func convertModelToAart(m *Model, filename string) map[string]interface{} {
	// Create a simplified version for now
	// TODO: Use actual fileformat.AartFile once imported
	frames := make([]map[string]interface{}, len(m.frames))
	for i, frame := range m.frames {
		cells := make([][]map[string]string, frame.Height)
		for y := 0; y < frame.Height; y++ {
			cells[y] = make([]map[string]string, frame.Width)
			for x := 0; x < frame.Width; x++ {
				cells[y][x] = map[string]string{
					"char": string(frame.Cells[y][x].Char),
					"fg":   frame.Cells[y][x].FG,
					"bg":   frame.Cells[y][x].BG,
				}
			}
		}
		frames[i] = map[string]interface{}{
			"index": i,
			"duration": frame.Delay,
			"cells": cells,
		}
	}
	
	return map[string]interface{}{
		"version": "1.0",
		"metadata": map[string]interface{}{
			"title": filename,
			"created": time.Now().Format(time.RFC3339),
			"modified": time.Now().Format(time.RFC3339),
		},
		"canvas": map[string]int{
			"width":  m.frames[0].Width,
			"height": m.frames[0].Height,
		},
		"frames": frames,
	}
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	
	var b strings.Builder
	
	// Zen mode - only canvas
	if m.zenMode {
		return m.renderCanvasOnly()
	}
	
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

func (m Model) renderCanvasOnly() string {
	frame := m.frames[m.currentFrame]
	var b strings.Builder
	
	// Just the canvas content, no borders
	for y := 0; y < min(frame.Height, m.height); y++ {
		for x := 0; x < min(frame.Width, m.width); x++ {
			cell := frame.Cells[y][x]
			
			// Show cursor in zen mode too
			if x == m.cursor.X && y == m.cursor.Y && m.mode != ModeCommand {
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("11")).
					Render("┃"))
			} else {
				b.WriteString(string(cell.Char))
			}
		}
		b.WriteString("\n")
	}
	
	// Minimal status at bottom in zen mode
	if m.mode == ModeCommand {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(fmt.Sprintf(":%s▌", m.command)))
	} else if m.mode == ModeInsert {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render("-- INSERT --"))
	} else {
		playIcon := "⏸"
		if m.playing {
			playIcon = "▶"
		}
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(fmt.Sprintf("%s %d/%d | z:exit zen", playIcon, m.currentFrame+1, len(m.frames))))
	}
	
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
			wheelWidth = 10 // Compact 4-letter wheel
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
	
	// Collapsed/Cycling wheel - compact 4-letter codes
	wheelLines := []string{
		"     ╔════╗",
		"     ║HELP║",
		"     ╠════╣",
		"     ║EXPT║",
		"     ╠════╣",
		"     ║IMPT║",
		"     ╠════╣",
		"     ║LAYR║",
		"     ╠════╣",
		"     ║TOOL║",
		"     ╠════╣",
		"     ║COLR║",
		"     ╚════╝",
	}
	
	// Add selection indicator
	var wheelLine string
	if line < len(wheelLines) {
		wheelLine = wheelLines[line]
		
		// Highlight selected section with ◄
		selectedLine := -1
		switch m.wheel.Selected {
		case WheelHelp:
			selectedLine = 1
		case WheelExport:
			selectedLine = 3
		case WheelImport:
			selectedLine = 5
		case WheelLayers:
			selectedLine = 7
		case WheelTools:
			selectedLine = 9
		case WheelColors:
			selectedLine = 11
		}
		
		if line == selectedLine {
			// Replace last border with indicator
			wheelLine = strings.TrimRight(wheelLine, "║")
			wheelLine += "◄"
			return lipgloss.NewStyle().
				Foreground(m.theme.AccentPrimary).
				Bold(true).
				Render(wheelLine)
		}
		
		return lipgloss.NewStyle().
			Foreground(m.theme.Border).
			Render(wheelLine)
	}
	return strings.Repeat(" ", 10)
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
