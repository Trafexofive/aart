package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

// ImportGIFScreen handles GIF import with options
type ImportGIFScreen struct {
	width      int
	height     int
	theme      Theme
	styles     Styles
	config     *config.Config
	returnTo   tea.Model
	
	// Import options
	url        string
	targetWidth  int
	targetHeight int
	fps          int
	method       string
	ratio        string
	
	// UI state
	inputMode    string // "url", "width", "height", "fps", "method", "ratio", "confirm"
	cursor       int
	methods      []string
	ratios       []string
}

// NewImportGIFScreen creates GIF import dialog
func NewImportGIFScreen(cfg *config.Config, returnTo tea.Model) ImportGIFScreen {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	// Get terminal size for smart defaults
	termWidth, termHeight := GetTerminalSize()
	
	// Use 80% of terminal size for canvas, accounting for UI chrome
	defaultWidth := int(float64(termWidth) * 0.8)
	defaultHeight := int(float64(termHeight-10) * 0.8) // Subtract UI elements
	
	// Clamp to reasonable values
	if defaultWidth < 40 {
		defaultWidth = 80
	}
	if defaultHeight < 20 {
		defaultHeight = 30
	}
	
	return ImportGIFScreen{
		theme:        theme,
		styles:       NewStyles(theme),
		config:       cfg,
		returnTo:     returnTo,
		targetWidth:  defaultWidth,
		targetHeight: defaultHeight,
		fps:          cfg.Editor.DefaultFPS,
		method:       "block",
		ratio:        "fill",
		inputMode:    "url",
		cursor:       0,
		methods:      []string{"luminosity", "average", "block", "dither"},
		ratios:       []string{"fill", "fit", "original"},
	}
}

func (g ImportGIFScreen) Init() tea.Cmd {
	return nil
}

func (g ImportGIFScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Return to previous screen
			if g.returnTo != nil {
				return g.returnTo, nil
			}
			return NewStartupPage(g.config), nil
		
		case "tab":
			// Cycle through input fields
			g.advanceInputMode()
		
		case "enter":
			if g.inputMode == "confirm" {
				// Start import
				return g.startImport()
			}
			g.advanceInputMode()
		
		case "backspace":
			if g.inputMode == "url" && len(g.url) > 0 {
				g.url = g.url[:len(g.url)-1]
			}
		
		case "up", "k":
			if g.inputMode == "method" {
				g.cursor = (g.cursor - 1 + len(g.methods)) % len(g.methods)
				g.method = g.methods[g.cursor]
			} else if g.inputMode == "ratio" {
				g.cursor = (g.cursor - 1 + len(g.ratios)) % len(g.ratios)
				g.ratio = g.ratios[g.cursor]
			}
		
		case "down", "j":
			if g.inputMode == "method" {
				g.cursor = (g.cursor + 1) % len(g.methods)
				g.method = g.methods[g.cursor]
			} else if g.inputMode == "ratio" {
				g.cursor = (g.cursor + 1) % len(g.ratios)
				g.ratio = g.ratios[g.cursor]
			}
		
		case "+":
			if g.inputMode == "width" {
				g.targetWidth += 10
			} else if g.inputMode == "height" {
				g.targetHeight += 5
			} else if g.inputMode == "fps" {
				g.fps += 5
			}
		
		case "-":
			if g.inputMode == "width" && g.targetWidth > 20 {
				g.targetWidth -= 10
			} else if g.inputMode == "height" && g.targetHeight > 10 {
				g.targetHeight -= 5
			} else if g.inputMode == "fps" && g.fps > 5 {
				g.fps -= 5
			}
		
		default:
			// Type into URL field
			if g.inputMode == "url" {
				g.url += msg.String()
			}
		}
	
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
	}
	
	return g, nil
}

func (g *ImportGIFScreen) advanceInputMode() {
	modes := []string{"url", "width", "height", "fps", "method", "ratio", "confirm"}
	for i, mode := range modes {
		if mode == g.inputMode {
			g.inputMode = modes[(i+1)%len(modes)]
			g.cursor = 0
			
			// Set cursor for method/ratio selection
			if g.inputMode == "method" {
				for i, m := range g.methods {
					if m == g.method {
						g.cursor = i
						break
					}
				}
			} else if g.inputMode == "ratio" {
				for i, r := range g.ratios {
					if r == g.ratio {
						g.cursor = i
						break
					}
				}
			}
			return
		}
	}
}

func (g ImportGIFScreen) startImport() (tea.Model, tea.Cmd) {
	if g.url == "" {
		return g, nil
	}
	
	// Create import options
	opts := &ImportOptions{
		URL:    g.url,
		Width:  g.targetWidth,
		Height: g.targetHeight,
		FPS:    g.fps,
		Method: g.method,
		Ratio:  g.ratio,
	}
	
	// Start import with progress screen
	return NewImportProgressScreen(g.config, opts, g.returnTo), nil
}

