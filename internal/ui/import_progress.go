package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
	"github.com/mlamkadm/aart/internal/converter"
)

// convertConverterToUIFrames converts converter frames to UI frames
func convertConverterToUIFrames(convFrames []*converter.Frame) []*Frame {
	uiFrames := make([]*Frame, len(convFrames))
	
	for i, cf := range convFrames {
		uf := &Frame{
			Width:  cf.Width,
			Height: cf.Height,
			Delay:  cf.Delay,
		}
		
		// Copy cells
		uf.Cells = make([][]Cell, cf.Height)
		for y := 0; y < cf.Height; y++ {
			uf.Cells[y] = make([]Cell, cf.Width)
			for x := 0; x < cf.Width; x++ {
				uf.Cells[y][x] = Cell{
					Char: cf.Cells[y][x].Char,
					FG:   cf.Cells[y][x].FG,
					BG:   cf.Cells[y][x].BG,
				}
			}
		}
		
		uiFrames[i] = uf
	}
	
	return uiFrames
}

// ImportProgressScreen shows GIF import progress
type ImportProgressScreen struct {
	width    int
	height   int
	theme    Theme
	styles   Styles
	config   *config.Config
	returnTo tea.Model
	opts     *ImportOptions
	
	// Progress tracking
	status       string
	progress     float64
	currentFrame int
	totalFrames  int
	done         bool
	err          error
	result       []*converter.Frame
	startTime    time.Time
}

// NewImportProgressScreen creates import progress display
func NewImportProgressScreen(cfg *config.Config, opts *ImportOptions, returnTo tea.Model) ImportProgressScreen {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	p := ImportProgressScreen{
		theme:     theme,
		styles:    NewStyles(theme),
		config:    cfg,
		returnTo:  returnTo,
		opts:      opts,
		status:    "Initializing...",
		startTime: time.Now(),
	}
	
	return p
}

func (p ImportProgressScreen) Init() tea.Cmd {
	return tea.Batch(
		p.startImport(),
		tickCmd(),
	)
}

type importDoneMsg struct {
	frames []*converter.Frame
	err    error
}

func (p ImportProgressScreen) startImport() tea.Cmd {
	return func() tea.Msg {
		// Import the GIF
		frames, err := converter.ImportGIF(
			p.opts.URL,
			p.opts.Width,
			p.opts.Height,
			p.opts.FPS,
			p.opts.Method,
			p.opts.Ratio,
		)
		
		return importDoneMsg{
			frames: frames,
			err:    err,
		}
	}
}

type importProgressMsg struct {
	status       string
	progress     float64
	currentFrame int
	totalFrames  int
}

func (p ImportProgressScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if p.done {
			switch msg.String() {
			case "enter", " ":
				if p.err == nil && len(p.result) > 0 {
					// Success - convert to UI frames and open editor
					uiFrames := convertConverterToUIFrames(p.result)
					return newModelWithConfig(uiFrames, "imported.aart", p.config), nil
				}
				// Error - return to previous
				if p.returnTo != nil {
					return p.returnTo, nil
				}
				return NewStartupPage(p.config), nil
			case "q", "esc":
				if p.returnTo != nil {
					return p.returnTo, nil
				}
				return NewStartupPage(p.config), nil
			}
		}
	
	case importDoneMsg:
		p.done = true
		p.err = msg.err
		p.result = msg.frames
		if msg.err == nil {
			p.status = "✓ Import complete!"
			p.progress = 1.0
		} else {
			p.status = fmt.Sprintf("✗ Error: %v", msg.err)
		}
	
	case importProgressMsg:
		p.status = msg.status
		p.progress = msg.progress
		p.currentFrame = msg.currentFrame
		p.totalFrames = msg.totalFrames
	
	case tickMsg:
		if !p.done {
			return p, tickCmd()
		}
	
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	}
	
	return p, nil
}

func (p ImportProgressScreen) View() string {
	// Use defaults if not set yet
	width := p.width
	height := p.height
	if width == 0 {
		width = 120
		height = 40
	}
	
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(p.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)
	
	b.WriteString(titleStyle.Render("🎬 Importing GIF"))
	b.WriteString("\n\n")
	
	// Source info
	infoStyle := lipgloss.NewStyle().
		Foreground(p.theme.FgSecondary)
	
	b.WriteString(infoStyle.Render(fmt.Sprintf("Source: %s", truncate(p.opts.URL, 60))))
	b.WriteString("\n")
	b.WriteString(infoStyle.Render(fmt.Sprintf("Target: %dx%d @ %dfps", p.opts.Width, p.opts.Height, p.opts.FPS)))
	b.WriteString("\n")
	b.WriteString(infoStyle.Render(fmt.Sprintf("Method: %s | Ratio: %s", p.opts.Method, p.opts.Ratio)))
	b.WriteString("\n\n")
	
	// Progress bar
	barWidth := 60
	filled := int(float64(barWidth) * p.progress)
	
	filledStyle := lipgloss.NewStyle().Foreground(p.theme.AccentPrimary)
	emptyStyle := lipgloss.NewStyle().Foreground(p.theme.FgMuted)
	
	for i := 0; i < barWidth; i++ {
		if i < filled {
			b.WriteString(filledStyle.Render("█"))
		} else {
			b.WriteString(emptyStyle.Render("░"))
		}
	}
	
	percentStyle := lipgloss.NewStyle().
		Foreground(p.theme.AccentSecondary).
		Bold(true)
	b.WriteString(" ")
	b.WriteString(percentStyle.Render(fmt.Sprintf("%3.0f%%", p.progress*100)))
	b.WriteString("\n\n")
	
	// Status
	var statusStyle lipgloss.Style
	if p.err != nil {
		statusStyle = lipgloss.NewStyle().
			Foreground(p.theme.AccentError).
			Bold(true)
	} else if p.done {
		statusStyle = lipgloss.NewStyle().
			Foreground(p.theme.AccentSuccess).
			Bold(true)
	} else {
		statusStyle = lipgloss.NewStyle().
			Foreground(p.theme.AccentInfo)
	}
	
	b.WriteString(statusStyle.Render(p.status))
	b.WriteString("\n")
	
	// Frame count if available
	if p.totalFrames > 0 {
		frameStyle := lipgloss.NewStyle().
			Foreground(p.theme.FgMuted)
		b.WriteString(frameStyle.Render(fmt.Sprintf("Processing frame %d/%d", p.currentFrame, p.totalFrames)))
		b.WriteString("\n")
	}
	
	// Elapsed time
	elapsed := time.Since(p.startTime)
	timeStyle := lipgloss.NewStyle().
		Foreground(p.theme.FgMuted)
	b.WriteString(timeStyle.Render(fmt.Sprintf("Elapsed: %.1fs", elapsed.Seconds())))
	b.WriteString("\n\n")
	
	// Instructions
	if p.done {
		hintStyle := lipgloss.NewStyle().
			Foreground(p.theme.AccentPrimary).
			Bold(true)
		
		if p.err == nil {
			b.WriteString(hintStyle.Render("Press Enter to open in editor"))
		} else {
			b.WriteString(hintStyle.Render("Press any key to return"))
		}
	} else {
		spinnerStyle := lipgloss.NewStyle().
			Foreground(p.theme.AccentInfo)
		
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinner := spinners[int(elapsed.Milliseconds()/100)%len(spinners)]
		
		b.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Processing...", spinner)))
	}
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.theme.BorderActive).
		Padding(2, 4).
		Width(80)
	
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(b.String()),
	)
}