func (g ImportGIFScreen) View() string {
	if g.width == 0 {
		return "Loading..."
	}
	
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(g.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)
	
	b.WriteString(titleStyle.Render("🎬 Import GIF to ASCII"))
	b.WriteString("\n\n")
	
	labelStyle := lipgloss.NewStyle().
		Foreground(g.theme.FgSecondary).
		Width(20)
	
	activeStyle := lipgloss.NewStyle().
		Foreground(g.theme.AccentPrimary).
		Bold(true)
	
	valueStyle := lipgloss.NewStyle().
		Foreground(g.theme.AccentSecondary).
		Bold(true)
	
	// URL input
	prefix := "  "
	if g.inputMode == "url" {
		prefix = "▶ "
		b.WriteString(activeStyle.Render(prefix + labelStyle.Render("GIF URL/Path:")))
		b.WriteString(" ")
		b.WriteString(valueStyle.Render(g.url))
		b.WriteString(lipgloss.NewStyle().Foreground(g.theme.Cursor).Render("▌"))
	} else {
		b.WriteString(labelStyle.Render(prefix + "GIF URL/Path:"))
		b.WriteString(" ")
		if g.url == "" {
			b.WriteString(lipgloss.NewStyle().Foreground(g.theme.FgMuted).Render("(empty)"))
		} else {
			b.WriteString(g.url)
		}
	}
	b.WriteString("\n\n")
	
	// Width
	prefix = "  "
	if g.inputMode == "width" {
		prefix = "▶ "
		b.WriteString(activeStyle.Render(prefix + labelStyle.Render("Width:")))
		b.WriteString(" ")
		b.WriteString(valueStyle.Render(fmt.Sprintf("%d", g.targetWidth)))
		b.WriteString(lipgloss.NewStyle().Foreground(g.theme.FgMuted).Render(" (+/- to adjust)"))
	} else {
		b.WriteString(labelStyle.Render(prefix + "Width:"))
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%d", g.targetWidth))
	}
	b.WriteString("\n")
	
	// Height
	prefix = "  "
	if g.inputMode == "height" {
		prefix = "▶ "
		b.WriteString(activeStyle.Render(prefix + labelStyle.Render("Height:")))
		b.WriteString(" ")
		b.WriteString(valueStyle.Render(fmt.Sprintf("%d", g.targetHeight)))
		b.WriteString(lipgloss.NewStyle().Foreground(g.theme.FgMuted).Render(" (+/- to adjust)"))
	} else {
		b.WriteString(labelStyle.Render(prefix + "Height:"))
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%d", g.targetHeight))
	}
	b.WriteString("\n")
	
	// FPS
	prefix = "  "
	if g.inputMode == "fps" {
		prefix = "▶ "
		b.WriteString(activeStyle.Render(prefix + labelStyle.Render("FPS:")))
		b.WriteString(" ")
		b.WriteString(valueStyle.Render(fmt.Sprintf("%d", g.fps)))
		b.WriteString(lipgloss.NewStyle().Foreground(g.theme.FgMuted).Render(" (+/- to adjust)"))
	} else {
		b.WriteString(labelStyle.Render(prefix + "FPS:"))
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%d", g.fps))
	}
	b.WriteString("\n\n")
	
	// Method selection
	prefix = "  "
	if g.inputMode == "method" {
		prefix = "▶ "
		b.WriteString(activeStyle.Render(prefix + labelStyle.Render("Method:")))
		b.WriteString("\n")
		for i, method := range g.methods {
			if i == g.cursor {
				b.WriteString(valueStyle.Render(fmt.Sprintf("     ▶ %s", method)))
			} else {
				b.WriteString(fmt.Sprintf("       %s", method))
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString(labelStyle.Render(prefix + "Method:"))
		b.WriteString(" ")
		b.WriteString(g.method)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	
	// Ratio selection
	prefix = "  "
	if g.inputMode == "ratio" {
		prefix = "▶ "
		b.WriteString(activeStyle.Render(prefix + labelStyle.Render("Aspect Ratio:")))
		b.WriteString("\n")
		for i, ratio := range g.ratios {
			desc := ""
			switch ratio {
			case "fill":
				desc = " - Fill canvas (may stretch)"
			case "fit":
				desc = " - Fit inside canvas (preserve ratio)"
			case "original":
				desc = " - Use original GIF size"
			}
			if i == g.cursor {
				b.WriteString(valueStyle.Render(fmt.Sprintf("     ▶ %s", ratio)))
				b.WriteString(lipgloss.NewStyle().Foreground(g.theme.FgMuted).Italic(true).Render(desc))
			} else {
				b.WriteString(fmt.Sprintf("       %s", ratio))
				b.WriteString(lipgloss.NewStyle().Foreground(g.theme.FgMuted).Italic(true).Render(desc))
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString(labelStyle.Render(prefix + "Aspect Ratio:"))
		b.WriteString(" ")
		b.WriteString(g.ratio)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	
	// Confirm button
	if g.inputMode == "confirm" {
		confirmStyle := lipgloss.NewStyle().
			Background(g.theme.AccentSuccess).
			Foreground(g.theme.BgPrimary).
			Bold(true).
			Padding(0, 2)
		b.WriteString("  ")
		b.WriteString(confirmStyle.Render("▶ START IMPORT ◀"))
	} else {
		b.WriteString("  [Press Tab to continue, Enter when ready]")
	}
	
	// Instructions
	b.WriteString("\n\n")
	hintStyle := lipgloss.NewStyle().
		Foreground(g.theme.FgMuted).
		Italic(true)
	
	hints := "Tab: next field │ Enter: confirm │ +/-: adjust values │ hjkl: navigate lists │ Esc: cancel"
	b.WriteString(hintStyle.Render(hints))
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(g.theme.BorderActive).
		Padding(2, 4).
		Width(90)
	
	return lipgloss.Place(
		g.width,
		g.height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(b.String()),
	)
}

// GetTerminalSize returns the current terminal dimensions
func GetTerminalSize() (width, height int) {
	// Default fallback
	width, height = 120, 40
	
	// Try stty size command
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err == nil {
		fmt.Sscanf(string(out), "%d %d", &height, &width)
	}
	
	// Ensure reasonable minimums
	if width < 80 {
		width = 120
	}
	if height < 24 {
		height = 40
	}
	
	return width, height
}

// ImportOptions holds GIF import parameters
type ImportOptions struct {
	URL    string
	Width  int
	Height int
	FPS    int
	Method string
	Ratio  string
}
